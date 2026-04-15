package deriver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fanfeilong/dot_20_arch_draft/internal/analyzer"
	"github.com/fanfeilong/dot_20_arch_draft/internal/installer"
	"github.com/fanfeilong/dot_20_arch_draft/internal/state"
)

func TestDeriveMini(t *testing.T) {
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

	resolved, err := DeriveMini(repo)
	if err != nil {
		t.Fatalf("DeriveMini returned error: %v", err)
	}

	if resolved != target {
		t.Fatalf("unexpected resolved target: got %q want %q", resolved, target)
	}

	designPath := filepath.Join(repo, "docs", "implementation", "01_mini_design.md")
	content, err := os.ReadFile(designPath)
	if err != nil {
		t.Fatalf("read mini design file: %v", err)
	}
	if !strings.Contains(string(content), target) {
		t.Fatalf("mini design file does not mention target repo")
	}

	srcArchPath := filepath.Join(repo, "src", "ARCHITECTURE.md")
	if _, err := os.Stat(srcArchPath); err != nil {
		t.Fatalf("expected src architecture file %s: %v", srcArchPath, err)
	}
	goMiniMain := filepath.Join(repo, "src", "go-mini", "cmd", "mini", "main.go")
	if _, err := os.Stat(goMiniMain); err != nil {
		t.Fatalf("expected go mini scaffold file %s: %v", goMiniMain, err)
	}
}

func TestDeriveMiniRequiresAnalyze(t *testing.T) {
	repo := t.TempDir()
	if _, err := installer.Install(repo); err != nil {
		t.Fatalf("Install returned error: %v", err)
	}

	if _, err := DeriveMini(repo); err == nil {
		t.Fatalf("expected error when analysis metadata is missing")
	}
}

func TestDeriveMiniRequiresChallengeComplete(t *testing.T) {
	repo := t.TempDir()
	if _, err := installer.Install(repo); err != nil {
		t.Fatalf("Install returned error: %v", err)
	}

	target := t.TempDir()
	if _, err := analyzer.Analyze(target, repo); err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}

	if _, err := DeriveMini(repo); err == nil {
		t.Fatalf("expected error when challenge phase is incomplete")
	}
}
