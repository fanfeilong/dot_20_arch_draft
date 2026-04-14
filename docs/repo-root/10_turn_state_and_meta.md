# Turn State And Meta Prompts

## Goal

Ensure that every question turn preserves enough context to avoid state loss.

## Required Per-Turn Header

Every question turn should begin with a compact status header.

Suggested shape:

```text
d2a repo: n8n
d2a repo path: /abs/path/to/n8n
d2a path: /abs/path/to/n8n/.d2a
d2a stage: architecture-in-progress
d2a flow: initialized -> analysis-prepared -> architecture-in-progress -> architecture-complete -> mini-derivation-prepared -> mini-design-complete -> test-plan-prepared -> testing-in-progress -> testing-complete -> report-ready
d2a question progress: 2/5
d2a next step after questions: update docs/architecture/02_driver.md and continue to docs/architecture/03_core_objects.md
```

## Turn-End Rule

After each question turn, the system should preserve:

- current question index
- total question count
- current stage
- next pending question
- next workflow step after the confirmation phase

This information should be recoverable from `.d2a/state.json` and recent history.

## Phase-End Rule

After the final question in a confirmation set:

1. output the short comprehension evaluation
2. update the stage machine
3. record the next recommended skill or file

## Persistence Rule

The state store should track both:

- stage progress
- question progress inside a stage

This prevents loss of context when the confirmation phase spans multiple interactions.

## Status

`.d2a/state.json` already includes the per-turn fields used by two-phase skills:

- `current_phase`
- `question_index`
- `question_total`
- `current_skill`
- `next_skill`
- `next_file`
