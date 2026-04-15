package state

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	StageInitialized                     = "initialized"
	StageAnalysisPrepared                = "analysis-prepared"
	StageArchitectureInProgress          = "architecture-in-progress"
	StageArchitectureComplete            = "architecture-complete"
	StageArchitectureChallengePrepared   = "architecture-challenge-prepared"
	StageArchitectureChallengeInProgress = "architecture-challenge-in-progress"
	StageArchitectureChallengeComplete   = "architecture-challenge-complete"
	StageMiniDerivationPrepared          = "mini-derivation-prepared"
	StageMiniDesignInProgress            = "mini-design-in-progress"
	StageMiniDesignComplete              = "mini-design-complete"
	StageTestPlanPrepared                = "test-plan-prepared"
	StageTestingInProgress               = "testing-in-progress"
	StageTestingComplete                 = "testing-complete"
	StageReportPrepared                  = "report-prepared"
	StageReportReady                     = "report-ready"
	StageServing                         = "serving"
)

var stageSequence = []string{
	StageInitialized,
	StageAnalysisPrepared,
	StageArchitectureInProgress,
	StageArchitectureComplete,
	StageArchitectureChallengePrepared,
	StageArchitectureChallengeInProgress,
	StageArchitectureChallengeComplete,
	StageMiniDerivationPrepared,
	StageMiniDesignInProgress,
	StageMiniDesignComplete,
	StageTestPlanPrepared,
	StageTestingInProgress,
	StageTestingComplete,
	StageReportPrepared,
	StageReportReady,
	StageServing,
}

var stageRank = func() map[string]int {
	out := make(map[string]int, len(stageSequence))
	for i, stage := range stageSequence {
		out[stage] = i
	}
	return out
}()

type Snapshot struct {
	Version                 int    `json:"version"`
	RepoName                string `json:"repo_name"`
	RepoPath                string `json:"repo_path"`
	D2APath                 string `json:"d2a_path"`
	CurrentStage            string `json:"current_stage"`
	CurrentSkill            string `json:"current_skill"`
	CurrentPhase            string `json:"current_phase"`
	QuestionIndex           int    `json:"question_index"`
	QuestionTotal           int    `json:"question_total"`
	LastCommand             string `json:"last_command"`
	LastSkill               string `json:"last_skill"`
	NextStep                string `json:"next_step"`
	NextSkill               string `json:"next_skill"`
	NextFile                string `json:"next_file"`
	CurrentDecision         string `json:"current_decision"`
	LastChallengeStrength   string `json:"last_challenge_strength"`
	ChallengeRecommendation string `json:"challenge_recommendation"`
	UpdatedAt               string `json:"updated_at"`
}

type Event struct {
	Timestamp     string `json:"timestamp"`
	ActorType     string `json:"actor_type"`
	ActorName     string `json:"actor_name"`
	StageBefore   string `json:"stage_before"`
	StageAfter    string `json:"stage_after"`
	PhaseAfter    string `json:"phase_after,omitempty"`
	QuestionIndex int    `json:"question_index,omitempty"`
	QuestionTotal int    `json:"question_total,omitempty"`
	Summary       string `json:"summary"`
}

type ChallengeSnapshot struct {
	RepoName        string `json:"repo_name"`
	RepoPath        string `json:"repo_path"`
	D2APath         string `json:"d2a_path"`
	CurrentStage    string `json:"current_stage"`
	CurrentDecision string `json:"current_decision"`
	DecisionIndex   int    `json:"decision_index"`
	DecisionTotal   int    `json:"decision_total"`
	LastStrength    string `json:"last_strength"`
	Recommendation  string `json:"recommendation"`
	LastObjection   string `json:"last_objection"`
	SkipReason      string `json:"skip_reason,omitempty"`
	UpdatedAt       string `json:"updated_at"`
}

type ChallengeEvent struct {
	Timestamp      string `json:"timestamp"`
	Decision       string `json:"decision"`
	DecisionIndex  int    `json:"decision_index"`
	DecisionTotal  int    `json:"decision_total"`
	Strength       string `json:"strength"`
	Recommendation string `json:"recommendation"`
	Summary        string `json:"summary"`
}

type SkillUpdate struct {
	Skill          string
	Status         string
	Stage          string
	Phase          string
	QuestionIndex  int
	QuestionTotal  int
	Question       string
	Answer         string
	Evaluation     string
	Explanation    string
	NextStep       string
	NextSkill      string
	NextFile       string
	Summary        string
	Decision       string
	Strength       string
	Recommendation string
	Objection      string
}

type QASnapshot struct {
	Skill         string `json:"skill"`
	Stage         string `json:"stage"`
	Phase         string `json:"phase"`
	QuestionIndex int    `json:"question_index"`
	QuestionTotal int    `json:"question_total"`
	UpdatedAt     string `json:"updated_at"`
}

type QAEvent struct {
	Timestamp     string `json:"timestamp"`
	Skill         string `json:"skill"`
	Stage         string `json:"stage"`
	Phase         string `json:"phase"`
	QuestionIndex int    `json:"question_index"`
	QuestionTotal int    `json:"question_total"`
	Question      string `json:"question,omitempty"`
	Answer        string `json:"answer,omitempty"`
	Evaluation    string `json:"evaluation,omitempty"`
	Explanation   string `json:"explanation,omitempty"`
	Summary       string `json:"summary,omitempty"`
}

func Initialize(repoRoot string) (Snapshot, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return Snapshot{}, fmt.Errorf("resolve repo root: %w", err)
	}

	s := defaultSnapshot(repoRoot)
	s.CurrentStage = StageInitialized
	s.LastCommand = "d2a init"
	s.NextStep = "Run d2a analyze."
	s.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := writeSnapshot(s); err != nil {
		return Snapshot{}, err
	}
	if err := appendEvent(repoRoot, Event{
		Timestamp:   s.UpdatedAt,
		ActorType:   "command",
		ActorName:   "d2a init",
		StageBefore: "",
		StageAfter:  StageInitialized,
		Summary:     "Initialized d2a in the repository and created the .d2a workspace.",
	}); err != nil {
		return Snapshot{}, err
	}

	return s, nil
}

func Load(repoRoot string) (Snapshot, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return Snapshot{}, fmt.Errorf("resolve repo root: %w", err)
	}

	path := snapshotPath(repoRoot)
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Snapshot{}, fmt.Errorf("state file not found: %s", path)
		}
		return Snapshot{}, fmt.Errorf("read state file %s: %w", path, err)
	}

	var s Snapshot
	if err := json.Unmarshal(content, &s); err != nil {
		return Snapshot{}, fmt.Errorf("parse state file %s: %w", path, err)
	}
	return s, nil
}

func RecordCommand(repoRoot, command, nextStage, nextStep, summary string) (Snapshot, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return Snapshot{}, fmt.Errorf("resolve repo root: %w", err)
	}

	s, err := loadOrBootstrap(repoRoot)
	if err != nil {
		return Snapshot{}, err
	}

	stageBefore := s.CurrentStage
	if strings.TrimSpace(nextStage) != "" {
		if err := validateStageTransition(stageBefore, nextStage); err != nil {
			return Snapshot{}, err
		}
		s.CurrentStage = nextStage
	}
	s.LastCommand = command
	s.CurrentPhase = ""
	s.QuestionIndex = 0
	s.QuestionTotal = 0
	if strings.TrimSpace(nextStep) != "" {
		s.NextStep = nextStep
	}
	s.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := writeSnapshot(s); err != nil {
		return Snapshot{}, err
	}
	if err := appendEvent(repoRoot, Event{
		Timestamp:   s.UpdatedAt,
		ActorType:   "command",
		ActorName:   command,
		StageBefore: stageBefore,
		StageAfter:  s.CurrentStage,
		PhaseAfter:  s.CurrentPhase,
		Summary:     summary,
	}); err != nil {
		return Snapshot{}, err
	}

	return s, nil
}

func RecordSkill(repoRoot string, update SkillUpdate) (Snapshot, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return Snapshot{}, fmt.Errorf("resolve repo root: %w", err)
	}
	if strings.TrimSpace(update.Skill) == "" {
		return Snapshot{}, fmt.Errorf("skill name must not be empty")
	}

	s, err := loadOrBootstrap(repoRoot)
	if err != nil {
		return Snapshot{}, err
	}

	stageBefore := s.CurrentStage
	if strings.TrimSpace(update.Stage) != "" {
		if err := validateStageTransition(stageBefore, update.Stage); err != nil {
			return Snapshot{}, err
		}
		s.CurrentStage = update.Stage
	}
	s.CurrentSkill = update.Skill
	s.LastSkill = update.Skill
	if update.Phase != "" {
		s.CurrentPhase = update.Phase
	}
	if update.QuestionIndex >= 0 {
		s.QuestionIndex = update.QuestionIndex
	}
	if update.QuestionTotal >= 0 {
		s.QuestionTotal = update.QuestionTotal
	}
	if update.NextStep != "" {
		s.NextStep = update.NextStep
	}
	if update.NextSkill != "" {
		s.NextSkill = update.NextSkill
	}
	if update.NextFile != "" {
		s.NextFile = update.NextFile
	}
	if update.Decision != "" {
		s.CurrentDecision = update.Decision
	}
	if update.Strength != "" {
		s.LastChallengeStrength = update.Strength
	}
	if update.Recommendation != "" {
		s.ChallengeRecommendation = update.Recommendation
	}
	if update.Status == "completed" && update.Phase == "" {
		s.CurrentPhase = ""
	}
	if update.Status == "completed" && update.QuestionIndex == 0 && update.QuestionTotal == 0 {
		s.QuestionIndex = 0
		s.QuestionTotal = 0
	}
	s.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := writeSnapshot(s); err != nil {
		return Snapshot{}, err
	}
	if err := appendEvent(repoRoot, Event{
		Timestamp:     s.UpdatedAt,
		ActorType:     "skill",
		ActorName:     update.Skill,
		StageBefore:   stageBefore,
		StageAfter:    s.CurrentStage,
		PhaseAfter:    s.CurrentPhase,
		QuestionIndex: s.QuestionIndex,
		QuestionTotal: s.QuestionTotal,
		Summary:       update.Summary,
	}); err != nil {
		return Snapshot{}, err
	}
	if isChallengeStage(s.CurrentStage) || update.Skill == "d2a-challenge-architecture" {
		if err := writeChallengeState(s, update); err != nil {
			return Snapshot{}, err
		}
		if err := appendChallengeEvent(repoRoot, s, update); err != nil {
			return Snapshot{}, err
		}
	}
	if s.CurrentPhase == "confirmation-questions" {
		if err := writeQASnapshot(s); err != nil {
			return Snapshot{}, err
		}
		if err := appendQAEvent(s, update); err != nil {
			return Snapshot{}, err
		}
	}

	return s, nil
}

func RecentHistory(repoRoot string, limit int) ([]Event, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve repo root: %w", err)
	}

	path := historyPath(repoRoot)
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("open history file %s: %w", path, err)
	}
	defer file.Close()

	events := make([]Event, 0, limit)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var event Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return nil, fmt.Errorf("parse history line in %s: %w", path, err)
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan history file %s: %w", path, err)
	}

	if limit <= 0 || len(events) <= limit {
		return events, nil
	}
	return events[len(events)-limit:], nil
}

func loadOrBootstrap(repoRoot string) (Snapshot, error) {
	s, err := Load(repoRoot)
	if err == nil {
		return s, nil
	}

	if !strings.Contains(err.Error(), "state file not found") {
		return Snapshot{}, err
	}

	d2aLabPath := filepath.Join(repoRoot, "LAB.md")
	if _, statErr := os.Stat(d2aLabPath); statErr != nil {
		return Snapshot{}, err
	}

	s = defaultSnapshot(repoRoot)
	s.CurrentStage = StageInitialized
	s.NextStep = "Run d2a analyze."
	s.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := writeSnapshot(s); err != nil {
		return Snapshot{}, err
	}
	return s, nil
}

func defaultSnapshot(repoRoot string) Snapshot {
	return Snapshot{
		Version:  1,
		RepoName: filepath.Base(repoRoot),
		RepoPath: repoRoot,
		D2APath:  filepath.Join(repoRoot, ".d2a"),
	}
}

func writeSnapshot(s Snapshot) error {
	if err := os.MkdirAll(s.D2APath, 0o755); err != nil {
		return fmt.Errorf("create d2a dir %s: %w", s.D2APath, err)
	}

	content, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state snapshot: %w", err)
	}
	content = append(content, '\n')

	path := snapshotPath(s.RepoPath)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write state file %s: %w", path, err)
	}
	return nil
}

func appendEvent(repoRoot string, event Event) error {
	path := historyPath(repoRoot)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create history dir %s: %w", filepath.Dir(path), err)
	}

	content, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal history event: %w", err)
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open history file %s: %w", path, err)
	}
	defer file.Close()

	if _, err := file.Write(append(content, '\n')); err != nil {
		return fmt.Errorf("append history file %s: %w", path, err)
	}
	return nil
}

func snapshotPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".d2a", "state.json")
}

func historyPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".d2a", "history.jsonl")
}

func ChallengeState(repoRoot string) (ChallengeSnapshot, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return ChallengeSnapshot{}, fmt.Errorf("resolve repo root: %w", err)
	}

	path := challengeStatePath(repoRoot)
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ChallengeSnapshot{}, fmt.Errorf("challenge state file not found: %s", path)
		}
		return ChallengeSnapshot{}, fmt.Errorf("read challenge state file %s: %w", path, err)
	}

	var s ChallengeSnapshot
	if err := json.Unmarshal(content, &s); err != nil {
		return ChallengeSnapshot{}, fmt.Errorf("parse challenge state file %s: %w", path, err)
	}
	return s, nil
}

func SkipChallenge(repoRoot, reason string) (Snapshot, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return Snapshot{}, fmt.Errorf("resolve repo root: %w", err)
	}
	if strings.TrimSpace(reason) == "" {
		return Snapshot{}, fmt.Errorf("skip reason must not be empty")
	}

	s, err := loadOrBootstrap(repoRoot)
	if err != nil {
		return Snapshot{}, err
	}

	stageBefore := s.CurrentStage
	s.CurrentStage = StageArchitectureChallengeComplete
	s.CurrentSkill = "d2a-challenge-architecture"
	s.LastSkill = "d2a-challenge-architecture"
	s.CurrentPhase = "challenge-skipped"
	s.QuestionIndex = 0
	s.QuestionTotal = 0
	s.CurrentDecision = ""
	s.LastChallengeStrength = ""
	s.ChallengeRecommendation = "skipped"
	s.NextStep = "Proceed to d2a derive-mini."
	s.NextSkill = "d2a-mini-1-scope"
	s.NextFile = "docs/implementation/00_mini_scope.md"
	s.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := writeSnapshot(s); err != nil {
		return Snapshot{}, err
	}
	if err := appendEvent(repoRoot, Event{
		Timestamp:   s.UpdatedAt,
		ActorType:   "command",
		ActorName:   "d2a derive-mini --skip-challenge-reason",
		StageBefore: stageBefore,
		StageAfter:  s.CurrentStage,
		PhaseAfter:  s.CurrentPhase,
		Summary:     "Skipped architecture challenge phase: " + reason,
	}); err != nil {
		return Snapshot{}, err
	}
	if err := writeChallengeSkipState(s, reason); err != nil {
		return Snapshot{}, err
	}
	if err := appendChallengeSkipEvent(repoRoot, s, reason); err != nil {
		return Snapshot{}, err
	}

	return s, nil
}

func writeChallengeState(snapshot Snapshot, update SkillUpdate) error {
	challenge := ChallengeSnapshot{
		RepoName:        snapshot.RepoName,
		RepoPath:        snapshot.RepoPath,
		D2APath:         snapshot.D2APath,
		CurrentStage:    snapshot.CurrentStage,
		CurrentDecision: snapshot.CurrentDecision,
		DecisionIndex:   snapshot.QuestionIndex,
		DecisionTotal:   snapshot.QuestionTotal,
		LastStrength:    snapshot.LastChallengeStrength,
		Recommendation:  snapshot.ChallengeRecommendation,
		LastObjection:   update.Objection,
		UpdatedAt:       snapshot.UpdatedAt,
	}

	content, err := json.MarshalIndent(challenge, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal challenge snapshot: %w", err)
	}
	content = append(content, '\n')

	path := challengeStatePath(snapshot.RepoPath)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write challenge state file %s: %w", path, err)
	}
	return nil
}

func writeChallengeSkipState(snapshot Snapshot, reason string) error {
	challenge := ChallengeSnapshot{
		RepoName:        snapshot.RepoName,
		RepoPath:        snapshot.RepoPath,
		D2APath:         snapshot.D2APath,
		CurrentStage:    snapshot.CurrentStage,
		CurrentDecision: "",
		DecisionIndex:   0,
		DecisionTotal:   0,
		LastStrength:    "",
		Recommendation:  "skipped",
		LastObjection:   "",
		SkipReason:      reason,
		UpdatedAt:       snapshot.UpdatedAt,
	}

	content, err := json.MarshalIndent(challenge, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal challenge skip snapshot: %w", err)
	}
	content = append(content, '\n')

	path := challengeStatePath(snapshot.RepoPath)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write challenge state file %s: %w", path, err)
	}
	return nil
}

func appendChallengeEvent(repoRoot string, snapshot Snapshot, update SkillUpdate) error {
	event := ChallengeEvent{
		Timestamp:      snapshot.UpdatedAt,
		Decision:       snapshot.CurrentDecision,
		DecisionIndex:  snapshot.QuestionIndex,
		DecisionTotal:  snapshot.QuestionTotal,
		Strength:       snapshot.LastChallengeStrength,
		Recommendation: snapshot.ChallengeRecommendation,
		Summary:        update.Summary,
	}

	content, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal challenge event: %w", err)
	}

	path := challengeLogPath(repoRoot)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open challenge log file %s: %w", path, err)
	}
	defer file.Close()

	if _, err := file.Write(append(content, '\n')); err != nil {
		return fmt.Errorf("append challenge log file %s: %w", path, err)
	}
	return nil
}

func appendChallengeSkipEvent(repoRoot string, snapshot Snapshot, reason string) error {
	event := ChallengeEvent{
		Timestamp:      snapshot.UpdatedAt,
		Decision:       "",
		DecisionIndex:  0,
		DecisionTotal:  0,
		Strength:       "",
		Recommendation: "skipped",
		Summary:        "Skipped architecture challenge phase: " + reason,
	}

	content, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal challenge skip event: %w", err)
	}

	path := challengeLogPath(snapshot.RepoPath)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open challenge log file %s: %w", path, err)
	}
	defer file.Close()

	if _, err := file.Write(append(content, '\n')); err != nil {
		return fmt.Errorf("append challenge log file %s: %w", path, err)
	}
	return nil
}

func challengeStatePath(repoRoot string) string {
	return filepath.Join(repoRoot, ".d2a", "challenge.json")
}

func challengeLogPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".d2a", "challenge_log.jsonl")
}

func writeQASnapshot(snapshot Snapshot) error {
	qa := QASnapshot{
		Skill:         snapshot.CurrentSkill,
		Stage:         snapshot.CurrentStage,
		Phase:         snapshot.CurrentPhase,
		QuestionIndex: snapshot.QuestionIndex,
		QuestionTotal: snapshot.QuestionTotal,
		UpdatedAt:     snapshot.UpdatedAt,
	}
	content, err := json.MarshalIndent(qa, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal qa snapshot: %w", err)
	}
	content = append(content, '\n')
	path := qaStatePath(snapshot.RepoPath, snapshot.CurrentSkill)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create qa dir %s: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write qa snapshot %s: %w", path, err)
	}
	return nil
}

func appendQAEvent(snapshot Snapshot, update SkillUpdate) error {
	event := QAEvent{
		Timestamp:     snapshot.UpdatedAt,
		Skill:         snapshot.CurrentSkill,
		Stage:         snapshot.CurrentStage,
		Phase:         snapshot.CurrentPhase,
		QuestionIndex: snapshot.QuestionIndex,
		QuestionTotal: snapshot.QuestionTotal,
		Question:      update.Question,
		Answer:        update.Answer,
		Evaluation:    update.Evaluation,
		Explanation:   update.Explanation,
		Summary:       update.Summary,
	}
	content, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal qa event: %w", err)
	}
	path := qaLogPath(snapshot.RepoPath, snapshot.CurrentSkill)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create qa dir %s: %w", filepath.Dir(path), err)
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open qa log file %s: %w", path, err)
	}
	defer file.Close()
	if _, err := file.Write(append(content, '\n')); err != nil {
		return fmt.Errorf("append qa log file %s: %w", path, err)
	}
	return nil
}

func qaStatePath(repoRoot, skill string) string {
	return filepath.Join(repoRoot, ".d2a", "qa", skill+".json")
}

func qaLogPath(repoRoot, skill string) string {
	return filepath.Join(repoRoot, ".d2a", "qa", skill+".jsonl")
}

func isChallengeStage(stage string) bool {
	return stage == StageArchitectureChallengePrepared ||
		stage == StageArchitectureChallengeInProgress ||
		stage == StageArchitectureChallengeComplete
}

func validateStageTransition(stageBefore, stageAfter string) error {
	afterRank, ok := stageRank[stageAfter]
	if !ok {
		return fmt.Errorf("unknown stage: %s", stageAfter)
	}

	if stageBefore == "" {
		return nil
	}

	beforeRank, ok := stageRank[stageBefore]
	if !ok {
		return nil
	}
	if afterRank < beforeRank {
		return fmt.Errorf("invalid stage regression: %s -> %s", stageBefore, stageAfter)
	}

	return nil
}
