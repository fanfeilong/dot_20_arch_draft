package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fanfeilong/dot_20_arch_draft/internal/analyzer"
	"github.com/fanfeilong/dot_20_arch_draft/internal/deriver"
	"github.com/fanfeilong/dot_20_arch_draft/internal/installer"
	"github.com/fanfeilong/dot_20_arch_draft/internal/reporter"
	"github.com/fanfeilong/dot_20_arch_draft/internal/server"
	"github.com/fanfeilong/dot_20_arch_draft/internal/state"
	"github.com/fanfeilong/dot_20_arch_draft/internal/tester"
)

const usage = `d2a initializes a repository-root workflow and installs built-in skills.

Usage:
  d2a help
  d2a init <target-repo-git-url> [--lang <zh|en>]
  d2a analyze [<target-repo>] [--repo <repo-dir>]
  d2a derive-mini [--repo <repo-dir>] [--skip-challenge-reason <text>]
  d2a test-mini [--repo <repo-dir>]
  d2a report [--repo <repo-dir>]
  d2a serve [--repo <repo-dir>]
  d2a status [--repo <repo-dir>]
  d2a skill-state <skill-name> [--repo <repo-dir>] [--status <started|progress|completed>] [--stage <stage>] [--phase <phase>] [--question-index <n>] [--question-total <n>] [--question <text>] [--answer <text>] [--evaluation <correct|partial|incorrect>] [--explanation <text>] [--next-step <text>] [--next-skill <name>] [--next-file <path>] [--decision <label>] [--strength <strong|partial|weak>] [--recommendation <proceed|review|revisit architecture>] [--objection <text>] [--summary <text>]
  d2a version
`

type repoContext struct {
	Name    string
	RepoDir string
	D2APath string
}

type initTargetMetadata struct {
	TargetRepo    string `json:"target_repo"`
	TargetRepoURL string `json:"target_repo_url,omitempty"`
	RepoRoot      string `json:"repo_root"`
	D2APath       string `json:"d2a_path"`
}

func Run(args []string, version string) error {
	return runWithIO(args, os.Stdout, version)
}

func runWithIO(args []string, stdout io.Writer, version string) error {
	if len(args) == 0 {
		printUsage(stdout)
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		printUsage(stdout)
		return nil
	case "version":
		_, err := fmt.Fprintf(stdout, "%s\n", version)
		return err
	case "init":
		targetRepoURL, language, err := parseInitArgs(args)
		if err != nil {
			return err
		}
		repoName, err := parseRepoNameFromGitURL(targetRepoURL)
		if err != nil {
			return err
		}

		workspaceDir, err := filepath.Abs(repoName + "_d2a")
		if err != nil {
			return fmt.Errorf("resolve workspace path: %w", err)
		}
		ctx, err := initRepoContext(workspaceDir)
		if err != nil {
			return err
		}
		if err := printRepoHeader(stdout, ctx); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(stdout, "target repo url: %s\n", targetRepoURL); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(stdout, "d2a language: %s\n", language); err != nil {
			return err
		}

		target, err := installer.InstallWithLanguage(workspaceDir, language)
		if err != nil {
			return err
		}
		if err := ensureWorkspaceGitRepo(target); err != nil {
			return err
		}
		clonedRepoPath, err := cloneTargetRepo(targetRepoURL, target, repoName)
		if err != nil {
			return err
		}
		if err := writeInitTargetMetadata(target, targetRepoURL, clonedRepoPath); err != nil {
			return err
		}
		if err := writeAgentsFile(target); err != nil {
			return err
		}
		if _, err := state.Initialize(target); err != nil {
			return err
		}

		_, err = fmt.Fprintf(stdout, "initialized d2a workspace %s\ncloned target repo to %s\n", target, clonedRepoPath)
		return err
	case "analyze":
		targetRepo, explicitRepo, err := parseAnalyzeArgs(args)
		if err != nil {
			return err
		}

		ctx, err := resolveRepoContext(explicitRepo)
		if err != nil {
			return err
		}
		if err := printRepoHeader(stdout, ctx); err != nil {
			return err
		}
		if targetRepo == "" {
			targetRepo, err = resolveDefaultTargetRepo(ctx.RepoDir)
			if err != nil {
				return err
			}
		}

		target, err := analyzer.Analyze(targetRepo, ctx.RepoDir)
		if err != nil {
			return err
		}
		if _, err := state.RecordCommand(
			ctx.RepoDir,
			"d2a analyze",
			state.StageAnalysisPrepared,
			"Use d2a architecture skills to fill docs/architecture/.",
			"Prepared architecture analysis task files under docs/architecture/.",
		); err != nil {
			return err
		}

		_, err = fmt.Fprintf(stdout, "prepared d2a analysis for %s in %s\n", target, ctx.D2APath)
		return err
	case "derive-mini":
		explicitRepo, skipReason, err := parseDeriveMiniArgs(args)
		if err != nil {
			return err
		}

		ctx, err := resolveRepoContext(explicitRepo)
		if err != nil {
			return err
		}
		if err := printRepoHeader(stdout, ctx); err != nil {
			return err
		}
		if skipReason != "" {
			if _, err := state.SkipChallenge(ctx.RepoDir, skipReason); err != nil {
				return err
			}
		}

		target, err := deriver.DeriveMini(ctx.RepoDir)
		if err != nil {
			return err
		}
		if _, err := state.RecordCommand(
			ctx.RepoDir,
			"d2a derive-mini",
			state.StageMiniDerivationPrepared,
			"Use d2a mini skills to fill docs/implementation/ and src/ARCHITECTURE.md.",
			"Prepared mini-implementation planning files under docs/implementation/.",
		); err != nil {
			return err
		}

		_, err = fmt.Fprintf(stdout, "prepared d2a mini derivation for %s in %s\n", target, ctx.D2APath)
		return err
	case "test-mini":
		explicitRepo, err := parseOptionalRepoArgs("test-mini", args)
		if err != nil {
			return err
		}

		ctx, err := resolveRepoContext(explicitRepo)
		if err != nil {
			return err
		}
		if err := printRepoHeader(stdout, ctx); err != nil {
			return err
		}

		target, err := tester.PrepareTests(ctx.RepoDir)
		if err != nil {
			return err
		}
		if _, err := state.RecordCommand(
			ctx.RepoDir,
			"d2a test-mini",
			state.StageTestPlanPrepared,
			"Use d2a testing skills to fill tests/ and create the first integration checks.",
			"Prepared test-planning files under tests/.",
		); err != nil {
			return err
		}

		_, err = fmt.Fprintf(stdout, "prepared d2a test plan for %s in %s\n", target, ctx.D2APath)
		return err
	case "report":
		explicitRepo, err := parseOptionalRepoArgs("report", args)
		if err != nil {
			return err
		}

		ctx, err := resolveRepoContext(explicitRepo)
		if err != nil {
			return err
		}
		if err := printRepoHeader(stdout, ctx); err != nil {
			return err
		}

		target, err := reporter.BuildReport(ctx.RepoDir)
		if err != nil {
			return err
		}
		if _, err := state.RecordCommand(
			ctx.RepoDir,
			"d2a report",
			state.StageReportPrepared,
			"Use d2a-report-build to refine report content, then run d2a serve.",
			"Prepared report artifacts under report/.",
		); err != nil {
			return err
		}

		_, err = fmt.Fprintf(stdout, "prepared d2a report for %s in %s\n", target, ctx.D2APath)
		return err
	case "serve":
		explicitRepo, err := parseOptionalRepoArgs("serve", args)
		if err != nil {
			return err
		}

		ctx, err := resolveRepoContext(explicitRepo)
		if err != nil {
			return err
		}
		if err := printRepoHeader(stdout, ctx); err != nil {
			return err
		}
		if _, err := state.RecordCommand(
			ctx.RepoDir,
			"d2a serve",
			state.StageServing,
			"Open the local report page and continue with d2a-report-build if presentation content needs refinement.",
			"Started serving the local d2a report surface.",
		); err != nil {
			return err
		}
		return server.Serve(ctx.RepoDir, server.DefaultAddr)
	case "status":
		explicitRepo, err := parseOptionalRepoArgs("status", args)
		if err != nil {
			return err
		}

		ctx, err := resolveRepoContext(explicitRepo)
		if err != nil {
			return err
		}
		if err := printRepoHeader(stdout, ctx); err != nil {
			return err
		}
		return printStatus(stdout, ctx.RepoDir)
	case "skill-state":
		skillName, explicitRepo, update, err := parseSkillStateArgs(args)
		if err != nil {
			return err
		}

		ctx, err := resolveRepoContext(explicitRepo)
		if err != nil {
			return err
		}
		if err := printRepoHeader(stdout, ctx); err != nil {
			return err
		}

		update.Skill = skillName
		s, err := state.RecordSkill(ctx.RepoDir, update)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(stdout, "recorded skill state for %s at stage %s\n", s.LastSkill, valueOrUnknown(s.CurrentStage))
		return err
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func printUsage(stdout io.Writer) {
	fmt.Fprint(stdout, usage)
}

func initRepoContext(target string) (repoContext, error) {
	path, err := filepath.Abs(target)
	if err != nil {
		return repoContext{}, fmt.Errorf("resolve repo path: %w", err)
	}
	return repoContext{
		Name:    filepath.Base(path),
		RepoDir: path,
		D2APath: filepath.Join(path, ".d2a"),
	}, nil
}

func resolveRepoContext(explicit string) (repoContext, error) {
	path := explicit
	if path == "" {
		wd, err := os.Getwd()
		if err != nil {
			return repoContext{}, fmt.Errorf("resolve working directory: %w", err)
		}
		path = wd
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return repoContext{}, fmt.Errorf("resolve repo path: %w", err)
	}

	d2aFile := filepath.Join(path, "LAB.md")
	if _, err := os.Stat(d2aFile); err != nil {
		if explicit == "" {
			return repoContext{}, errors.New("no active d2a workspace could be determined; run this command inside a d2a workspace or specify --repo <repo-dir>")
		}
		return repoContext{}, fmt.Errorf("specified repository is not initialized for d2a workspace: %s", path)
	}

	return repoContext{
		Name:    filepath.Base(path),
		RepoDir: path,
		D2APath: filepath.Join(path, ".d2a"),
	}, nil
}

func printRepoHeader(stdout io.Writer, ctx repoContext) error {
	_, err := fmt.Fprintf(stdout, "d2a repo: %s\nd2a repo path: %s\nd2a path: %s\n", ctx.Name, ctx.RepoDir, ctx.D2APath)
	return err
}

func printStatus(stdout io.Writer, repoRoot string) error {
	s, err := state.Load(repoRoot)
	if err != nil {
		return err
	}
	history, err := state.RecentHistory(repoRoot, 5)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(stdout,
		"current stage: %s\ncurrent skill: %s\ncurrent phase: %s\nquestion progress: %s\nlast command: %s\nlast skill: %s\nnext step: %s\nnext skill: %s\nnext file: %s\n",
		valueOrUnknown(s.CurrentStage),
		valueOrUnknown(s.CurrentSkill),
		valueOrUnknown(s.CurrentPhase),
		questionProgress(s.QuestionIndex, s.QuestionTotal),
		valueOrUnknown(s.LastCommand),
		valueOrUnknown(s.LastSkill),
		valueOrUnknown(s.NextStep),
		valueOrUnknown(s.NextSkill),
		valueOrUnknown(s.NextFile),
	); err != nil {
		return err
	}
	if stateIsChallenge(s.CurrentStage) {
		if err := printChallengeStatus(stdout, repoRoot); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(stdout, "recent history:"); err != nil {
		return err
	}

	if len(history) == 0 {
		_, err = fmt.Fprintln(stdout, "- none")
		return err
	}

	for _, event := range history {
		if _, err := fmt.Fprintf(stdout, "- [%s] %s: %s\n", event.StageAfter, event.ActorName, event.Summary); err != nil {
			return err
		}
	}
	return nil
}

func printChallengeStatus(stdout io.Writer, repoRoot string) error {
	challenge, err := state.ChallengeState(repoRoot)
	if err != nil {
		if _, writeErr := fmt.Fprintln(stdout, "challenge progress: unavailable"); writeErr != nil {
			return writeErr
		}
		return nil
	}
	_, err = fmt.Fprintf(stdout,
		"challenge progress: %s\ncurrent decision: %s\nlast challenge strength: %s\nchallenge recommendation: %s\n",
		questionProgress(challenge.DecisionIndex, challenge.DecisionTotal),
		valueOrUnknown(challenge.CurrentDecision),
		valueOrUnknown(challenge.LastStrength),
		valueOrUnknown(challenge.Recommendation),
	)
	return err
}

func valueOrUnknown(value string) string {
	if value == "" {
		return "unknown"
	}
	return value
}

func questionProgress(index, total int) string {
	if total <= 0 {
		return "none"
	}
	return fmt.Sprintf("%d/%d", index, total)
}

func stateIsChallenge(stage string) bool {
	return stage == state.StageArchitectureChallengePrepared ||
		stage == state.StageArchitectureChallengeInProgress ||
		stage == state.StageArchitectureChallengeComplete
}

func parseAnalyzeArgs(args []string) (targetRepo, explicitRepo string, err error) {
	switch len(args) {
	case 1:
		return "", "", nil
	case 2:
		if args[1] == "--repo" {
			return "", "", errors.New("analyze requires: d2a analyze [<target-repo>] [--repo <repo-dir>]")
		}
		return args[1], "", nil
	case 3:
		if args[1] != "--repo" {
			return "", "", errors.New("analyze requires: d2a analyze [<target-repo>] [--repo <repo-dir>]")
		}
		return "", args[2], nil
	case 4:
		if args[2] != "--repo" {
			return "", "", errors.New("analyze requires: d2a analyze [<target-repo>] [--repo <repo-dir>]")
		}
		return args[1], args[3], nil
	default:
		return "", "", errors.New("analyze requires: d2a analyze [<target-repo>] [--repo <repo-dir>]")
	}
}

func parseInitArgs(args []string) (targetRepoURL, language string, err error) {
	if len(args) == 2 {
		return args[1], installer.LanguageZH, nil
	}
	if len(args) == 4 && args[2] == "--lang" {
		return args[1], args[3], nil
	}
	return "", "", errors.New("init requires: d2a init <target-repo-git-url> [--lang <zh|en>]")
}

func parseOptionalRepoArgs(cmd string, args []string) (string, error) {
	if len(args) == 1 {
		return "", nil
	}
	if len(args) == 3 && args[1] == "--repo" {
		return args[2], nil
	}
	return "", fmt.Errorf("%s requires: d2a %s [--repo <repo-dir>]", cmd, cmd)
}

func parseDeriveMiniArgs(args []string) (explicitRepo, skipReason string, err error) {
	if len(args) == 1 {
		return "", "", nil
	}

	for i := 1; i < len(args); i += 2 {
		if i+1 >= len(args) {
			return "", "", errors.New("derive-mini arguments must use flag/value pairs")
		}
		switch args[i] {
		case "--repo":
			explicitRepo = args[i+1]
		case "--skip-challenge-reason":
			skipReason = args[i+1]
		default:
			return "", "", fmt.Errorf("derive-mini requires: d2a derive-mini [--repo <repo-dir>] [--skip-challenge-reason <text>]")
		}
	}

	return explicitRepo, skipReason, nil
}

func parseSkillStateArgs(args []string) (string, string, state.SkillUpdate, error) {
	if len(args) < 2 {
		return "", "", state.SkillUpdate{}, errors.New("skill-state requires a skill name")
	}

	update := state.SkillUpdate{
		QuestionIndex: -1,
		QuestionTotal: -1,
	}
	explicitRepo := ""
	skillName := args[1]

	for i := 2; i < len(args); i += 2 {
		if i+1 >= len(args) {
			return "", "", state.SkillUpdate{}, errors.New("skill-state arguments must use flag/value pairs")
		}

		flag := args[i]
		value := args[i+1]

		switch flag {
		case "--repo":
			explicitRepo = value
		case "--status":
			update.Status = value
		case "--stage":
			update.Stage = value
		case "--phase":
			update.Phase = value
		case "--question-index":
			n, err := parseNonNegativeInt(value, flag)
			if err != nil {
				return "", "", state.SkillUpdate{}, err
			}
			update.QuestionIndex = n
		case "--question-total":
			n, err := parseNonNegativeInt(value, flag)
			if err != nil {
				return "", "", state.SkillUpdate{}, err
			}
			update.QuestionTotal = n
		case "--next-step":
			update.NextStep = value
		case "--next-skill":
			update.NextSkill = value
		case "--next-file":
			update.NextFile = value
		case "--question":
			update.Question = value
		case "--answer":
			update.Answer = value
		case "--evaluation":
			update.Evaluation = value
		case "--explanation":
			update.Explanation = value
		case "--summary":
			update.Summary = value
		case "--decision":
			update.Decision = value
		case "--strength":
			update.Strength = value
		case "--recommendation":
			update.Recommendation = value
		case "--objection":
			update.Objection = value
		default:
			return "", "", state.SkillUpdate{}, fmt.Errorf("unknown skill-state flag %q", flag)
		}
	}

	return skillName, explicitRepo, update, nil
}

func parseNonNegativeInt(value, flag string) (int, error) {
	var n int
	if _, err := fmt.Sscanf(value, "%d", &n); err != nil || n < 0 {
		return 0, fmt.Errorf("%s requires a non-negative integer", flag)
	}
	return n, nil
}

func parseRepoNameFromGitURL(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", errors.New("target repo git url must not be empty")
	}
	if !isGitURL(value) {
		return "", fmt.Errorf("init requires a git repo url, got %q", raw)
	}

	path := value
	if u, err := url.Parse(value); err == nil && u.Scheme != "" {
		path = u.Path
	}
	path = strings.TrimSuffix(path, "/")
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, ".git")
	base = strings.TrimSpace(base)
	if base == "" || base == "." || base == string(filepath.Separator) {
		return "", fmt.Errorf("cannot determine repo name from git url %q", raw)
	}
	return base, nil
}

func isGitURL(value string) bool {
	if strings.Contains(value, "://") {
		return true
	}
	return strings.HasPrefix(value, "git@")
}

func ensureWorkspaceGitRepo(repoRoot string) error {
	if _, err := os.Stat(filepath.Join(repoRoot, ".git")); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("check workspace git dir: %w", err)
	}
	cmd := exec.Command("git", "init", repoRoot)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("initialize workspace git repository: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func cloneTargetRepo(targetRepoURL, repoRoot, repoName string) (string, error) {
	reposRoot := filepath.Join(repoRoot, "repos")
	if err := os.MkdirAll(reposRoot, 0o755); err != nil {
		return "", fmt.Errorf("create repos dir %s: %w", reposRoot, err)
	}
	targetPath := filepath.Join(reposRoot, repoName)
	if _, err := os.Stat(targetPath); err == nil {
		return targetPath, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("check target repo path %s: %w", targetPath, err)
	}
	cmd := exec.Command("git", "clone", "--depth", "1", targetRepoURL, targetPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("clone target repository: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return targetPath, nil
}

func writeInitTargetMetadata(repoRoot, targetRepoURL, targetRepoPath string) error {
	path := filepath.Join(repoRoot, ".d2a", "target.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create target metadata dir %s: %w", filepath.Dir(path), err)
	}
	content, err := json.MarshalIndent(initTargetMetadata{
		TargetRepo:    targetRepoPath,
		TargetRepoURL: targetRepoURL,
		RepoRoot:      repoRoot,
		D2APath:       filepath.Join(repoRoot, ".d2a"),
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal target metadata: %w", err)
	}
	content = append(content, '\n')
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write target metadata %s: %w", path, err)
	}
	return nil
}

func resolveDefaultTargetRepo(repoRoot string) (string, error) {
	path := filepath.Join(repoRoot, ".d2a", "target.json")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("target repository is not configured; run d2a init <target-repo-git-url> again or pass d2a analyze <target-repo>")
		}
		return "", fmt.Errorf("read target metadata %s: %w", path, err)
	}
	var meta initTargetMetadata
	if err := json.Unmarshal(content, &meta); err != nil {
		return "", fmt.Errorf("parse target metadata %s: %w", path, err)
	}
	if strings.TrimSpace(meta.TargetRepo) == "" {
		return "", fmt.Errorf("target metadata missing target_repo: %s", path)
	}
	return meta.TargetRepo, nil
}

func writeAgentsFile(repoRoot string) error {
	path := filepath.Join(repoRoot, "AGENTS.md")
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("check AGENTS.md %s: %w", path, err)
	}
	content := "# AGENTS\n\n## Rules\n\n1. Use `d2a-step` as the primary orchestration skill in every turn.\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write AGENTS.md %s: %w", path, err)
	}
	return nil
}
