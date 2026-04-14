package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	var out bytes.Buffer

	if err := runWithIO([]string{"help"}, &out, "test"); err != nil {
		t.Fatalf("runWithIO returned error: %v", err)
	}

	if !strings.Contains(out.String(), "d2a init <repo-dir>") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestVersion(t *testing.T) {
	var out bytes.Buffer

	if err := runWithIO([]string{"version"}, &out, "v0.0.1"); err != nil {
		t.Fatalf("runWithIO returned error: %v", err)
	}

	if got := out.String(); got != "v0.0.1\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestInit(t *testing.T) {
	var out bytes.Buffer
	repo := t.TempDir()

	if err := runWithIO([]string{"init", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO returned error: %v", err)
	}

	if !strings.Contains(out.String(), "d2a repo: "+filepath.Base(repo)) {
		t.Fatalf("unexpected output: %q", out.String())
	}
	if !strings.Contains(out.String(), "d2a repo path: "+repo) {
		t.Fatalf("unexpected output: %q", out.String())
	}
	if !strings.Contains(out.String(), "d2a path: "+filepath.Join(repo, ".d2a")) {
		t.Fatalf("unexpected output: %q", out.String())
	}
	if !strings.Contains(out.String(), "initialized d2a in repository") {
		t.Fatalf("unexpected output: %q", out.String())
	}
	for _, rel := range []string{
		filepath.Join(".d2a", "state.json"),
		filepath.Join(".d2a", "history.jsonl"),
	} {
		if _, err := os.Stat(filepath.Join(repo, rel)); err != nil {
			t.Fatalf("expected state artifact %s: %v", rel, err)
		}
	}
}

func TestAnalyze(t *testing.T) {
	var out bytes.Buffer
	repo := t.TempDir()
	target := t.TempDir()

	if err := runWithIO([]string{"init", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO init returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"analyze", target, "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO analyze returned error: %v", err)
	}

	if !strings.Contains(out.String(), "d2a repo: "+filepath.Base(repo)) {
		t.Fatalf("unexpected output: %q", out.String())
	}
	if !strings.Contains(out.String(), "prepared d2a analysis for") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestDeriveMini(t *testing.T) {
	var out bytes.Buffer
	repo := t.TempDir()
	target := t.TempDir()

	if err := runWithIO([]string{"init", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO init returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"analyze", target, "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO analyze returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{
		"skill-state", "d2a-challenge-architecture",
		"--repo", repo,
		"--status", "completed",
		"--stage", "architecture-challenge-complete",
		"--phase", "challenge-dialogue",
		"--question-index", "6",
		"--question-total", "6",
		"--recommendation", "proceed",
		"--summary", "Challenge phase complete.",
	}, &out, "test"); err != nil {
		t.Fatalf("runWithIO skill-state returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"derive-mini", "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO derive-mini returned error: %v", err)
	}

	if !strings.Contains(out.String(), "d2a repo: "+filepath.Base(repo)) {
		t.Fatalf("unexpected output: %q", out.String())
	}
	if !strings.Contains(out.String(), "prepared d2a mini derivation for") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestTestMini(t *testing.T) {
	var out bytes.Buffer
	repo := t.TempDir()
	target := t.TempDir()

	if err := runWithIO([]string{"init", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO init returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"analyze", target, "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO analyze returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{
		"skill-state", "d2a-challenge-architecture",
		"--repo", repo,
		"--status", "completed",
		"--stage", "architecture-challenge-complete",
		"--phase", "challenge-dialogue",
		"--question-index", "6",
		"--question-total", "6",
		"--recommendation", "proceed",
		"--summary", "Challenge phase complete.",
	}, &out, "test"); err != nil {
		t.Fatalf("runWithIO skill-state returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"derive-mini", "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO derive-mini returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"test-mini", "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO test-mini returned error: %v", err)
	}

	if !strings.Contains(out.String(), "d2a repo: "+filepath.Base(repo)) {
		t.Fatalf("unexpected output: %q", out.String())
	}
	if !strings.Contains(out.String(), "prepared d2a test plan for") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestReport(t *testing.T) {
	var out bytes.Buffer
	repo := t.TempDir()
	target := t.TempDir()

	if err := runWithIO([]string{"init", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO init returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"analyze", target, "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO analyze returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{
		"skill-state", "d2a-challenge-architecture",
		"--repo", repo,
		"--status", "completed",
		"--stage", "architecture-challenge-complete",
		"--phase", "challenge-dialogue",
		"--question-index", "6",
		"--question-total", "6",
		"--recommendation", "proceed",
		"--summary", "Challenge phase complete.",
	}, &out, "test"); err != nil {
		t.Fatalf("runWithIO skill-state returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"derive-mini", "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO derive-mini returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"test-mini", "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO test-mini returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"report", "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO report returned error: %v", err)
	}

	if !strings.Contains(out.String(), "d2a repo: "+filepath.Base(repo)) {
		t.Fatalf("unexpected output: %q", out.String())
	}
	if !strings.Contains(out.String(), "prepared d2a report for") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestAnalyzeUsesCurrentDirectoryRepo(t *testing.T) {
	var out bytes.Buffer
	repo := t.TempDir()
	target := t.TempDir()

	if err := runWithIO([]string{"init", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO init returned error: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()
	if err := os.Chdir(repo); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"analyze", target}, &out, "test"); err != nil {
		t.Fatalf("runWithIO analyze returned error: %v", err)
	}

	if !strings.Contains(out.String(), "d2a repo: "+filepath.Base(repo)) {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestAnalyzeFailsWithoutActiveRepo(t *testing.T) {
	var out bytes.Buffer
	target := t.TempDir()
	nonRepo := t.TempDir()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()
	if err := os.Chdir(nonRepo); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}

	err = runWithIO([]string{"analyze", target}, &out, "test")
	if err == nil {
		t.Fatalf("expected error without active repository")
	}
	if !strings.Contains(err.Error(), "no active d2a repository could be determined") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeriveMiniCanSkipChallengeWithReason(t *testing.T) {
	var out bytes.Buffer
	repo := t.TempDir()
	target := t.TempDir()

	if err := runWithIO([]string{"init", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO init returned error: %v", err)
	}
	out.Reset()
	if err := runWithIO([]string{"analyze", target, "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO analyze returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"derive-mini", "--repo", repo, "--skip-challenge-reason", "live demo fast path"}, &out, "test"); err != nil {
		t.Fatalf("runWithIO derive-mini returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"status", "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO status returned error: %v", err)
	}
	text := out.String()
	if !strings.Contains(text, "current stage: mini-derivation-prepared") {
		t.Fatalf("unexpected status output: %q", text)
	}
}

func TestStatus(t *testing.T) {
	var out bytes.Buffer
	repo := t.TempDir()
	target := t.TempDir()

	if err := runWithIO([]string{"init", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO init returned error: %v", err)
	}
	out.Reset()
	if err := runWithIO([]string{"analyze", target, "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO analyze returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"status", "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO status returned error: %v", err)
	}

	text := out.String()
	if !strings.Contains(text, "current stage: analysis-prepared") {
		t.Fatalf("unexpected status output: %q", text)
	}
	if !strings.Contains(text, "question progress: none") {
		t.Fatalf("unexpected status output: %q", text)
	}
	if !strings.Contains(text, "last command: d2a analyze") {
		t.Fatalf("unexpected status output: %q", text)
	}
	if !strings.Contains(text, "recent history:") {
		t.Fatalf("unexpected status output: %q", text)
	}
}

func TestSkillState(t *testing.T) {
	var out bytes.Buffer
	repo := t.TempDir()

	if err := runWithIO([]string{"init", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO init returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{
		"skill-state", "d2a-project-scope",
		"--repo", repo,
		"--status", "progress",
		"--stage", "architecture-in-progress",
		"--phase", "confirmation-questions",
		"--question-index", "2",
		"--question-total", "5",
		"--next-step", "Continue confirmation questions.",
		"--next-skill", "d2a-runtime-view",
		"--next-file", ".d2a/docs/architecture/02_driver.md",
		"--summary", "Question loop is in progress.",
	}, &out, "test"); err != nil {
		t.Fatalf("runWithIO skill-state returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"status", "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO status returned error: %v", err)
	}

	text := out.String()
	if !strings.Contains(text, "current skill: d2a-project-scope") {
		t.Fatalf("unexpected status output: %q", text)
	}
	if !strings.Contains(text, "current phase: confirmation-questions") {
		t.Fatalf("unexpected status output: %q", text)
	}
	if !strings.Contains(text, "question progress: 2/5") {
		t.Fatalf("unexpected status output: %q", text)
	}
	if !strings.Contains(text, "next skill: d2a-runtime-view") {
		t.Fatalf("unexpected status output: %q", text)
	}
}

func TestChallengeSkillState(t *testing.T) {
	var out bytes.Buffer
	repo := t.TempDir()

	if err := runWithIO([]string{"init", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO init returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{
		"skill-state", "d2a-challenge-architecture",
		"--repo", repo,
		"--status", "progress",
		"--stage", "architecture-challenge-in-progress",
		"--phase", "challenge-dialogue",
		"--question-index", "3",
		"--question-total", "6",
		"--decision", "primary driver",
		"--strength", "partial",
		"--recommendation", "review",
		"--objection", "Why not a simpler timer-based trigger?",
		"--next-step", "Continue architecture challenge dialogue.",
		"--next-skill", "d2a-mini-scope",
		"--next-file", ".d2a/docs/implementation/00_mini_scope.md",
		"--summary", "Challenge round 3 is active.",
	}, &out, "test"); err != nil {
		t.Fatalf("runWithIO skill-state returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"status", "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO status returned error: %v", err)
	}

	text := out.String()
	if !strings.Contains(text, "challenge progress: 3/6") {
		t.Fatalf("unexpected status output: %q", text)
	}
	if !strings.Contains(text, "current decision: primary driver") {
		t.Fatalf("unexpected status output: %q", text)
	}
	if !strings.Contains(text, "last challenge strength: partial") {
		t.Fatalf("unexpected status output: %q", text)
	}
	if !strings.Contains(text, "challenge recommendation: review") {
		t.Fatalf("unexpected status output: %q", text)
	}
}
