package installer

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fanfeilong/dot_20_arch_draft/assets"
)

var toolDirs = []string{
	".codex",
	".claude",
	".cursor",
	".opencode",
	".trae",
	".neocode",
}

const (
	LanguageZH = "zh"
	LanguageEN = "en"
)

func Install(targetDir string) (string, error) {
	return InstallWithLanguage(targetDir, LanguageZH)
}

func InstallWithLanguage(targetDir, language string) (string, error) {
	targetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return "", fmt.Errorf("resolve target dir: %w", err)
	}
	language, err = normalizeLanguage(language)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", fmt.Errorf("create target dir: %w", err)
	}

	skills, err := loadSkills(language)
	if err != nil {
		return "", err
	}

	if err := installLab(targetDir, language); err != nil {
		return "", err
	}
	if language == LanguageZH {
		if err := ensureZhLegacyDocAliases(targetDir); err != nil {
			return "", err
		}
	}
	if err := writeLanguageConfig(targetDir, language); err != nil {
		return "", err
	}
	if err := ensureGitignoreEntries(targetDir); err != nil {
		return "", err
	}

	for _, toolDir := range toolDirs {
		base := filepath.Join(targetDir, toolDir, "skills")
		if err := os.MkdirAll(base, 0o755); err != nil {
			return "", fmt.Errorf("create skills dir %s: %w", base, err)
		}

		for name, content := range skills {
			skillDir := filepath.Join(base, name)
			if err := os.MkdirAll(skillDir, 0o755); err != nil {
				return "", fmt.Errorf("create skill dir %s: %w", skillDir, err)
			}

			skillFile := filepath.Join(skillDir, "SKILL.md")
			if err := os.WriteFile(skillFile, content, 0o644); err != nil {
				return "", fmt.Errorf("write skill file %s: %w", skillFile, err)
			}
		}
	}

	return targetDir, nil
}

func ToolDirs() []string {
	out := make([]string, len(toolDirs))
	copy(out, toolDirs)
	return out
}

func SkillNames() ([]string, error) {
	skills, err := loadSkills(LanguageZH)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(skills))
	for name := range skills {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

func loadSkills(language string) (map[string][]byte, error) {
	root, err := skillsRoot(language)
	if err != nil {
		return nil, err
	}
	entries, err := fs.ReadDir(assets.Files, root)
	if err != nil {
		return nil, fmt.Errorf("read embedded skills: %w", err)
	}

	skills := make(map[string][]byte, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		path := filepath.ToSlash(filepath.Join(root, name, "SKILL.md"))
		content, err := fs.ReadFile(assets.Files, path)
		if err != nil {
			return nil, fmt.Errorf("read embedded skill %s: %w", name, err)
		}

		skills[name] = content
	}

	return skills, nil
}

func skillsRoot(language string) (string, error) {
	switch language {
	case LanguageZH:
		return "skills_zh_cn", nil
	case LanguageEN:
		return "skills_en", nil
	default:
		return "", fmt.Errorf("unsupported language pack: %s", language)
	}
}

func installLab(targetDir, language string) error {
	root, err := labRoot(language)
	if err != nil {
		return err
	}
	return fs.WalkDir(assets.Files, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk embedded d2a assets: %w", err)
		}
		if path == root {
			return nil
		}

		rel := strings.TrimPrefix(path, root+"/")
		targetPath := filepath.Join(targetDir, filepath.FromSlash(rel))

		if d.IsDir() {
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return fmt.Errorf("create lab dir %s: %w", targetPath, err)
			}
			return nil
		}

		content, err := fs.ReadFile(assets.Files, path)
		if err != nil {
			return fmt.Errorf("read embedded lab file %s: %w", path, err)
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return fmt.Errorf("create lab parent dir %s: %w", filepath.Dir(targetPath), err)
		}
		if err := os.WriteFile(targetPath, content, 0o644); err != nil {
			return fmt.Errorf("write lab file %s: %w", targetPath, err)
		}

		return nil
	})
}

func labRoot(language string) (string, error) {
	switch language {
	case LanguageZH:
		return "lab_zh_cn", nil
	case LanguageEN:
		return "lab_en", nil
	default:
		return "", fmt.Errorf("unsupported language pack: %s", language)
	}
}

func ensureGitignoreEntries(targetDir string) error {
	path := filepath.Join(targetDir, ".gitignore")
	content, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read gitignore %s: %w", path, err)
	}

	text := string(content)
	lines := strings.Split(text, "\n")
	toAdd := make([]string, 0, len(toolDirs))
	for _, dir := range toolDirs {
		name := strings.TrimPrefix(dir, ".")
		if gitignoreHas(lines, name) {
			continue
		}
		toAdd = append(toAdd, dir+"/")
	}
	if !gitignoreHas(lines, "repos") {
		toAdd = append(toAdd, "repos/")
	}
	if len(toAdd) == 0 {
		return nil
	}

	var b strings.Builder
	b.WriteString(text)
	if b.Len() > 0 && !strings.HasSuffix(b.String(), "\n") {
		b.WriteString("\n")
	}
	b.WriteString("\n# d2a skills directories\n")
	for _, entry := range toAdd {
		b.WriteString(entry)
		b.WriteString("\n")
	}

	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		return fmt.Errorf("write gitignore %s: %w", path, err)
	}
	return nil
}

func ensureZhLegacyDocAliases(targetDir string) error {
	type fileAlias struct {
		srcRel string
		dstRel string
	}
	aliases := []fileAlias{
		{srcRel: filepath.Join("docs", "1.架构拆解", "00_总览.md"), dstRel: filepath.Join("docs", "architecture", "00_overview.md")},
		{srcRel: filepath.Join("docs", "1.架构拆解", "01_边界.md"), dstRel: filepath.Join("docs", "architecture", "01_boundary.md")},
		{srcRel: filepath.Join("docs", "1.架构拆解", "02_驱动.md"), dstRel: filepath.Join("docs", "architecture", "02_driver.md")},
		{srcRel: filepath.Join("docs", "1.架构拆解", "03_核心对象.md"), dstRel: filepath.Join("docs", "architecture", "03_core_objects.md")},
		{srcRel: filepath.Join("docs", "1.架构拆解", "04_状态演化.md"), dstRel: filepath.Join("docs", "architecture", "04_state_evolution.md")},
		{srcRel: filepath.Join("docs", "1.架构拆解", "05_协作.md"), dstRel: filepath.Join("docs", "architecture", "05_cooperation.md")},
		{srcRel: filepath.Join("docs", "1.架构拆解", "06_约束.md"), dstRel: filepath.Join("docs", "architecture", "06_constraints.md")},
		{srcRel: filepath.Join("docs", "1.架构拆解", "99_代码地图.md"), dstRel: filepath.Join("docs", "architecture", "99_code_map.md")},
		{srcRel: filepath.Join("docs", "2.mini实现", "00_最小范围.md"), dstRel: filepath.Join("docs", "implementation", "00_mini_scope.md")},
		{srcRel: filepath.Join("docs", "2.mini实现", "01_最小设计.md"), dstRel: filepath.Join("docs", "implementation", "01_mini_design.md")},
		{srcRel: filepath.Join("docs", "2.mini实现", "02_构建计划.md"), dstRel: filepath.Join("docs", "implementation", "02_build_plan.md")},
		{srcRel: filepath.Join("docs", "2.mini实现", "03_测试计划.md"), dstRel: filepath.Join("docs", "implementation", "03_test_plan.md")},
		{srcRel: filepath.Join("docs", "3.报告", "00_报告大纲.md"), dstRel: filepath.Join("docs", "report", "00_report_outline.md")},
	}
	for _, a := range aliases {
		src := filepath.Join(targetDir, a.srcRel)
		dst := filepath.Join(targetDir, a.dstRel)
		if err := copyMissingFile(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func copyMissingFile(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read source file %s: %w", src, err)
	}
	if _, err := os.Stat(dst); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat target file %s: %w", dst, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("create target dir %s: %w", filepath.Dir(dst), err)
	}
	if err := os.WriteFile(dst, content, 0o644); err != nil {
		return fmt.Errorf("write target file %s: %w", dst, err)
	}
	return nil
}

func gitignoreHas(lines []string, name string) bool {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		normalized := strings.TrimPrefix(trimmed, "/")
		normalized = strings.TrimSuffix(normalized, "/")
		normalized = strings.TrimPrefix(normalized, ".")
		if normalized == name {
			return true
		}
	}
	return false
}

func normalizeLanguage(language string) (string, error) {
	value := strings.ToLower(strings.TrimSpace(language))
	switch value {
	case "", "zh", "zh-cn":
		return LanguageZH, nil
	case "en", "en-us":
		return LanguageEN, nil
	default:
		return "", fmt.Errorf("unsupported language pack %q; use zh or en", language)
	}
}

func writeLanguageConfig(targetDir, language string) error {
	path := filepath.Join(targetDir, ".d2a", "language.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create language config dir %s: %w", filepath.Dir(path), err)
	}
	content, err := json.MarshalIndent(map[string]string{
		"language": language,
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal language config: %w", err)
	}
	content = append(content, '\n')
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write language config %s: %w", path, err)
	}
	return nil
}
