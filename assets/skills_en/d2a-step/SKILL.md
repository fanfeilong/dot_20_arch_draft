---
name: d2a-step
description: State-driven orchestrator skill that resumes from .d2a state and routes to the correct next d2a sub-skill.
---

# d2a-step

## Goal

Provide one single continuation entry for the user. Resume from persisted `.d2a` state and drive the correct next skill without requiring manual skill selection.

## Language Rule

All user-facing text must be in English.

## Required Output Format

Every `d2a-step` reply must use this compact shape:

```text
==================================================
[<Layer-1> <N/Total> | <Layer-2>] <repo>
next: <next_skill> -> <next_file>
==================================================

<short body, ideally 2-4 lines>

--------------------------------------------------
done: <what this turn completed>
state: <current skeleton position> -> <next skeleton position> · Continue with $d2a-step
--------------------------------------------------
```

Rules:

1. Keep the opening to exactly two lines between separators.
2. Keep the ending to exactly two lines between separators.
3. Keep the body concise and operational; avoid long raw state dumps unless debugging is explicitly requested.
4. If `next_file` is unknown, print `unknown` and keep the same two-line shape.

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

### Two-Layer Skeleton Mapping (English Display)

1. Layer-1:
   - `Analysis`
   - `Implementation`
   - `Report`
2. Layer-2:
   - Analysis: `Boundary` / `Driver` / `Core Objects` / `State Machine` / `Module Cooperation` / `Constraints`
   - Implementation: `Minimal Scope` / `Minimal Design` / `Minimal Build` / `Minimal Test`
   - Report: `Report Build`
3. Header/footer skeleton positions should use these labels instead of raw stage ids.

1. If `current_phase` is `confirmation-questions` or `challenge-dialogue` and `question_index < question_total`, resume `current_skill`.
2. Otherwise route by stage:
   - `initialized` or `analysis-prepared` or `architecture-in-progress` -> `d2a-arch-1-project-scope`
   - `architecture-complete` or `architecture-challenge-prepared` or `architecture-challenge-in-progress` -> `d2a-challenge-architecture`
   - `architecture-challenge-complete` or `mini-derivation-prepared` -> `d2a-mini-1-scope`
   - `mini-design-in-progress` -> `d2a-mini-2-design`
   - `mini-design-complete` -> `d2a-mini-3-build`
   - `test-plan-prepared` or `testing-in-progress` -> `d2a-mini-4-test`
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

5. End with: `Continue with $d2a-step`.

## Output

- Current recovered state snapshot
- Selected next skill
- Reason for routing decision
- Next file to continue
- Continuation hint using `d2a-step`
