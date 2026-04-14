# d2a State Machine

## Goal

Treat the full `d2a` workflow as a persistent state machine stored under `.d2a/`.

The purpose is to avoid state loss across:

- multiple Codex sessions
- interrupted skill runs
- partial command execution
- handoff between skills

## Core Principle

`d2a` is not only a file generator. It is also a state tracker.

At any moment, `.d2a/` should answer:

1. what repository is being processed
2. what stage has been reached
3. which commands have already run
4. which `d2a-*` skills have already been used
5. what the recommended next step is

## State Machine Stages

The first version of the state machine should include these stages:

1. `initialized`
2. `analysis-prepared`
3. `architecture-in-progress`
4. `architecture-complete`
5. `architecture-challenge-prepared`
6. `architecture-challenge-in-progress`
7. `architecture-challenge-complete`
8. `mini-derivation-prepared`
9. `mini-design-in-progress`
10. `mini-design-complete`
11. `test-plan-prepared`
12. `testing-in-progress`
13. `testing-complete`
14. `report-prepared`
15. `report-ready`
16. `serving`

These stages may later be refined, but the key requirement is that progress is explicit and persisted.

## Transition Rule

Every command and every skill should either:

- keep the current stage
- move the state machine to a later valid stage

No command or skill should silently skip state recording.

The implementation enforces:

- known-stage validation
- non-regression transitions (cannot move backward)

## Intra-Stage Progress

For skills that include a confirmation-question phase, the state machine should also track sub-progress inside the current stage.

Examples:

- current phase: `analysis-generation`
- current phase: `confirmation-questions`
- question index
- total question count

This is required because a single stage may span multiple conversational turns.
