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
		filepath.Join("LAB.md"),
		filepath.Join("docs", "1.架构拆解", "00_总览.md"),
		filepath.Join("docs", "1.架构拆解", "01_边界.md"),
		filepath.Join("docs", "1.架构拆解", "02_驱动.md"),
		filepath.Join("docs", "1.架构拆解", "03_核心对象.md"),
		filepath.Join("docs", "1.架构拆解", "04_状态演化.md"),
		filepath.Join("docs", "1.架构拆解", "05_协作.md"),
		filepath.Join("docs", "1.架构拆解", "06_约束.md"),
		filepath.Join("docs", "1.架构拆解", "99_代码地图.md"),
		filepath.Join("docs", "2.mini实现", "00_最小范围.md"),
		filepath.Join("docs", "2.mini实现", "01_最小设计.md"),
		filepath.Join("docs", "2.mini实现", "02_构建计划.md"),
		filepath.Join("docs", "2.mini实现", "03_测试计划.md"),
		filepath.Join("docs", "3.报告", "00_报告大纲.md"),
		filepath.Join("docs", "architecture", "00_overview.md"),
		filepath.Join("docs", "architecture", "01_boundary.md"),
		filepath.Join("docs", "architecture", "02_driver.md"),
		filepath.Join("docs", "architecture", "03_core_objects.md"),
		filepath.Join("docs", "architecture", "04_state_evolution.md"),
		filepath.Join("docs", "architecture", "05_cooperation.md"),
		filepath.Join("docs", "architecture", "06_constraints.md"),
		filepath.Join("docs", "architecture", "99_code_map.md"),
		filepath.Join("docs", "implementation", "00_mini_scope.md"),
		filepath.Join("docs", "implementation", "01_mini_design.md"),
		filepath.Join("docs", "implementation", "02_build_plan.md"),
		filepath.Join("docs", "implementation", "03_test_plan.md"),
		filepath.Join("docs", "report", "00_report_outline.md"),
		filepath.Join("src", "README.md"),
		filepath.Join("tests", "README.md"),
		filepath.Join("report", "README.md"),
		filepath.Join("report", "vue-app", "package.json"),
		filepath.Join("report", "vue-app", "src", "App.vue"),
	}

	for _, rel := range d2aFiles {
		path := filepath.Join(target, rel)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected d2a file %s: %v", path, err)
		}
	}
	labPath := filepath.Join(target, "LAB.md")
	labContent, err := os.ReadFile(labPath)
	if err != nil {
		t.Fatalf("read zh lab file: %v", err)
	}
	if !strings.Contains(string(labContent), "本地架构实验空间") {
		t.Fatalf("expected zh lab template content, got: %s", string(labContent))
	}
	languagePath := filepath.Join(target, ".d2a", "language.json")
	content, err := os.ReadFile(languagePath)
	if err != nil {
		t.Fatalf("read language config: %v", err)
	}
	if !strings.Contains(string(content), `"language": "zh"`) {
		t.Fatalf("unexpected language config: %s", string(content))
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

func TestInstallWithEnglishLanguagePack(t *testing.T) {
	target := t.TempDir()
	if _, err := InstallWithLanguage(target, "en"); err != nil {
		t.Fatalf("InstallWithLanguage returned error: %v", err)
	}

	skill := filepath.Join(target, ".codex", "skills", "d2a-step", "SKILL.md")
	content, err := os.ReadFile(skill)
	if err != nil {
		t.Fatalf("read installed en skill: %v", err)
	}
	if !strings.Contains(string(content), "All user-facing text must be in English.") {
		t.Fatalf("expected english language rule in d2a-step skill")
	}
	if strings.Contains(string(content), "Simplified Chinese") {
		t.Fatalf("unexpected chinese language rule in english skill")
	}
	labPath := filepath.Join(target, "LAB.md")
	labContent, err := os.ReadFile(labPath)
	if err != nil {
		t.Fatalf("read en lab file: %v", err)
	}
	if !strings.Contains(string(labContent), "local architecture lab") {
		t.Fatalf("expected english lab template content, got: %s", string(labContent))
	}

	languagePath := filepath.Join(target, ".d2a", "language.json")
	langContent, err := os.ReadFile(languagePath)
	if err != nil {
		t.Fatalf("read language config: %v", err)
	}
	if !strings.Contains(string(langContent), `"language": "en"`) {
		t.Fatalf("unexpected language config: %s", string(langContent))
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
