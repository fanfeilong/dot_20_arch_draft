package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fanfeilong/dot_20_arch_draft/internal/state"
)

type targetMetadata struct {
	TargetRepo string `json:"target_repo"`
	RepoRoot   string `json:"repo_root"`
	D2APath    string `json:"d2a_path"`
}

type testManifest struct {
	TargetRepo string   `json:"target_repo"`
	RepoRoot   string   `json:"repo_root"`
	D2APath    string   `json:"d2a_path"`
	Inputs     []string `json:"inputs"`
	Outputs    []string `json:"outputs"`
}

type summary struct {
	TargetRepo         string   `json:"target_repo"`
	RepoRoot           string   `json:"repo_root"`
	D2APath            string   `json:"d2a_path"`
	ArchitectureDocs   []string `json:"architecture_docs"`
	ChallengeArtifacts []string `json:"challenge_artifacts"`
	ImplementationDocs []string `json:"implementation_docs"`
	TestDocs           []string `json:"test_docs"`
	SourceDocs         []string `json:"source_docs"`
	AvailableArtifacts []string `json:"available_artifacts"`
}

func BuildReport(repoRoot string) (string, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", fmt.Errorf("resolve repo root: %w", err)
	}

	meta, err := loadTargetMetadata(repoRoot)
	if err != nil {
		return "", err
	}
	manifest, err := loadTestManifest(repoRoot)
	if err != nil {
		return "", err
	}
	challenge, err := loadChallengeState(repoRoot)
	if err != nil {
		return "", err
	}

	d2aDir := filepath.Join(repoRoot, ".d2a")
	dataDir := filepath.Join(d2aDir, "report", "data")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return "", fmt.Errorf("create report data dir: %w", err)
	}

	s := summary{
		TargetRepo:         meta.TargetRepo,
		RepoRoot:           repoRoot,
		D2APath:            d2aDir,
		ArchitectureDocs:   architectureDocs(),
		ChallengeArtifacts: []string{".d2a/challenge.json", ".d2a/challenge_log.jsonl"},
		ImplementationDocs: implementationDocs(),
		TestDocs:           manifest.Outputs,
		SourceDocs:         []string{"src/README.md", "src/ARCHITECTURE.md"},
		AvailableArtifacts: []string{".d2a/LAB.md", ".d2a/target.json", ".d2a/test-mini.json", ".d2a/challenge.json", ".d2a/challenge_log.jsonl"},
	}

	if err := writeJSON(filepath.Join(dataDir, "summary.json"), s); err != nil {
		return "", err
	}
	if err := writeJSON(filepath.Join(dataDir, "target.json"), meta); err != nil {
		return "", err
	}
	if err := writeJSON(filepath.Join(dataDir, "tests.json"), manifest); err != nil {
		return "", err
	}
	if err := writeJSON(filepath.Join(dataDir, "challenge.json"), challenge); err != nil {
		return "", err
	}

	indexPath := filepath.Join(d2aDir, "report", "index.md")
	if err := os.WriteFile(indexPath, []byte(reportIndex(s, challenge)), 0o644); err != nil {
		return "", fmt.Errorf("write report index %s: %w", indexPath, err)
	}
	htmlPath := filepath.Join(d2aDir, "report", "index.html")
	if err := os.WriteFile(htmlPath, []byte(reportHTML(s, challenge)), 0o644); err != nil {
		return "", fmt.Errorf("write report html %s: %w", htmlPath, err)
	}
	if err := ensureVueSkeleton(filepath.Join(d2aDir, "report")); err != nil {
		return "", err
	}

	return meta.TargetRepo, nil
}

func loadTargetMetadata(repoRoot string) (targetMetadata, error) {
	d2aFile := filepath.Join(repoRoot, ".d2a", "LAB.md")
	if _, err := os.Stat(d2aFile); err != nil {
		if os.IsNotExist(err) {
			return targetMetadata{}, fmt.Errorf("d2a not initialized in repository: %s", repoRoot)
		}
		return targetMetadata{}, fmt.Errorf("check d2a file %s: %w", d2aFile, err)
	}

	path := filepath.Join(repoRoot, ".d2a", "target.json")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return targetMetadata{}, fmt.Errorf("analysis not prepared: %s", repoRoot)
		}
		return targetMetadata{}, fmt.Errorf("read target metadata %s: %w", path, err)
	}

	var meta targetMetadata
	if err := json.Unmarshal(content, &meta); err != nil {
		return targetMetadata{}, fmt.Errorf("parse target metadata %s: %w", path, err)
	}
	if strings.TrimSpace(meta.TargetRepo) == "" {
		return targetMetadata{}, fmt.Errorf("target metadata missing target repo: %s", path)
	}
	return meta, nil
}

func loadTestManifest(repoRoot string) (testManifest, error) {
	path := filepath.Join(repoRoot, ".d2a", "test-mini.json")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return testManifest{}, fmt.Errorf("test planning not prepared: %s", repoRoot)
		}
		return testManifest{}, fmt.Errorf("read test manifest %s: %w", path, err)
	}

	var m testManifest
	if err := json.Unmarshal(content, &m); err != nil {
		return testManifest{}, fmt.Errorf("parse test manifest %s: %w", path, err)
	}
	if strings.TrimSpace(m.TargetRepo) == "" {
		return testManifest{}, fmt.Errorf("test manifest missing target repo: %s", path)
	}
	return m, nil
}

func loadChallengeState(repoRoot string) (state.ChallengeSnapshot, error) {
	challenge, err := state.ChallengeState(repoRoot)
	if err == nil {
		return challenge, nil
	}
	if strings.Contains(err.Error(), "challenge state file not found") {
		return state.ChallengeSnapshot{
			RepoPath:       repoRoot,
			D2APath:        filepath.Join(repoRoot, ".d2a"),
			Recommendation: "missing",
		}, nil
	}
	return state.ChallengeSnapshot{}, err
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

func architectureDocs() []string {
	return []string{
		"docs/architecture/00_overview.md",
		"docs/architecture/01_boundary.md",
		"docs/architecture/02_driver.md",
		"docs/architecture/03_core_objects.md",
		"docs/architecture/04_state_evolution.md",
		"docs/architecture/05_cooperation.md",
		"docs/architecture/06_constraints.md",
		"docs/architecture/99_code_map.md",
	}
}

func implementationDocs() []string {
	return []string{
		"docs/implementation/00_mini_scope.md",
		"docs/implementation/01_mini_design.md",
		"docs/implementation/02_build_plan.md",
		"docs/implementation/03_test_plan.md",
	}
}

func reportIndex(s summary, challenge state.ChallengeSnapshot) string {
	return fmt.Sprintf(`# d2a Report

## Target

- Target repo: %s
- Repo root: %s
- d2a path: %s

## Report Sections

### Architecture

%s

### Challenge

- Current stage: %s
- Challenge progress: %s
- Current decision: %s
- Recommendation: %s
- Skip reason: %s

### Implementation

%s

### Tests

%s

### Source

%s

## Data Files

- report/data/summary.json
- report/data/target.json
- report/data/tests.json
- report/data/challenge.json

## Next Step

Use these data files as the input contract for the future Vue report app and local serve command.
`,
		s.TargetRepo,
		s.RepoRoot,
		s.D2APath,
		bulletList(s.ArchitectureDocs),
		valueOrUnknown(challenge.CurrentStage),
		progressString(challenge.DecisionIndex, challenge.DecisionTotal),
		valueOrUnknown(challenge.CurrentDecision),
		valueOrUnknown(challenge.Recommendation),
		valueOrUnknown(challenge.SkipReason),
		bulletList(s.ImplementationDocs),
		bulletList(s.TestDocs),
		bulletList(s.SourceDocs),
	)
}

func bulletList(items []string) string {
	lines := make([]string, 0, len(items))
	for _, item := range items {
		lines = append(lines, "- "+item)
	}
	return strings.Join(lines, "\n")
}

func htmlList(items []string) string {
	lines := make([]string, 0, len(items))
	for _, item := range items {
		lines = append(lines, "<li><code>"+item+"</code></li>")
	}
	return strings.Join(lines, "\n")
}

func reportHTML(s summary, challenge state.ChallengeSnapshot) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>d2a Report</title>
  <style>
    :root {
      --bg: #f4f1e8;
      --panel: #fffdf8;
      --ink: #1f1a14;
      --muted: #5f574d;
      --line: #d7cdbd;
      --accent: #8a5a2b;
    }
    body {
      margin: 0;
      font-family: Georgia, "Iowan Old Style", "Palatino Linotype", serif;
      background: linear-gradient(180deg, #efe8d9 0%%, var(--bg) 100%%);
      color: var(--ink);
    }
    main {
      max-width: 960px;
      margin: 0 auto;
      padding: 40px 20px 64px;
    }
    .hero, section {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 16px;
      padding: 24px;
      box-shadow: 0 10px 30px rgba(60, 40, 20, 0.05);
      margin-bottom: 20px;
    }
    h1, h2 {
      margin-top: 0;
    }
    p, li {
      color: var(--muted);
      line-height: 1.6;
    }
    code {
      background: #f2eadc;
      padding: 2px 6px;
      border-radius: 6px;
      color: var(--accent);
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      gap: 16px;
    }
    a {
      color: var(--accent);
    }
    #data-status {
      font-size: 13px;
      color: var(--accent);
      margin-top: 8px;
    }
  </style>
</head>
<body>
  <main>
    <section class="hero">
      <h1>d2a Report</h1>
      <p>This is the current local report surface for the repository-root d2a workflow.</p>
      <p><strong>Target repo:</strong> <code id="target-repo">%s</code></p>
      <p><strong>Repo root:</strong> <code id="repo-root">%s</code></p>
      <p><strong>d2a path:</strong> <code id="d2a-path">%s</code></p>
      <p id="data-status">Loading report data from ./data/*.json ...</p>
    </section>

    <div class="grid">
      <section>
        <h2>Architecture</h2>
        <ul id="architecture-list">%s</ul>
      </section>
      <section>
        <h2>Implementation</h2>
        <ul id="implementation-list">%s</ul>
      </section>
      <section>
        <h2>Challenge</h2>
        <ul>
          <li><strong>Stage:</strong> <code id="challenge-stage">%s</code></li>
          <li><strong>Progress:</strong> <code id="challenge-progress">%s</code></li>
          <li><strong>Decision:</strong> <code id="challenge-decision">%s</code></li>
          <li><strong>Recommendation:</strong> <code id="challenge-recommendation">%s</code></li>
          <li><strong>Skip reason:</strong> <code id="challenge-skip">%s</code></li>
        </ul>
      </section>
      <section>
        <h2>Tests</h2>
        <ul id="tests-list">%s</ul>
      </section>
      <section>
        <h2>Source</h2>
        <ul id="source-list">%s</ul>
      </section>
    </div>

    <section>
      <h2>Data Files</h2>
      <ul>
        <li><a href="./data/summary.json">summary.json</a></li>
        <li><a href="./data/target.json">target.json</a></li>
        <li><a href="./data/tests.json">tests.json</a></li>
        <li><a href="./data/challenge.json">challenge.json</a></li>
        <li><a href="./index.md">index.md</a></li>
      </ul>
      <p>This page reads these files at runtime without requiring Node.js.</p>
    </section>
  </main>
  <script>
    function text(id, value) {
      var el = document.getElementById(id);
      if (!el) return;
      el.textContent = value || "unknown";
    }

    function list(id, items) {
      var el = document.getElementById(id);
      if (!el || !Array.isArray(items)) return;
      el.innerHTML = items.map(function(item) {
        return "<li><code>" + String(item) + "</code></li>";
      }).join("");
    }

    function progress(index, total) {
      if (!total || total <= 0) return "none";
      return String(index) + "/" + String(total);
    }

    Promise.all([
      fetch("./data/summary.json").then(function(r) { return r.json(); }),
      fetch("./data/target.json").then(function(r) { return r.json(); }),
      fetch("./data/tests.json").then(function(r) { return r.json(); }),
      fetch("./data/challenge.json").then(function(r) { return r.json(); })
    ]).then(function(parts) {
      var summary = parts[0] || {};
      var target = parts[1] || {};
      var tests = parts[2] || {};
      var challenge = parts[3] || {};

      text("target-repo", target.target_repo || summary.target_repo);
      text("repo-root", target.repo_root || summary.repo_root);
      text("d2a-path", target.d2a_path || summary.d2a_path);

      list("architecture-list", summary.architecture_docs || []);
      list("implementation-list", summary.implementation_docs || []);
      list("tests-list", tests.outputs || summary.test_docs || []);
      list("source-list", summary.source_docs || []);

      text("challenge-stage", challenge.current_stage);
      text("challenge-progress", progress(challenge.decision_index, challenge.decision_total));
      text("challenge-decision", challenge.current_decision);
      text("challenge-recommendation", challenge.recommendation);
      text("challenge-skip", challenge.skip_reason);

      text("data-status", "Report data loaded from local JSON files.");
    }).catch(function(err) {
      text("data-status", "Failed to load report data: " + String(err));
    });
  </script>
</body>
</html>
`,
		s.TargetRepo,
		s.RepoRoot,
		s.D2APath,
		htmlList(s.ArchitectureDocs),
		htmlList(s.ImplementationDocs),
		valueOrUnknown(challenge.CurrentStage),
		progressString(challenge.DecisionIndex, challenge.DecisionTotal),
		valueOrUnknown(challenge.CurrentDecision),
		valueOrUnknown(challenge.Recommendation),
		valueOrUnknown(challenge.SkipReason),
		htmlList(s.TestDocs),
		htmlList(s.SourceDocs),
	)
}

func progressString(index, total int) string {
	if total <= 0 {
		return "none"
	}
	return fmt.Sprintf("%d/%d", index, total)
}

func valueOrUnknown(value string) string {
	if value == "" {
		return "unknown"
	}
	return value
}

func ensureVueSkeleton(reportDir string) error {
	files := map[string]string{
		"vue-app/package.json": `{
  "name": "d2a-report-vue-app",
  "private": true,
  "version": "0.0.1",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "vue": "^3.5.13"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.2.1",
    "vite": "^5.4.11"
  }
}
`,
		"vue-app/index.html": `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>d2a Report App</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.js"></script>
  </body>
</html>
`,
		"vue-app/vite.config.js": `import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  plugins: [vue()],
});
`,
		"vue-app/src/main.js": `import { createApp } from "vue";
import App from "./App.vue";
import "./style.css";

createApp(App).mount("#app");
`,
		"vue-app/src/App.vue": `<script setup>
import { onMounted, ref } from "vue";

const status = ref("Loading report data ...");
const target = ref({});
const summary = ref({});
const tests = ref({});
const challenge = ref({});

onMounted(async () => {
  try {
    const [summaryRes, targetRes, testsRes, challengeRes] = await Promise.all([
      fetch("/data/summary.json"),
      fetch("/data/target.json"),
      fetch("/data/tests.json"),
      fetch("/data/challenge.json"),
    ]);

    summary.value = await summaryRes.json();
    target.value = await targetRes.json();
    tests.value = await testsRes.json();
    challenge.value = await challengeRes.json();
    status.value = "Loaded data from /data/*.json";
  } catch (err) {
    status.value = ` + "`" + `Failed to load data: ${String(err)}` + "`" + `;
  }
});
</script>

<template>
  <main class="wrap">
    <section class="card">
      <p class="eyebrow">d2a report app</p>
      <h1>Vue Report App</h1>
      <p>{{ status }}</p>

      <h2>Target</h2>
      <ul>
        <li><code>{{ target.target_repo || "unknown" }}</code></li>
        <li><code>{{ target.repo_root || "unknown" }}</code></li>
        <li><code>{{ target.d2a_path || "unknown" }}</code></li>
      </ul>

      <h2>Architecture Docs</h2>
      <ul>
        <li v-for="item in summary.architecture_docs || []" :key="item">
          <code>{{ item }}</code>
        </li>
      </ul>

      <h2>Test Outputs</h2>
      <ul>
        <li v-for="item in tests.outputs || []" :key="item">
          <code>{{ item }}</code>
        </li>
      </ul>

      <h2>Challenge</h2>
      <ul>
        <li><code>{{ challenge.current_stage || "unknown" }}</code></li>
        <li><code>{{ challenge.recommendation || "unknown" }}</code></li>
      </ul>
    </section>
  </main>
</template>
`,
		"vue-app/src/style.css": `:root {
  --bg0: #fff8ec;
  --bg1: #f5e8cf;
  --ink: #1f1a14;
  --muted: #5f574d;
  --line: #d7cdbd;
  --accent: #8a5a2b;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: "Iowan Old Style", "Palatino Linotype", Georgia, serif;
  color: var(--ink);
  background: radial-gradient(circle at 10% 10%, var(--bg0), var(--bg1));
}

.wrap {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
}

.card {
  width: min(760px, 100%);
  border: 1px solid var(--line);
  border-radius: 20px;
  background: #fffdf8;
  padding: 28px;
  box-shadow: 0 10px 30px rgba(60, 40, 20, 0.08);
}

.eyebrow {
  margin: 0 0 8px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-size: 12px;
  color: var(--accent);
}

h1 {
  margin: 0 0 10px;
}

p,
li {
  color: var(--muted);
  line-height: 1.6;
}
`,
	}

	for rel, content := range files {
		path := filepath.Join(reportDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("create vue skeleton dir %s: %w", filepath.Dir(path), err)
		}
		if _, err := os.Stat(path); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("check vue skeleton file %s: %w", path, err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return fmt.Errorf("write vue skeleton file %s: %w", path, err)
		}
	}

	return nil
}
