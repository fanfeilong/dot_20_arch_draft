package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type metadata struct {
	TargetRepo    string `json:"target_repo"`
	TargetRepoURL string `json:"target_repo_url,omitempty"`
	RepoRoot      string `json:"repo_root"`
	D2APath       string `json:"d2a_path"`
}

func Analyze(targetRepo, repoRoot string) (string, error) {
	if strings.TrimSpace(targetRepo) == "" {
		return "", fmt.Errorf("target repo must not be empty")
	}

	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", fmt.Errorf("resolve repo root: %w", err)
	}

	if err := ensureD2ARepo(repoRoot); err != nil {
		return "", err
	}

	resolvedTarget, err := resolveTarget(targetRepo)
	if err != nil {
		return "", err
	}

	d2aDir := filepath.Join(repoRoot, ".d2a")

	meta := metadata{
		TargetRepo: resolvedTarget,
		RepoRoot:   repoRoot,
		D2APath:    d2aDir,
	}
	if err := writeJSON(filepath.Join(d2aDir, "target.json"), meta); err != nil {
		return "", err
	}

	for rel, content := range architectureTasks(resolvedTarget) {
		path := filepath.Join(repoRoot, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return "", fmt.Errorf("create architecture dir %s: %w", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return "", fmt.Errorf("write analysis file %s: %w", path, err)
		}
	}

	return resolvedTarget, nil
}

func ensureD2ARepo(repoRoot string) error {
	d2aFile := filepath.Join(repoRoot, "LAB.md")
	if _, err := os.Stat(d2aFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("d2a not initialized in repository: %s", repoRoot)
		}
		return fmt.Errorf("check d2a file %s: %w", d2aFile, err)
	}
	return nil
}

func resolveTarget(targetRepo string) (string, error) {
	if filepath.IsAbs(targetRepo) {
		return targetRepo, nil
	}

	if stat, err := os.Stat(targetRepo); err == nil && stat.IsDir() {
		abs, err := filepath.Abs(targetRepo)
		if err != nil {
			return "", fmt.Errorf("resolve target repo: %w", err)
		}
		return abs, nil
	}

	return targetRepo, nil
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

func architectureTasks(target string) map[string]string {
	return map[string]string{
		"docs/architecture/00_overview.md":        overviewTask(target),
		"docs/architecture/01_boundary.md":        boundaryTask(target),
		"docs/architecture/02_driver.md":          driverTask(target),
		"docs/architecture/03_core_objects.md":    coreObjectsTask(target),
		"docs/architecture/04_state_evolution.md": stateEvolutionTask(target),
		"docs/architecture/05_cooperation.md":     cooperationTask(target),
		"docs/architecture/06_constraints.md":     constraintsTask(target),
		"docs/architecture/99_code_map.md":        codeMapTask(target),
	}
}

func sharedHeader(title, target, skill string) string {
	return fmt.Sprintf(`# %s

## Target

- Target repo: %s
- Primary skill: %s

## Output Rules

- Keep the result small and numbered.
- Prefer three to six bullets or one short flow.
- After drafting, run:
  - compression pass
  - de-jargon pass
  - code-evidence pass

`, title, target, skill)
}

func overviewTask(target string) string {
	return sharedHeader("00. Overview", target, "d2a-step") + `## To Produce

- One-sentence system definition
- The one capability that must remain if 80% of the code were deleted
- The core architecture idea worth preserving in a mini implementation
- The four things a reader should understand first

## Inputs

- docs/architecture/01_boundary.md
- docs/architecture/02_driver.md
- docs/architecture/03_core_objects.md
- docs/architecture/04_state_evolution.md
- docs/architecture/05_cooperation.md
- docs/architecture/06_constraints.md
`
}

func boundaryTask(target string) string {
	return sharedHeader("01. Boundary", target, "d2a-arch-1-project-scope") + `## Atomic Questions

1. What kind of system is this?
2. What is the one-sentence definition?
3. What capability is non-removable?
4. What are the one to three best entry points?
5. What is inside the system boundary?
6. What is outside the system boundary?
`
}

func driverTask(target string) string {
	return sharedHeader("02. Driver", target, "d2a-arch-2-runtime-view") + `## Atomic Questions

1. What is the dominant runtime driver?
2. What is the core runtime loop?
3. Which module is the engine?
4. Which supporting modules are required to understand the engine?
`
}

func coreObjectsTask(target string) string {
	return sharedHeader("03. Core Objects", target, "d2a-arch-3-core-objects") + `## Atomic Questions

1. What are the at most three core objects?
2. Who creates, consumes, persists, or drives them?
3. Where is the state center?
`
}

func stateEvolutionTask(target string) string {
	return sharedHeader("04. State Evolution", target, "d2a-arch-4-state-evolution") + `## Atomic Questions

1. What single object or workflow should be tracked?
2. What are its three to six state stages?
3. What triggers the main state transitions?
4. Where is state stored, observed, or reconstructed?
`
}

func cooperationTask(target string) string {
	return sharedHeader("05. Cooperation", target, "d2a-arch-5-module-view") + `## Atomic Questions

1. What are the three to six top-level modules?
2. What is the responsibility of each module?
3. What is the minimum dependency structure?
4. What is the single most important cooperation chain?
5. Which module is most complex and why?
`
}

func constraintsTask(target string) string {
	return sharedHeader("06. Constraints", target, "d2a-arch-6-tradeoff-view") + `## Atomic Questions

1. What are the two to four hard constraints?
2. Which one is dominant?
3. What tradeoff does it force?
4. Which structures must remain if the system is rewritten?
5. Which large-looking parts are implementation detail rather than architecture core?
`
}

func codeMapTask(target string) string {
	return sharedHeader("99. Code Map", target, "d2a-step") + `## To Produce

- Claim-to-code mapping for the six architecture elements
- A suggested reading order
- A short list of parts that can be skipped at first

## Required Evidence Shape

- Prefer directories and files over broad prose.
- Map each major claim back to code locations.
`
}
