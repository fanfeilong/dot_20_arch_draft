package installer

import (
	"os"
	"path/filepath"
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
}
