---
name: d2a-step
description: State-driven orchestrator skill that resumes from .d2a state and routes to the correct next d2a sub-skill.
---

# d2a-step

## Goal

Provide one single continuation entry for the user. Resume from persisted `.d2a` state and drive the correct next skill without requiring manual skill selection.

## Required Start Header

Always print this header first:

- `d2a repo: ...`
- `d2a repo path: ...`
- `d2a path: ...`
- `d2a stage: ...`
- `d2a phase: ...`
- `d2a progress: ...`
- `d2a next step: ...`

If the active repository is unknown, stop and ask the user which repository should be used.

## State Recovery

1. Call `d2a status` or read `.d2a/state.json` to get:
   - `current_stage`
   - `current_phase`
   - `question_index`
   - `question_total`
   - `current_skill`
   - `next_skill`
   - `next_file`
2. Read recent `.d2a/history.jsonl` when resumption details are unclear.
3. If `.d2a/state.json` is missing, tell the user to run `d2a init` first and stop.

## Routing Rules

1. If `current_phase` is `confirmation-questions` or `challenge-dialogue` and `question_index < question_total`, resume `current_skill`.
2. Otherwise route by stage:
   - `initialized` or `analysis-prepared` or `architecture-in-progress` -> `d2a-architecture-walkthrough`
   - `architecture-complete` or `architecture-challenge-prepared` or `architecture-challenge-in-progress` -> `d2a-challenge-architecture`
   - `architecture-challenge-complete` or `mini-derivation-prepared` -> `d2a-mini-scope`
   - `mini-design-in-progress` -> `d2a-mini-design`
   - `mini-design-complete` -> `d2a-mini-build`
   - `test-plan-prepared` or `testing-in-progress` -> `d2a-mini-test`
   - `testing-complete` or `report-prepared` or `report-ready` -> `d2a-report-build`
   - `serving` -> `d2a-status`
3. If `next_skill` exists and does not conflict with the stage routing, prefer `next_skill`.
4. If routing is ambiguous, ask one short clarification question and do not guess.

## Step Execution

1. Before handing off, call:

   `d2a skill-state d2a-step --status started --phase analysis-generation --next-step "Resume from persisted d2a state and route to the next skill." --summary "Started d2a-step orchestration."`

2. Tell the user exactly which skill to run now and why (stage + phase evidence).
3. If needed, tell the user the exact file to continue from (`next_file`).
4. After emitting the step decision, call:

   `d2a skill-state d2a-step --status completed --phase analysis-generation --next-step "Continue with the routed skill." --next-skill "<ROUTED_SKILL>" --next-file "<ROUTED_FILE>" --summary "d2a-step routed to <ROUTED_SKILL> based on persisted state."`

5. End with: `继续请使用 $d2a-step`.

## Output

- Current recovered state snapshot
- Selected next skill
- Reason for routing decision
- Next file to continue
- Continuation hint using `d2a-step`
