package tester

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fanfeilong/dot_20_arch_draft/internal/analyzer"
	"github.com/fanfeilong/dot_20_arch_draft/internal/deriver"
	"github.com/fanfeilong/dot_20_arch_draft/internal/installer"
	"github.com/fanfeilong/dot_20_arch_draft/internal/state"
)

func TestPrepareTests(t *testing.T) {
	repo := t.TempDir()
	if _, err := installer.Install(repo); err != nil {
		t.Fatalf("Install returned error: %v", err)
	}

	target := t.TempDir()
	if _, err := analyzer.Analyze(target, repo); err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}
	if _, err := state.RecordSkill(repo, state.SkillUpdate{
		Skill:          "d2a-challenge-architecture",
		Status:         "completed",
		Stage:          state.StageArchitectureChallengeComplete,
		Phase:          "challenge-dialogue",
		QuestionIndex:  6,
		QuestionTotal:  6,
		Recommendation: "proceed",
		Summary:        "Challenge phase complete.",
	}); err != nil {
		t.Fatalf("RecordSkill returned error: %v", err)
	}
	if _, err := deriver.DeriveMini(repo); err != nil {
		t.Fatalf("DeriveMini returned error: %v", err)
	}

	resolved, err := PrepareTests(repo)
	if err != nil {
		t.Fatalf("PrepareTests returned error: %v", err)
	}

	if resolved != target {
		t.Fatalf("unexpected resolved target: got %q want %q", resolved, target)
	}

	manifestPath := filepath.Join(repo, ".d2a", "test-mini.json")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf("expected test manifest %s: %v", manifestPath, err)
	}

	tasksPath := filepath.Join(repo, "tests", "01_integration_tasks.md")
	content, err := os.ReadFile(tasksPath)
	if err != nil {
		t.Fatalf("read integration tasks file: %v", err)
	}
	if !strings.Contains(string(content), target) {
		t.Fatalf("integration tasks file does not mention target repo")
	}
}

func TestPrepareTestsRequiresDeriveMini(t *testing.T) {
	repo := t.TempDir()
	if _, err := installer.Install(repo); err != nil {
		t.Fatalf("Install returned error: %v", err)
	}

	target := t.TempDir()
	if _, err := analyzer.Analyze(target, repo); err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}

	if _, err := PrepareTests(repo); err == nil {
		t.Fatalf("expected error when mini derivation is missing")
	}
}
