package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstall(t *testing.T) {
	target := t.TempDir()

	installed, err := Install(target)
	if err != nil {
		t.Fatalf("Install returned error: %v", err)
	}

	if installed != target {
		t.Fatalf("unexpected install target: got %q want %q", installed, target)
	}

	skills, err := SkillNames()
	if err != nil {
		t.Fatalf("SkillNames returned error: %v", err)
	}

	for _, toolDir := range ToolDirs() {
		for _, skill := range skills {
			path := filepath.Join(target, toolDir, "skills", skill, "SKILL.md")
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("expected installed skill file %s: %v", path, err)
			}
		}
	}

	d2aFiles := []string{
		filepath.Join(".d2a", "LAB.md"),
		filepath.Join(".d2a", "docs", "architecture", "00_overview.md"),
		filepath.Join(".d2a", "docs", "architecture", "01_boundary.md"),
		filepath.Join(".d2a", "docs", "architecture", "02_driver.md"),
		filepath.Join(".d2a", "docs", "architecture", "03_core_objects.md"),
		filepath.Join(".d2a", "docs", "architecture", "04_state_evolution.md"),
		filepath.Join(".d2a", "docs", "architecture", "05_cooperation.md"),
		filepath.Join(".d2a", "docs", "architecture", "06_constraints.md"),
		filepath.Join(".d2a", "docs", "architecture", "99_code_map.md"),
		filepath.Join(".d2a", "docs", "implementation", "00_mini_scope.md"),
		filepath.Join(".d2a", "docs", "implementation", "01_mini_design.md"),
		filepath.Join(".d2a", "docs", "implementation", "02_build_plan.md"),
		filepath.Join(".d2a", "docs", "implementation", "03_test_plan.md"),
		filepath.Join(".d2a", "docs", "report", "00_report_outline.md"),
		filepath.Join(".d2a", "src", "README.md"),
		filepath.Join(".d2a", "tests", "README.md"),
		filepath.Join(".d2a", "report", "README.md"),
		filepath.Join(".d2a", "report", "vue-app", "package.json"),
		filepath.Join(".d2a", "report", "vue-app", "src", "App.vue"),
	}

	for _, rel := range d2aFiles {
		path := filepath.Join(target, rel)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected d2a file %s: %v", path, err)
		}
	}

	gitignorePath := filepath.Join(target, ".gitignore")
	gitignore, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("read .gitignore: %v", err)
	}
	for _, toolDir := range ToolDirs() {
		entry := toolDir + "/"
		if !strings.Contains(string(gitignore), entry) {
			t.Fatalf("expected .gitignore to contain %q", entry)
		}
	}
}

func TestInstallUpdatesExistingGitignoreWithoutDuplicateEntries(t *testing.T) {
	target := t.TempDir()
	gitignorePath := filepath.Join(target, ".gitignore")
	initial := "node_modules/\n.codex/\n"
	if err := os.WriteFile(gitignorePath, []byte(initial), 0o644); err != nil {
		t.Fatalf("write initial .gitignore: %v", err)
	}

	if _, err := Install(target); err != nil {
		t.Fatalf("Install returned error: %v", err)
	}
	if _, err := Install(target); err != nil {
		t.Fatalf("second Install returned error: %v", err)
	}

	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("read .gitignore: %v", err)
	}
	text := string(content)
	if strings.Count(text, ".codex/") != 1 {
		t.Fatalf("expected .codex/ to appear once, got %d", strings.Count(text, ".codex/"))
	}
	for _, toolDir := range ToolDirs() {
		entry := toolDir + "/"
		if !strings.Contains(text, entry) {
			t.Fatalf("expected .gitignore to contain %q", entry)
		}
	}
}
