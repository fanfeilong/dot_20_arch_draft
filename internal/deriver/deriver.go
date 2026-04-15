package deriver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fanfeilong/dot_20_arch_draft/internal/state"
)

type metadata struct {
	TargetRepo string `json:"target_repo"`
	RepoRoot   string `json:"repo_root"`
	D2APath    string `json:"d2a_path"`
}

func DeriveMini(repoRoot string) (string, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", fmt.Errorf("resolve repo root: %w", err)
	}

	meta, err := loadMetadata(repoRoot)
	if err != nil {
		return "", err
	}
	if err := ensureChallengeComplete(repoRoot); err != nil {
		return "", err
	}

	for rel, content := range implementationTasks(meta.TargetRepo) {
		path := filepath.Join(repoRoot, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return "", fmt.Errorf("create implementation dir %s: %w", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return "", fmt.Errorf("write implementation file %s: %w", path, err)
		}
	}

	srcArchPath := filepath.Join(repoRoot, "src", "ARCHITECTURE.md")
	if err := os.MkdirAll(filepath.Dir(srcArchPath), 0o755); err != nil {
		return "", fmt.Errorf("create src dir %s: %w", filepath.Dir(srcArchPath), err)
	}
	if err := os.WriteFile(srcArchPath, []byte(srcArchitecture(meta.TargetRepo)), 0o644); err != nil {
		return "", fmt.Errorf("write src architecture file %s: %w", srcArchPath, err)
	}
	if err := writeGoMiniScaffold(repoRoot, meta.TargetRepo); err != nil {
		return "", err
	}

	return meta.TargetRepo, nil
}

func ensureChallengeComplete(repoRoot string) error {
	s, err := state.Load(repoRoot)
	if err != nil {
		return fmt.Errorf("load d2a state: %w", err)
	}
	if s.CurrentStage != state.StageArchitectureChallengeComplete &&
		s.CurrentStage != state.StageMiniDerivationPrepared &&
		s.CurrentStage != state.StageTestPlanPrepared &&
		s.CurrentStage != state.StageReportReady &&
		s.CurrentStage != state.StageServing {
		return fmt.Errorf("architecture challenge phase is not complete: %s", s.CurrentStage)
	}
	return nil
}

func loadMetadata(repoRoot string) (metadata, error) {
	d2aFile := filepath.Join(repoRoot, "LAB.md")
	if _, err := os.Stat(d2aFile); err != nil {
		if os.IsNotExist(err) {
			return metadata{}, fmt.Errorf("d2a not initialized in repository: %s", repoRoot)
		}
		return metadata{}, fmt.Errorf("check d2a file %s: %w", d2aFile, err)
	}

	metaPath := filepath.Join(repoRoot, ".d2a", "target.json")
	content, err := os.ReadFile(metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			return metadata{}, fmt.Errorf("analysis not prepared: %s", repoRoot)
		}
		return metadata{}, fmt.Errorf("read metadata file %s: %w", metaPath, err)
	}

	var meta metadata
	if err := json.Unmarshal(content, &meta); err != nil {
		return metadata{}, fmt.Errorf("parse metadata file %s: %w", metaPath, err)
	}
	if meta.TargetRepo == "" {
		return metadata{}, fmt.Errorf("analysis metadata missing target repo: %s", metaPath)
	}

	return meta, nil
}

func implementationTasks(target string) map[string]string {
	return map[string]string{
		"docs/implementation/00_mini_scope.md":  miniScopeTask(target),
		"docs/implementation/01_mini_design.md": miniDesignTask(target),
		"docs/implementation/02_build_plan.md":  buildPlanTask(target),
		"docs/implementation/03_test_plan.md":   testPlanTask(target),
	}
}

func sharedHeader(title, target string) string {
	return fmt.Sprintf(`# %s

## Target

- Target repo: %s

## Output Rules

- Keep the result small and implementation-oriented.
- Preserve only the architecture idea that matters most.
- Prefer a runnable 20 percent slice over broad feature coverage.
- Keep the target stack aligned with the original project when practical.

`, title, target)
}

func miniScopeTask(target string) string {
	return sharedHeader("00. Mini Scope", target) + `## Inputs

- docs/architecture/00_overview.md
- docs/architecture/02_driver.md
- docs/architecture/03_core_objects.md
- docs/architecture/04_state_evolution.md
- docs/architecture/05_cooperation.md
- docs/architecture/06_constraints.md

## Atomic Questions

1. What single architecture idea should the mini version preserve?
2. Which 20 percent slice is enough to demonstrate that idea?
3. What will the mini version intentionally omit?
4. Which stack should be used to stay close to the target project?
`
}

func miniDesignTask(target string) string {
	return sharedHeader("01. Mini Design", target) + `## Inputs

- docs/implementation/00_mini_scope.md
- docs/architecture/02_driver.md
- docs/architecture/03_core_objects.md
- docs/architecture/04_state_evolution.md
- docs/architecture/05_cooperation.md

## Atomic Questions

1. What are the main modules of the mini version?
2. What interfaces or entry points are required?
3. What is the runtime flow of the mini version?
4. What state model must be preserved?
`
}

func buildPlanTask(target string) string {
	return sharedHeader("02. Build Plan", target) + `## Inputs

- docs/implementation/00_mini_scope.md
- docs/implementation/01_mini_design.md

## Atomic Questions

1. Which built-in provider (if any) matches the target stack, and what minimal scaffold does it imply?
2. What is the first runnable slice that can be completed within a strict timebox?
3. How does this slice prove the architecture intent anchors (object/state/cooperation chain)?
4. What is intentionally left unimplemented to protect the timebox?
`
}

func testPlanTask(target string) string {
	return sharedHeader("03. Test Plan", target) + `## Inputs

- docs/implementation/00_mini_scope.md
- docs/implementation/01_mini_design.md
- docs/implementation/02_build_plan.md

## Atomic Questions

1. What is the first end-to-end test for the mini version?
2. What main scenarios should be covered next?
3. What outputs should be observable in tests?
4. Which failure cases are worth keeping?
`
}

func srcArchitecture(target string) string {
	return fmt.Sprintf(`# Mini Architecture

## Target

- Target repo: %s

## Purpose

This file should summarize the chosen mini implementation after docs/implementation/ is filled.

## To Fill

- The architecture idea being preserved
- The chosen 20 percent runnable slice
- The modules to implement under src/
- The tests that must pass first
`, target)
}

func writeGoMiniScaffold(repoRoot, target string) error {
	files := map[string]string{
		"src/go-mini/README.md": goMiniReadme(target),
		"src/go-mini/go.mod": `module d2a-mini

go 1.22
`,
		"src/go-mini/cmd/mini/main.go": `package main

import (
	"fmt"

	"d2a-mini/internal/mini"
)

func main() {
	runner := mini.NewRunner()
	out := runner.Run("demo-input")
	fmt.Println(out)
}
`,
		"src/go-mini/internal/mini/runner.go": `package mini

type Runner struct{}

func NewRunner() Runner {
	return Runner{}
}

func (Runner) Run(input string) string {
	if input == "" {
		return "mini-failure: empty input"
	}
	return "mini-success: " + input
}
`,
	}

	for rel, content := range files {
		path := filepath.Join(repoRoot, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("create scaffold dir %s: %w", filepath.Dir(path), err)
		}
		if _, err := os.Stat(path); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("check scaffold file %s: %w", path, err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return fmt.Errorf("write scaffold file %s: %w", path, err)
		}
	}

	return nil
}

func goMiniReadme(target string) string {
	return fmt.Sprintf(`# Go Mini Scaffold

This is the first supported runnable mini stack scaffold (Go).

- target repo: %s
- purpose: provide an executable baseline that skills can evolve

Run:

`+"```bash\ncd src/go-mini\ngo run ./cmd/mini\n```"+`

Expected output:

`+"```text\nmini-success: demo-input\n```"+`
`, target)
}
