package analyzer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fanfeilong/dot_20_arch_draft/internal/installer"
)

func TestAnalyze(t *testing.T) {
	repo := t.TempDir()
	if _, err := installer.Install(repo); err != nil {
		t.Fatalf("Install returned error: %v", err)
	}

	target := t.TempDir()
	resolved, err := Analyze(target, repo)
	if err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}

	if resolved != target {
		t.Fatalf("unexpected resolved target: got %q want %q", resolved, target)
	}

	metaPath := filepath.Join(repo, ".d2a", "target.json")
	if _, err := os.Stat(metaPath); err != nil {
		t.Fatalf("expected metadata file %s: %v", metaPath, err)
	}

	overviewPath := filepath.Join(repo, "docs", "architecture", "00_overview.md")
	content, err := os.ReadFile(overviewPath)
	if err != nil {
		t.Fatalf("read overview file: %v", err)
	}

	text := string(content)
	if !strings.Contains(text, target) {
		t.Fatalf("overview file does not mention target repo: %q", text)
	}
	if !strings.Contains(text, "d2a-step") {
		t.Fatalf("overview file does not mention primary skill: %q", text)
	}
}

func TestAnalyzeRequiresInitializedRepo(t *testing.T) {
	repo := t.TempDir()

	if _, err := Analyze("example/repo", repo); err == nil {
		t.Fatalf("expected error for uninitialized repository")
	}
}
