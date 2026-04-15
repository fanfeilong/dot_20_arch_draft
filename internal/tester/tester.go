package tester

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type metadata struct {
	TargetRepo string `json:"target_repo"`
	RepoRoot   string `json:"repo_root"`
	D2APath    string `json:"d2a_path"`
}

type manifest struct {
	TargetRepo string   `json:"target_repo"`
	RepoRoot   string   `json:"repo_root"`
	D2APath    string   `json:"d2a_path"`
	Inputs     []string `json:"inputs"`
	Outputs    []string `json:"outputs"`
}

func PrepareTests(repoRoot string) (string, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", fmt.Errorf("resolve repo root: %w", err)
	}

	meta, err := loadMetadata(repoRoot)
	if err != nil {
		return "", err
	}

	if err := ensureDerivation(repoRoot); err != nil {
		return "", err
	}

	for rel, content := range testTasks(meta.TargetRepo) {
		path := filepath.Join(repoRoot, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return "", fmt.Errorf("create tests dir %s: %w", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return "", fmt.Errorf("write test file %s: %w", path, err)
		}
	}

	m := manifest{
		TargetRepo: meta.TargetRepo,
		RepoRoot:   repoRoot,
		D2APath:    filepath.Join(repoRoot, ".d2a"),
		Inputs: []string{
			"docs/implementation/00_mini_scope.md",
			"docs/implementation/01_mini_design.md",
			"docs/implementation/02_build_plan.md",
			"docs/implementation/03_test_plan.md",
			"src/ARCHITECTURE.md",
		},
		Outputs: []string{
			"tests/README.md",
			"tests/01_integration_tasks.md",
		},
	}
	if err := writeJSON(filepath.Join(repoRoot, ".d2a", "test-mini.json"), m); err != nil {
		return "", err
	}

	return meta.TargetRepo, nil
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

func ensureDerivation(repoRoot string) error {
	path := filepath.Join(repoRoot, "src", "ARCHITECTURE.md")
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("mini derivation not prepared: %s", repoRoot)
		}
		return fmt.Errorf("check src architecture file %s: %w", path, err)
	}
	return nil
}

func writeJSON(path string, value any) error {
	content, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	content = append(content, '\n')

	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write json file %s: %w", path, err)
	}

	return nil
}

func testTasks(target string) map[string]string {
	return map[string]string{
		"tests/README.md":               testsReadme(target),
		"tests/01_integration_tasks.md": integrationTasks(target),
	}
}

func testsReadme(target string) string {
	return fmt.Sprintf(`# tests

## Target

- Target repo: %s

## Purpose

This directory holds the incremental integration testing surface for the mini implementation under src/.

## Inputs

- docs/implementation/00_mini_scope.md
- docs/implementation/01_mini_design.md
- docs/implementation/02_build_plan.md
- docs/implementation/03_test_plan.md
- src/ARCHITECTURE.md

## Output Rules

- Start with one end-to-end test for the smallest runnable slice.
- Prefer observable system behavior over internal unit detail.
- Add tests only when they clarify the preserved architecture idea.
`, target)
}

func integrationTasks(target string) string {
	return fmt.Sprintf(`# 01. Integration Tasks

## Target

- Target repo: %s

## Atomic Questions

1. What is the first end-to-end test for the smallest runnable slice?
2. What setup or fixture is required to run it?
3. What input should drive the mini system?
4. What observable output or state change proves the architecture idea is working?
5. What single failure case should be tested next?

## Deliverables

- One first integration scenario
- One next scenario to add after the first pass
- The observable success signal
- The first failure signal worth keeping
`, target)
}
