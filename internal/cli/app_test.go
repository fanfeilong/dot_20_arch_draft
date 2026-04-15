package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	var out bytes.Buffer

	if err := runWithIO([]string{"help"}, &out, "test"); err != nil {
		t.Fatalf("runWithIO returned error: %v", err)
	}

	if !strings.Contains(out.String(), "d2a init <target-repo-git-url>") {
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
	fx := initWorkspace(t, &out, "zh")
	repo := fx.workspace

	if !strings.Contains(out.String(), "d2a repo: "+filepath.Base(repo)) {
		t.Fatalf("unexpected output: %q", out.String())
	}
	if !strings.Contains(out.String(), "d2a repo path: "+repo) {
		t.Fatalf("unexpected output: %q", out.String())
	}
	if !strings.Contains(out.String(), "d2a path: "+filepath.Join(repo, ".d2a")) {
		t.Fatalf("unexpected output: %q", out.String())
	}
	if !strings.Contains(out.String(), "d2a language: zh") {
		t.Fatalf("unexpected output: %q", out.String())
	}
	if !strings.Contains(out.String(), "initialized d2a workspace") {
		t.Fatalf("unexpected output: %q", out.String())
	}
	for _, rel := range []string{
		filepath.Join(".d2a", "state.json"),
		filepath.Join(".d2a", "history.jsonl"),
		filepath.Join("AGENTS.md"),
	} {
		if _, err := os.Stat(filepath.Join(repo, rel)); err != nil {
			t.Fatalf("expected state artifact %s: %v", rel, err)
		}
	}
	if _, err := os.Stat(fx.targetPath); err != nil {
		t.Fatalf("expected cloned target repo %s: %v", fx.targetPath, err)
	}
}

func TestInitWithEnglishLanguagePack(t *testing.T) {
	var out bytes.Buffer
	fx := initWorkspace(t, &out, "en")
	repo := fx.workspace
	if !strings.Contains(out.String(), "d2a language: en") {
		t.Fatalf("unexpected output: %q", out.String())
	}

	skillPath := filepath.Join(repo, ".codex", "skills", "d2a-step", "SKILL.md")
	content, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("read skill file: %v", err)
	}
	if !strings.Contains(string(content), "All user-facing text must be in English.") {
		t.Fatalf("expected english skill language rule")
	}
}

func TestInitRejectsInvalidLanguagePack(t *testing.T) {
	var out bytes.Buffer
	targetURL, _ := createGitRemoteWithCommit(t)

	err := runInitInTempWorkspace(t, []string{"init", targetURL, "--lang", "jp"}, &out)
	if err == nil {
		t.Fatalf("expected error for invalid language pack")
	}
	if !strings.Contains(err.Error(), "unsupported language pack") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnalyze(t *testing.T) {
	var out bytes.Buffer
	fx := initWorkspace(t, &out, "zh")
	repo := fx.workspace
	target := fx.targetPath

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
	fx := initWorkspace(t, &out, "zh")
	repo := fx.workspace
	target := fx.targetPath

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
	fx := initWorkspace(t, &out, "zh")
	repo := fx.workspace
	target := fx.targetPath

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
	fx := initWorkspace(t, &out, "zh")
	repo := fx.workspace
	target := fx.targetPath

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
	fx := initWorkspace(t, &out, "zh")
	repo := fx.workspace

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
	if err := runWithIO([]string{"analyze"}, &out, "test"); err != nil {
		t.Fatalf("runWithIO analyze returned error: %v", err)
	}

	if !strings.Contains(out.String(), "d2a repo: "+filepath.Base(repo)) {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestAnalyzeFailsWithoutActiveRepo(t *testing.T) {
	var out bytes.Buffer
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

	err = runWithIO([]string{"analyze"}, &out, "test")
	if err == nil {
		t.Fatalf("expected error without active repository")
	}
	if !strings.Contains(err.Error(), "no active d2a workspace could be determined") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeriveMiniCanSkipChallengeWithReason(t *testing.T) {
	var out bytes.Buffer
	fx := initWorkspace(t, &out, "zh")
	repo := fx.workspace
	target := fx.targetPath
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
	fx := initWorkspace(t, &out, "zh")
	repo := fx.workspace
	target := fx.targetPath
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
	fx := initWorkspace(t, &out, "zh")
	repo := fx.workspace

	out.Reset()
	if err := runWithIO([]string{
		"skill-state", "d2a-arch-1-project-scope",
		"--repo", repo,
		"--status", "progress",
		"--stage", "architecture-in-progress",
		"--phase", "confirmation-questions",
		"--question-index", "2",
		"--question-total", "5",
		"--next-step", "Continue confirmation questions.",
		"--next-skill", "d2a-arch-2-runtime-view",
		"--next-file", "docs/architecture/02_driver.md",
		"--summary", "Question loop is in progress.",
	}, &out, "test"); err != nil {
		t.Fatalf("runWithIO skill-state returned error: %v", err)
	}

	out.Reset()
	if err := runWithIO([]string{"status", "--repo", repo}, &out, "test"); err != nil {
		t.Fatalf("runWithIO status returned error: %v", err)
	}

	text := out.String()
	if !strings.Contains(text, "current skill: d2a-arch-1-project-scope") {
		t.Fatalf("unexpected status output: %q", text)
	}
	if !strings.Contains(text, "current phase: confirmation-questions") {
		t.Fatalf("unexpected status output: %q", text)
	}
	if !strings.Contains(text, "question progress: 2/5") {
		t.Fatalf("unexpected status output: %q", text)
	}
	if !strings.Contains(text, "next skill: d2a-arch-2-runtime-view") {
		t.Fatalf("unexpected status output: %q", text)
	}
}

func TestChallengeSkillState(t *testing.T) {
	var out bytes.Buffer
	fx := initWorkspace(t, &out, "zh")
	repo := fx.workspace

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
		"--next-skill", "d2a-mini-1-scope",
		"--next-file", "docs/implementation/00_mini_scope.md",
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

type initFixture struct {
	workspace  string
	targetURL  string
	targetPath string
}

func initWorkspace(t *testing.T, out *bytes.Buffer, language string) initFixture {
	t.Helper()
	targetURL, repoName := createGitRemoteWithCommit(t)
	workspaceBase := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	if err := os.Chdir(workspaceBase); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()

	args := []string{"init", targetURL}
	if language != "" && language != "zh" {
		args = append(args, "--lang", language)
	}
	if err := runWithIO(args, out, "test"); err != nil {
		t.Fatalf("runWithIO init returned error: %v", err)
	}

	workspace, err := filepath.Abs(repoName + "_d2a")
	if err != nil {
		t.Fatalf("Abs returned error: %v", err)
	}
	return initFixture{
		workspace:  workspace,
		targetURL:  targetURL,
		targetPath: filepath.Join(workspace, "repos", repoName),
	}
}

func runInitInTempWorkspace(t *testing.T, args []string, out *bytes.Buffer) error {
	t.Helper()
	workspaceBase := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	if err := os.Chdir(workspaceBase); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()
	return runWithIO(args, out, "test")
}

func createGitRemoteWithCommit(t *testing.T) (string, string) {
	t.Helper()
	base := t.TempDir()
	remote := filepath.Join(base, "sample.git")
	work := filepath.Join(base, "work")

	runCmd(t, "", "git", "init", "--bare", remote)
	runCmd(t, "", "git", "clone", remote, work)
	if err := os.WriteFile(filepath.Join(work, "README.md"), []byte("# sample\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	runCmd(t, work, "git", "add", "README.md")
	runCmd(t, work, "git", "-c", "user.name=test", "-c", "user.email=test@example.com", "commit", "-m", "init")
	runCmd(t, work, "git", "push", "origin", "HEAD")

	return "file://" + remote, "sample"
}

func runCmd(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %s %s: %v\n%s", name, strings.Join(args, " "), err, string(output))
	}
}
