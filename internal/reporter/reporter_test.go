package reporter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fanfeilong/dot_20_arch_draft/internal/analyzer"
	"github.com/fanfeilong/dot_20_arch_draft/internal/deriver"
	"github.com/fanfeilong/dot_20_arch_draft/internal/installer"
	"github.com/fanfeilong/dot_20_arch_draft/internal/state"
	"github.com/fanfeilong/dot_20_arch_draft/internal/tester"
)

func TestBuildReport(t *testing.T) {
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
	if _, err := tester.PrepareTests(repo); err != nil {
		t.Fatalf("PrepareTests returned error: %v", err)
	}

	resolved, err := BuildReport(repo)
	if err != nil {
		t.Fatalf("BuildReport returned error: %v", err)
	}

	if resolved != target {
		t.Fatalf("unexpected resolved target: got %q want %q", resolved, target)
	}

	indexPath := filepath.Join(repo, "report", "index.md")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read report index: %v", err)
	}
	if !strings.Contains(string(content), target) {
		t.Fatalf("report index does not mention target repo")
	}

	for _, rel := range []string{
		filepath.Join("report", "data", "summary.json"),
		filepath.Join("report", "data", "target.json"),
		filepath.Join("report", "data", "tests.json"),
		filepath.Join("report", "data", "challenge.json"),
		filepath.Join("report", "vue-app", "package.json"),
		filepath.Join("report", "vue-app", "src", "App.vue"),
	} {
		if _, err := os.Stat(filepath.Join(repo, rel)); err != nil {
			t.Fatalf("expected report data file %s: %v", rel, err)
		}
	}
	if !strings.Contains(string(content), "Challenge") {
		t.Fatalf("report index does not contain challenge section")
	}
}

func TestBuildReportRequiresTestMini(t *testing.T) {
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

	if _, err := BuildReport(repo); err == nil {
		t.Fatalf("expected error when test-mini has not been prepared")
	}
}
