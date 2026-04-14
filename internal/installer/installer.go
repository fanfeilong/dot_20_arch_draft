package installer

import (
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

func Install(targetDir string) (string, error) {
	targetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return "", fmt.Errorf("resolve target dir: %w", err)
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", fmt.Errorf("create target dir: %w", err)
	}

	skills, err := loadSkills()
	if err != nil {
		return "", err
	}

	if err := installLab(targetDir); err != nil {
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
	skills, err := loadSkills()
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

func loadSkills() (map[string][]byte, error) {
	entries, err := fs.ReadDir(assets.Files, "skills")
	if err != nil {
		return nil, fmt.Errorf("read embedded skills: %w", err)
	}

	skills := make(map[string][]byte, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		path := filepath.ToSlash(filepath.Join("skills", name, "SKILL.md"))
		content, err := fs.ReadFile(assets.Files, path)
		if err != nil {
			return nil, fmt.Errorf("read embedded skill %s: %w", name, err)
		}

		skills[name] = content
	}

	return skills, nil
}

func installLab(targetDir string) error {
	return fs.WalkDir(assets.Files, "lab", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk embedded d2a assets: %w", err)
		}
		if path == "lab" {
			return nil
		}

		rel := strings.TrimPrefix(path, "lab/")
		targetPath := filepath.Join(targetDir, ".d2a", filepath.FromSlash(rel))

		if d.IsDir() {
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return fmt.Errorf("create d2a dir %s: %w", targetPath, err)
			}
			return nil
		}

		content, err := fs.ReadFile(assets.Files, path)
		if err != nil {
			return fmt.Errorf("read embedded d2a file %s: %w", path, err)
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return fmt.Errorf("create d2a parent dir %s: %w", filepath.Dir(targetPath), err)
		}
		if err := os.WriteFile(targetPath, content, 0o644); err != nil {
			return fmt.Errorf("write d2a file %s: %w", targetPath, err)
		}

		return nil
	})
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
