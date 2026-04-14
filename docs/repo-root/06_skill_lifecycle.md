# Skill Lifecycle Rules

## Goal

Require every `d2a-*` skill to behave like a state-machine participant.

Each skill should also behave like a teaching loop:

1. generate analysis
2. verify learner understanding
3. challenge and defend the architecture when required
4. summarize progress

## Required Start Behavior

At the start of every skill run, the skill must print:

1. repository name
2. repository path
3. d2a path
4. current stage
5. whole-stage flow
6. current recommended step
7. current phase inside the skill

Example shape:

```text
d2a repo: n8n
d2a repo path: /abs/path/to/n8n
d2a path: /abs/path/to/n8n/.d2a
d2a stage: architecture-in-progress
d2a flow: initialized -> analysis-prepared -> architecture-in-progress -> architecture-complete -> architecture-challenge-in-progress -> architecture-challenge-complete -> mini-derivation-prepared -> mini-design-complete -> test-plan-prepared -> testing-in-progress -> testing-complete -> report-ready
d2a phase: analysis-generation
d2a next step: fill docs/architecture/02_driver.md
```

## Required Middle Behavior

When the skill enters the confirmation-question phase, it must:

1. switch the printed phase to `confirmation-questions`
2. ask one multiple-choice question at a time
3. show question progress such as `2/5`
4. evaluate the answer before moving to the next question

When the skill enters the challenge phase, it must:

1. switch the printed phase to `challenge-dialogue`
2. show challenge progress such as `3/6`
3. show the current architecture decision under challenge
4. answer the learner objection
5. record whether the objection is strong, partial, or weak

## Required End Behavior

At the end of every skill run, the AI coding tool must summarize:

1. what was completed
2. what stage has now been reached
3. what the next step should be
4. which file or command should be used next
5. a short `理解度打分`

## Failure Rule

If the current state cannot be determined, the skill must stop and ask the user to restore or confirm context.

The skill must not continue with an unknown stage.

## Persistence Rule

Every meaningful skill run should update `.d2a/state.json` and append an event to `.d2a/history.jsonl`.

This includes:

- phase switches
- question progress
- final comprehension evaluation

## Continuation Rule

At the end of every learner-facing turn, the skill should:

1. persist turn progress via `d2a skill-state`
2. instruct the learner to continue with `$d2a-step`
