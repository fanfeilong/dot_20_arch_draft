package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

type reportPrefs struct {
	Lang  string `json:"lang"`
	Theme string `json:"theme"`
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
	reportDir := filepath.Join(repoRoot, "report")
	dataDir := filepath.Join(reportDir, "data")
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
		AvailableArtifacts: []string{"LAB.md", ".d2a/target.json", ".d2a/test-mini.json", ".d2a/challenge.json", ".d2a/challenge_log.jsonl"},
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

	indexPath := filepath.Join(reportDir, "index.md")
	if err := os.WriteFile(indexPath, []byte(reportIndex(s, challenge)), 0o644); err != nil {
		return "", fmt.Errorf("write report index %s: %w", indexPath, err)
	}
	htmlPath := filepath.Join(reportDir, "index.html")
	if err := os.WriteFile(htmlPath, []byte(reportHTML(s, challenge)), 0o644); err != nil {
		return "", fmt.Errorf("write report html %s: %w", htmlPath, err)
	}
	if err := writeBriefArtifacts(repoRoot, s, manifest, challenge); err != nil {
		return "", err
	}
	if err := ensureVueSkeleton(reportDir); err != nil {
		return "", err
	}

	return meta.TargetRepo, nil
}

func loadTargetMetadata(repoRoot string) (targetMetadata, error) {
	d2aFile := filepath.Join(repoRoot, "LAB.md")
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
- report/brief.md
- report/brief.html

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

func writeBriefArtifacts(repoRoot string, s summary, manifest testManifest, challenge state.ChallengeSnapshot) error {
	reportDir := filepath.Join(repoRoot, "report")
	prefs := loadReportPrefs(repoRoot)
	page1 := buildArchitecturePage(repoRoot)
	page2 := buildMiniPage(repoRoot, manifest, challenge)

	mdPath := filepath.Join(reportDir, "brief.md")
	mdContent := buildBriefMarkdown(s, page1, page2, prefs)
	if err := os.WriteFile(mdPath, []byte(mdContent), 0o644); err != nil {
		return fmt.Errorf("write brief markdown %s: %w", mdPath, err)
	}

	htmlPath := filepath.Join(reportDir, "brief.html")
	htmlContent := buildBriefHTML(s, page1, page2, prefs)
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0o644); err != nil {
		return fmt.Errorf("write brief html %s: %w", htmlPath, err)
	}
	return nil
}

func loadReportPrefs(repoRoot string) reportPrefs {
	p := reportPrefs{
		Lang:  "zh",
		Theme: "ocean",
	}
	path := filepath.Join(repoRoot, ".d2a", "report_prefs.json")
	content, err := os.ReadFile(path)
	if err != nil {
		return p
	}
	var raw reportPrefs
	if err := json.Unmarshal(content, &raw); err != nil {
		return p
	}
	lang := strings.ToLower(strings.TrimSpace(raw.Lang))
	theme := strings.ToLower(strings.TrimSpace(raw.Theme))
	switch lang {
	case "zh", "en":
		p.Lang = lang
	}
	switch theme {
	case "ocean", "sunset", "slate":
		p.Theme = theme
	}
	return p
}

type briefPageOne struct {
	StateDiagram string
	SixElements  [][2]string
}

type briefPageTwo struct {
	TargetStack             string
	ProviderSummary         string
	TimeboxSummary          string
	IntentSummary           string
	RunnableSliceSummary    string
	BuildSummary            string
	TestEvidence            string
	IntentionalOmissions    string
	ChallengeRecommendation string
}

func buildArchitecturePage(repoRoot string) briefPageOne {
	p := briefPageOne{
		StateDiagram: extractMermaidStateDiagram(filepath.Join(repoRoot, "docs", "architecture", "04_state_evolution.md")),
	}
	if p.StateDiagram == "" {
		p.StateDiagram = "stateDiagram-v2\n  [*] --> Unknown\n  Unknown --> Unknown: fill docs/architecture/04_state_evolution.md"
	}
	p.SixElements = [][2]string{
		{"Boundary", summarizeDoc(filepath.Join(repoRoot, "docs", "architecture", "01_boundary.md"), 140)},
		{"Driver", summarizeDoc(filepath.Join(repoRoot, "docs", "architecture", "02_driver.md"), 140)},
		{"Core Objects", summarizeDoc(filepath.Join(repoRoot, "docs", "architecture", "03_core_objects.md"), 140)},
		{"State Evolution", summarizeDoc(filepath.Join(repoRoot, "docs", "architecture", "04_state_evolution.md"), 140)},
		{"Cooperation", summarizeDoc(filepath.Join(repoRoot, "docs", "architecture", "05_cooperation.md"), 140)},
		{"Constraints", summarizeDoc(filepath.Join(repoRoot, "docs", "architecture", "06_constraints.md"), 140)},
	}
	return p
}

func buildMiniPage(repoRoot string, manifest testManifest, challenge state.ChallengeSnapshot) briefPageTwo {
	stack, provider, timebox, intent := loadMiniGateSummary(filepath.Join(repoRoot, ".d2a", "mini_gate", "d2a-mini-1-scope.json"))
	return briefPageTwo{
		TargetStack:             defaultUnknown(stack),
		ProviderSummary:         defaultUnknown(provider),
		TimeboxSummary:          defaultUnknown(timebox),
		IntentSummary:           defaultUnknown(intent),
		RunnableSliceSummary:    summarizeDoc(filepath.Join(repoRoot, "docs", "implementation", "00_mini_scope.md"), 180),
		BuildSummary:            summarizeDoc(filepath.Join(repoRoot, "docs", "implementation", "02_build_plan.md"), 180),
		TestEvidence:            summarizeList(manifest.Outputs, 4),
		IntentionalOmissions:    summarizeOmissions(filepath.Join(repoRoot, "docs", "implementation", "02_build_plan.md"), filepath.Join(repoRoot, "src", "ARCHITECTURE.md")),
		ChallengeRecommendation: valueOrUnknown(challenge.Recommendation),
	}
}

func buildBriefMarkdown(s summary, page1 briefPageOne, page2 briefPageTwo, prefs reportPrefs) string {
	labels := briefLabels(prefs.Lang)
	rows := make([]string, 0, len(page1.SixElements))
	for _, pair := range page1.SixElements {
		rows = append(rows, fmt.Sprintf("| %s | %s |", pair[0], sanitizeInline(pair[1])))
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", labels["brief_title"])
	fmt.Fprintf(&b, "%s: %s\n\n", labels["target_repo"], s.TargetRepo)
	fmt.Fprintf(&b, "## %s\n\n", labels["page1_title"])
	fmt.Fprintf(&b, "### %s\n\n", labels["diagram_title"])
	fmt.Fprintf(&b, "```mermaid\n%s\n```\n\n", page1.StateDiagram)
	fmt.Fprintf(&b, "### %s\n\n", labels["six_elements"])
	fmt.Fprintf(&b, "| %s | %s |\n| --- | --- |\n%s\n\n", labels["element"], labels["key_point"], strings.Join(rows, "\n"))
	fmt.Fprintf(&b, "## %s\n\n", labels["page2_title"])
	fmt.Fprintf(&b, "- %s: %s\n", labels["target_stack"], limitText(page2.TargetStack, 120))
	fmt.Fprintf(&b, "- %s: %s\n", labels["provider_gate"], limitText(page2.ProviderSummary, 180))
	fmt.Fprintf(&b, "- %s: %s\n", labels["timebox_gate"], limitText(page2.TimeboxSummary, 180))
	fmt.Fprintf(&b, "- %s: %s\n", labels["intent_gate"], limitText(page2.IntentSummary, 180))
	fmt.Fprintf(&b, "- %s: %s\n", labels["slice"], limitText(page2.RunnableSliceSummary, 220))
	fmt.Fprintf(&b, "- %s: %s\n", labels["build_summary"], limitText(page2.BuildSummary, 220))
	fmt.Fprintf(&b, "- %s: %s\n", labels["test_evidence"], limitText(page2.TestEvidence, 220))
	fmt.Fprintf(&b, "- %s: %s\n", labels["omissions"], limitText(page2.IntentionalOmissions, 220))
	fmt.Fprintf(&b, "- %s: %s\n", labels["challenge_rec"], limitText(page2.ChallengeRecommendation, 80))
	return b.String()
}

func buildBriefHTML(s summary, page1 briefPageOne, page2 briefPageTwo, prefs reportPrefs) string {
	labels := briefLabels(prefs.Lang)
	theme := briefTheme(prefs.Theme)
	rows := make([]string, 0, len(page1.SixElements))
	for _, pair := range page1.SixElements {
		rows = append(rows, fmt.Sprintf("<tr><th scope=\"row\">%s</th><td>%s</td></tr>", escapeHTML(pair[0]), escapeHTML(limitText(pair[1], 170))))
	}
	return fmt.Sprintf(`<!doctype html>
<html lang="%s">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>d2a Brief</title>
  <style>
    :root {
      --ink: #1f2a30;
      --subtle: #5f6d73;
      --line: #d6dde2;
      --panel: #ffffff;
      --accent: %s;
      --accent-soft: %s;
      --paper: #f3f7f8;
      --bg1: %s;
      --bg2: %s;
      --bar1: %s;
      --bar2: %s;
      --bar3: %s;
    }
    @page { size: A4; margin: 12mm; }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: "Source Serif 4", "Iowan Old Style", "Palatino Linotype", Georgia, serif;
      color: var(--ink);
      background:
        radial-gradient(1200px 500px at 10%% -10%%, var(--bg1), transparent 60%%),
        radial-gradient(1200px 500px at 95%% 0%%, var(--bg2), transparent 62%%),
        var(--paper);
      padding: 20px 0 40px;
    }
    .page {
      width: 186mm;
      min-height: 272mm;
      margin: 0 auto 14px;
      padding: 8mm;
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 12px;
      box-shadow: 0 18px 40px rgba(22, 35, 42, 0.08);
      page-break-after: always;
      position: relative;
      overflow: hidden;
    }
    .page:last-child { page-break-after: auto; }
    .page::before {
      content: "";
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 6mm;
      background: linear-gradient(90deg, #136f63, #2a8f7e, #5da399);
      background: linear-gradient(90deg, var(--bar1), var(--bar2), var(--bar3));
      opacity: 0.9;
    }
    .top {
      display: flex;
      justify-content: space-between;
      align-items: flex-end;
      margin-top: 7mm;
      margin-bottom: 5mm;
      gap: 8mm;
    }
    .kicker {
      margin: 0;
      font-size: 11px;
      text-transform: uppercase;
      letter-spacing: 0.08em;
      color: var(--accent);
      font-weight: 700;
    }
    h1 {
      margin: 2px 0 0;
      font-size: 24px;
      line-height: 1.1;
      letter-spacing: 0.01em;
    }
    .repo {
      margin: 0;
      color: var(--subtle);
      font-size: 11px;
      max-width: 60%%;
      text-align: right;
      word-break: break-word;
    }
    h2 {
      margin: 0 0 6px;
      font-size: 14px;
      letter-spacing: 0.01em;
    }
    .card {
      border: 1px solid var(--line);
      border-radius: 10px;
      background: #fff;
      padding: 9px;
      margin-bottom: 8px;
    }
    .diagram pre {
      margin: 0;
      background: #f7fafb;
      border: 1px solid #e2e9ed;
      border-radius: 8px;
      padding: 8px;
      white-space: pre-wrap;
      font-size: 10.5px;
      line-height: 1.32;
      max-height: 95mm;
      overflow: hidden;
      color: #20353b;
    }
    table {
      width: 100%%;
      border-collapse: collapse;
      table-layout: fixed;
      font-size: 11px;
    }
    th, td {
      border: 1px solid var(--line);
      padding: 6px 7px;
      vertical-align: top;
      line-height: 1.3;
      text-align: left;
    }
    th[scope="row"] {
      width: 28%%;
      font-weight: 700;
      background: var(--accent-soft);
    }
    .mini-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 8px;
    }
    .mini-item {
      border: 1px solid var(--line);
      border-radius: 8px;
      padding: 7px;
      background: #fff;
    }
    .mini-item strong {
      display: block;
      margin-bottom: 3px;
      font-size: 11px;
      color: #174d45;
    }
    .mini-item p {
      margin: 0;
      font-size: 10.5px;
      line-height: 1.3;
      color: #32424a;
    }
    .footer-note {
      margin-top: 8px;
      font-size: 10px;
      color: var(--subtle);
      text-align: right;
    }
    @media (max-width: 900px) {
      body { padding: 0; background: var(--paper); }
      .page {
        width: 100%%;
        min-height: auto;
        margin: 0;
        border-radius: 0;
        border-left: none;
        border-right: none;
        padding: 16px 14px 20px;
      }
      .page::before { height: 5px; }
      .top {
        margin-top: 12px;
        flex-direction: column;
        align-items: flex-start;
      }
      .repo {
        text-align: left;
        max-width: 100%%;
      }
      .mini-grid { grid-template-columns: 1fr; }
    }
    @media print {
      body { background: #fff; padding: 0; }
      .page {
        margin: 0;
        border: none;
        box-shadow: none;
        border-radius: 0;
      }
    }
  </style>
</head>
<body>
  <section class="page">
    <header class="top">
      <div>
        <p class="kicker">%s</p>
        <h1>%s</h1>
      </div>
      <p class="repo">%s: %s</p>
    </header>
    <section class="card diagram">
      <h2>%s</h2>
      <pre>%s</pre>
    </section>
    <section class="card">
      <h2>%s</h2>
      <table>
        <tbody>%s</tbody>
      </table>
    </section>
    <div class="footer-note">%s</div>
  </section>
  <section class="page">
    <header class="top">
      <div>
        <p class="kicker">%s</p>
        <h1>%s</h1>
      </div>
      <p class="repo">%s: %s</p>
    </header>
    <section class="mini-grid">
      <article class="mini-item"><strong>%s</strong><p>%s</p></article>
      <article class="mini-item"><strong>%s</strong><p>%s</p></article>
      <article class="mini-item"><strong>%s</strong><p>%s</p></article>
      <article class="mini-item"><strong>%s</strong><p>%s</p></article>
      <article class="mini-item"><strong>%s</strong><p>%s</p></article>
      <article class="mini-item"><strong>%s</strong><p>%s</p></article>
      <article class="mini-item"><strong>%s</strong><p>%s</p></article>
      <article class="mini-item"><strong>%s</strong><p>%s</p></article>
    </section>
    <section class="card">
      <h2>%s</h2>
      <p>%s</p>
    </section>
    <div class="footer-note">%s</div>
  </section>
</body>
</html>
`,
		escapeHTML(prefs.Lang),
		escapeHTML(theme["accent"]),
		escapeHTML(theme["accent_soft"]),
		escapeHTML(theme["bg1"]),
		escapeHTML(theme["bg2"]),
		escapeHTML(theme["bar1"]),
		escapeHTML(theme["bar2"]),
		escapeHTML(theme["bar3"]),
		escapeHTML(labels["kicker_arch"]),
		escapeHTML(labels["page1_h1"]),
		escapeHTML(labels["target_repo"]),
		escapeHTML(s.TargetRepo),
		escapeHTML(labels["diagram_title"]),
		escapeHTML(limitText(page1.StateDiagram, 1400)),
		escapeHTML(labels["six_elements"]),
		strings.Join(rows, ""),
		escapeHTML(labels["footer_arch"]),
		escapeHTML(labels["kicker_mini"]),
		escapeHTML(labels["page2_h1"]),
		escapeHTML(labels["target_repo"]),
		escapeHTML(s.TargetRepo),
		escapeHTML(labels["target_stack"]),
		escapeHTML(limitText(page2.TargetStack, 120)),
		escapeHTML(labels["provider_gate"]),
		escapeHTML(limitText(page2.ProviderSummary, 180)),
		escapeHTML(labels["timebox_gate"]),
		escapeHTML(limitText(page2.TimeboxSummary, 180)),
		escapeHTML(labels["intent_gate"]),
		escapeHTML(limitText(page2.IntentSummary, 180)),
		escapeHTML(labels["slice"]),
		escapeHTML(limitText(page2.RunnableSliceSummary, 220)),
		escapeHTML(labels["build_summary"]),
		escapeHTML(limitText(page2.BuildSummary, 220)),
		escapeHTML(labels["test_evidence"]),
		escapeHTML(limitText(page2.TestEvidence, 220)),
		escapeHTML(labels["omissions"]),
		escapeHTML(limitText(page2.IntentionalOmissions, 220)),
		escapeHTML(labels["challenge_rec"]),
		escapeHTML(limitText(page2.ChallengeRecommendation, 80)),
		escapeHTML(labels["footer_mini"]),
	)
}

func briefLabels(lang string) map[string]string {
	if lang == "en" {
		return map[string]string{
			"brief_title":   "d2a Brief (2-Page A4)",
			"target_repo":   "Target repo",
			"page1_title":   "Page 1: Architecture Core",
			"diagram_title": "State Machine Diagram",
			"six_elements":  "Six Elements (Compact)",
			"element":       "Element",
			"key_point":     "Key Point",
			"page2_title":   "Page 2: Mini Implementation Brief",
			"target_stack":  "Target stack",
			"provider_gate": "Provider gate",
			"timebox_gate":  "Timebox gate",
			"intent_gate":   "Intent gate",
			"slice":         "Runnable 20% slice",
			"build_summary": "Build summary",
			"test_evidence": "Test evidence",
			"omissions":     "Intentional omissions",
			"challenge_rec": "Challenge recommendation",
			"kicker_arch":   "d2a architecture brief · page 1",
			"page1_h1":      "Architecture Core",
			"footer_arch":   "d2a · architecture skeleton snapshot",
			"kicker_mini":   "d2a mini brief · page 2",
			"page2_h1":      "Mini Implementation Brief",
			"footer_mini":   "d2a · mini delivery snapshot",
		}
	}
	return map[string]string{
		"brief_title":   "d2a 简报（2 页 A4）",
		"target_repo":   "目标仓库",
		"page1_title":   "第 1 页：架构核心",
		"diagram_title": "状态机图",
		"six_elements":  "六要素（极简）",
		"element":       "要素",
		"key_point":     "关键点",
		"page2_title":   "第 2 页：Mini 实现简报",
		"target_stack":  "目标技术栈",
		"provider_gate": "Provider Gate",
		"timebox_gate":  "Timebox Gate",
		"intent_gate":   "Intent Gate",
		"slice":         "可运行 20%% 切片",
		"build_summary": "构建摘要",
		"test_evidence": "测试证据",
		"omissions":     "刻意未实现项",
		"challenge_rec": "挑战建议",
		"kicker_arch":   "d2a 架构简报 · 第 1 页",
		"page1_h1":      "架构核心",
		"footer_arch":   "d2a · 架构骨架快照",
		"kicker_mini":   "d2a mini 简报 · 第 2 页",
		"page2_h1":      "Mini 实现简报",
		"footer_mini":   "d2a · mini 交付快照",
	}
}

func briefTheme(name string) map[string]string {
	switch name {
	case "sunset":
		return map[string]string{
			"accent":      "#a64b2a",
			"accent_soft": "#fbe9df",
			"bg1":         "#f4d9c8",
			"bg2":         "#efe1b9",
			"bar1":        "#a64b2a",
			"bar2":        "#c06a3f",
			"bar3":        "#d28c54",
		}
	case "slate":
		return map[string]string{
			"accent":      "#2f4a5f",
			"accent_soft": "#e8eef3",
			"bg1":         "#dbe3ea",
			"bg2":         "#d8d6df",
			"bar1":        "#2f4a5f",
			"bar2":        "#44657f",
			"bar3":        "#5d829f",
		}
	default:
		return map[string]string{
			"accent":      "#136f63",
			"accent_soft": "#e7f4f1",
			"bg1":         "#dbeef0",
			"bg2":         "#f2e9dc",
			"bar1":        "#136f63",
			"bar2":        "#2a8f7e",
			"bar3":        "#5da399",
		}
	}
}

func extractMermaidStateDiagram(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	text := string(content)
	const open = "```mermaid"
	start := strings.Index(text, open)
	if start < 0 {
		return ""
	}
	rest := text[start+len(open):]
	end := strings.Index(rest, "```")
	if end < 0 {
		return ""
	}
	return strings.TrimSpace(rest[:end])
}

func summarizeDoc(path string, maxChars int) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return "unknown"
	}
	lines := strings.Split(string(content), "\n")
	parts := make([]string, 0, 3)
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "```") {
			continue
		}
		line = strings.TrimLeft(line, "-*0123456789. ")
		if line == "" {
			continue
		}
		parts = append(parts, line)
		if len(parts) >= 2 {
			break
		}
	}
	if len(parts) == 0 {
		return "unknown"
	}
	return limitText(strings.Join(parts, " ; "), maxChars)
}

func loadMiniGateSummary(path string) (stack, provider, timebox, intent string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", "", "", ""
	}
	var raw map[string]any
	if err := json.Unmarshal(content, &raw); err != nil {
		return "", "", "", ""
	}

	stack = firstString(raw, "stack", "final_stack", "selected_stack", "target_stack", "tech_stack")
	provider = firstString(raw, "provider", "provider_match", "provider_decision")
	timebox = firstString(raw, "timebox", "timebox_budget", "timebox_decision")
	intent = firstString(raw, "intent", "intent_anchor", "intent_decision")
	if provider == "" || timebox == "" || intent == "" {
		parts := collectStringPairs(raw, 8)
		if provider == "" {
			provider = findPair(parts, "provider")
		}
		if timebox == "" {
			timebox = findPair(parts, "timebox")
		}
		if intent == "" {
			intent = findPair(parts, "intent")
		}
		if stack == "" {
			stack = findPair(parts, "stack")
		}
	}
	return stack, provider, timebox, intent
}

func firstString(data map[string]any, keys ...string) string {
	for _, key := range keys {
		v, ok := data[key]
		if !ok {
			continue
		}
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func collectStringPairs(data map[string]any, limit int) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, limit)
	for _, key := range keys {
		if len(parts) >= limit {
			break
		}
		v := data[key]
		s, ok := v.(string)
		if !ok {
			continue
		}
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%s", key, s))
	}
	return parts
}

func findPair(parts []string, needle string) string {
	for _, p := range parts {
		if strings.Contains(strings.ToLower(p), needle) {
			return p
		}
	}
	return ""
}

func summarizeOmissions(primaryPath, fallbackPath string) string {
	primary := summarizeDoc(primaryPath, 220)
	if primary != "unknown" {
		return primary
	}
	return summarizeDoc(fallbackPath, 220)
}

func summarizeList(items []string, maxItems int) string {
	if len(items) == 0 {
		return "unknown"
	}
	if len(items) > maxItems {
		items = items[:maxItems]
	}
	return strings.Join(items, ", ")
}

func sanitizeInline(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "|", "/")
	return strings.TrimSpace(s)
}

func limitText(s string, maxChars int) string {
	runes := []rune(strings.TrimSpace(s))
	if len(runes) <= maxChars || maxChars <= 0 {
		if len(runes) == 0 {
			return "unknown"
		}
		return string(runes)
	}
	return string(runes[:maxChars-1]) + "…"
}

func defaultUnknown(s string) string {
	if strings.TrimSpace(s) == "" {
		return "unknown"
	}
	return strings.TrimSpace(s)
}

func escapeHTML(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&#39;",
	)
	return replacer.Replace(s)
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
