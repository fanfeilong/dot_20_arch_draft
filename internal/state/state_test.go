package state

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitializeAndRecordCommand(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".d2a"), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "LAB.md"), []byte("# d2a\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	s, err := Initialize(repo)
	if err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}
	if s.CurrentStage != StageInitialized {
		t.Fatalf("unexpected stage after initialize: %q", s.CurrentStage)
	}

	s, err = RecordCommand(repo, "d2a analyze", StageAnalysisPrepared, "Next.", "Prepared analysis.")
	if err != nil {
		t.Fatalf("RecordCommand returned error: %v", err)
	}
	if s.CurrentStage != StageAnalysisPrepared {
		t.Fatalf("unexpected stage after record: %q", s.CurrentStage)
	}
	if s.LastCommand != "d2a analyze" {
		t.Fatalf("unexpected last command: %q", s.LastCommand)
	}

	s, err = RecordSkill(repo, SkillUpdate{
		Skill:         "d2a-arch-1-project-scope",
		Status:        "progress",
		Stage:         StageAnalysisPrepared,
		Phase:         "confirmation-questions",
		QuestionIndex: 2,
		QuestionTotal: 5,
		NextStep:      "Continue the confirmation questions.",
		Summary:       "Moved into confirmation questions.",
	})
	if err != nil {
		t.Fatalf("RecordSkill returned error: %v", err)
	}
	if s.CurrentSkill != "d2a-arch-1-project-scope" {
		t.Fatalf("unexpected current skill: %q", s.CurrentSkill)
	}
	if s.CurrentPhase != "confirmation-questions" {
		t.Fatalf("unexpected current phase: %q", s.CurrentPhase)
	}
	if s.QuestionIndex != 2 || s.QuestionTotal != 5 {
		t.Fatalf("unexpected question progress: %d/%d", s.QuestionIndex, s.QuestionTotal)
	}

	history, err := RecentHistory(repo, 10)
	if err != nil {
		t.Fatalf("RecentHistory returned error: %v", err)
	}
	if len(history) != 3 {
		t.Fatalf("unexpected history length: got %d want 3", len(history))
	}
	if history[len(history)-1].ActorName != "d2a-arch-1-project-scope" {
		t.Fatalf("unexpected final history event: %+v", history[len(history)-1])
	}
}

func TestRecordChallengeSkill(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".d2a"), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "LAB.md"), []byte("# d2a\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	if _, err := Initialize(repo); err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	s, err := RecordSkill(repo, SkillUpdate{
		Skill:          "d2a-challenge-architecture",
		Status:         "progress",
		Stage:          StageArchitectureChallengeInProgress,
		Phase:          "challenge-dialogue",
		QuestionIndex:  3,
		QuestionTotal:  6,
		Decision:       "primary driver",
		Strength:       "partial",
		Recommendation: "review",
		Objection:      "Why not a simpler timer-based trigger?",
		Summary:        "Challenge round 3 is active.",
	})
	if err != nil {
		t.Fatalf("RecordSkill returned error: %v", err)
	}
	if s.CurrentDecision != "primary driver" {
		t.Fatalf("unexpected current decision: %q", s.CurrentDecision)
	}
	if s.LastChallengeStrength != "partial" {
		t.Fatalf("unexpected challenge strength: %q", s.LastChallengeStrength)
	}

	challenge, err := ChallengeState(repo)
	if err != nil {
		t.Fatalf("ChallengeState returned error: %v", err)
	}
	if challenge.DecisionIndex != 3 || challenge.DecisionTotal != 6 {
		t.Fatalf("unexpected challenge progress: %d/%d", challenge.DecisionIndex, challenge.DecisionTotal)
	}
	if challenge.Recommendation != "review" {
		t.Fatalf("unexpected recommendation: %q", challenge.Recommendation)
	}
	if _, err := os.Stat(filepath.Join(repo, ".d2a", "challenge_log.jsonl")); err != nil {
		t.Fatalf("expected challenge log file: %v", err)
	}
}

func TestRecordSkillPersistsQAArtifactsDuringConfirmation(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".d2a"), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "LAB.md"), []byte("# d2a\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if _, err := Initialize(repo); err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	if _, err := RecordSkill(repo, SkillUpdate{
		Skill:         "d2a-arch-2-runtime-view",
		Status:        "progress",
		Stage:         StageArchitectureInProgress,
		Phase:         "confirmation-questions",
		QuestionIndex: 1,
		QuestionTotal: 4,
		Question:      "What is the dominant runtime driver?",
		Answer:        "B",
		Evaluation:    "partial",
		Explanation:   "The loop trigger was close but not exact.",
		Summary:       "Question 1 reviewed.",
	}); err != nil {
		t.Fatalf("RecordSkill returned error: %v", err)
	}

	qaStatePath := filepath.Join(repo, ".d2a", "qa", "d2a-arch-2-runtime-view.json")
	if _, err := os.Stat(qaStatePath); err != nil {
		t.Fatalf("expected qa state file %s: %v", qaStatePath, err)
	}

	qaLogPath := filepath.Join(repo, ".d2a", "qa", "d2a-arch-2-runtime-view.jsonl")
	content, err := os.ReadFile(qaLogPath)
	if err != nil {
		t.Fatalf("read qa log file: %v", err)
	}
	text := string(content)
	if !strings.Contains(text, `"evaluation":"partial"`) {
		t.Fatalf("qa log does not contain evaluation: %q", text)
	}
	if !strings.Contains(text, `"question":"What is the dominant runtime driver?"`) {
		t.Fatalf("qa log does not contain question: %q", text)
	}
}

func TestRecordCommandRejectsUnknownStage(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".d2a"), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "LAB.md"), []byte("# d2a\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if _, err := Initialize(repo); err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	if _, err := RecordCommand(repo, "d2a fake", "unknown-stage", "Next.", "Nope."); err == nil {
		t.Fatalf("expected unknown stage error")
	}
}

func TestRecordSkillRejectsStageRegression(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".d2a"), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "LAB.md"), []byte("# d2a\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if _, err := Initialize(repo); err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}
	if _, err := RecordCommand(repo, "d2a report", StageReportReady, "Next.", "Prepared report."); err != nil {
		t.Fatalf("RecordCommand returned error: %v", err)
	}

	if _, err := RecordSkill(repo, SkillUpdate{
		Skill:   "d2a-mini-1-scope",
		Status:  "progress",
		Stage:   StageMiniDerivationPrepared,
		Summary: "Try to regress stage.",
	}); err == nil {
		t.Fatalf("expected stage regression error")
	}
}

func TestRecordSkillNormalizesNextFileToZhDocPath(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".d2a"), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(repo, "docs", "1.架构拆解"), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "LAB.md"), []byte("# d2a\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if _, err := Initialize(repo); err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	s, err := RecordSkill(repo, SkillUpdate{
		Skill:    "d2a-arch-2-runtime-view",
		Status:   "progress",
		Stage:    StageArchitectureInProgress,
		Phase:    "analysis-generation",
		NextFile: ".d2a/docs/architecture/02_driver.md",
		Summary:  "Normalize next file path.",
	})
	if err != nil {
		t.Fatalf("RecordSkill returned error: %v", err)
	}
	if s.NextFile != "docs/1.架构拆解/02_驱动.md" {
		t.Fatalf("unexpected normalized next file: %q", s.NextFile)
	}
}

func TestRecordSkillNormalizesNextFileWithoutZhLayout(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".d2a"), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "LAB.md"), []byte("# d2a\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if _, err := Initialize(repo); err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}

	s, err := RecordSkill(repo, SkillUpdate{
		Skill:    "d2a-arch-2-runtime-view",
		Status:   "progress",
		Stage:    StageArchitectureInProgress,
		Phase:    "analysis-generation",
		NextFile: ".d2a/docs/architecture/02_driver.md",
		Summary:  "Normalize next file path.",
	})
	if err != nil {
		t.Fatalf("RecordSkill returned error: %v", err)
	}
	if s.NextFile != "docs/architecture/02_driver.md" {
		t.Fatalf("unexpected normalized next file: %q", s.NextFile)
	}
}
