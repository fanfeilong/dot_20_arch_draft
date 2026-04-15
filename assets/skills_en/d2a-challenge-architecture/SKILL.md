---
name: d2a-challenge-architecture
description: Built-in d2a skill for d2a-challenge-architecture stage guidance and state updates.
---

# d2a-challenge-architecture

## Goal

Pressure-test the completed architecture through repeated learner objections, then record the challenge outcome without mutating the architecture docs by default.

## Language Rule

All user-facing text must be in English.

## Required Output Format (Same Envelope as d2a-step)

All non-`d2a-step` skills must reuse the same envelope format as `d2a-step`:

```text
==================================================
[<Layer-1> <N/Total> | <Layer-2>] <repo>
next: <next_skill> -> <next_file>
==================================================

<body of this skill>

--------------------------------------------------
done: <what this turn completed>
state: <current skeleton position> -> <next skeleton position> · Continue with $d2a-step
--------------------------------------------------
```

Rules:

1. Keep the opening to exactly two lines between separators.
2. Keep the ending to exactly two lines between separators.
3. Do not print the old multi-line start header list anymore (repo/path/stage/flow as separate lines).
4. If `next_file` is unknown, print `unknown` and keep the same two-line shape.

If the active repository is unknown, stop and ask the user which repository should be used.

## Human In Loop Marker Rule

When the current turn asks the user a question and waits for user input, the last line of the reply body must append:

`[human_in_loop]`

## Inputs

Read these files before starting:

- `.d2a/docs/architecture/00_overview.md`
- `.d2a/docs/architecture/01_boundary.md`
- `.d2a/docs/architecture/02_driver.md`
- `.d2a/docs/architecture/03_core_objects.md`
- `.d2a/docs/architecture/04_state_evolution.md`
- `.d2a/docs/architecture/05_cooperation.md`
- `.d2a/docs/architecture/06_constraints.md`
- `.d2a/docs/architecture/99_code_map.md`

## Phase 1: Challenge Preparation

1. After context is confirmed, call:

   `d2a skill-state d2a-challenge-architecture --status started --stage architecture-challenge-prepared --phase challenge-preparation --next-step "Prepare the six architecture decisions for challenge dialogue." --next-skill "d2a-challenge-architecture" --next-file ".d2a/docs/architecture/00_overview.md" --summary "Started architecture challenge preparation."`

2. Extract these `6` architecture decisions to challenge:
   - system boundary
   - primary driver
   - core objects
   - state evolution
   - cooperation pattern
   - dominant constraint / tradeoff
3. Do not rewrite any architecture file in this phase.
4. If needed, create a temporary challenge checklist in the conversation, but keep the architecture docs unchanged.
5. When the challenge set is ready, call:

   `d2a skill-state d2a-challenge-architecture --status progress --stage architecture-challenge-in-progress --phase challenge-dialogue --question-index 0 --question-total 6 --next-step "Start the first architecture challenge round." --next-skill "d2a-challenge-architecture" --next-file ".d2a/docs/architecture/00_overview.md" --summary "Architecture challenge set prepared; moving into challenge dialogue."`

## Phase 2: Challenge Dialogue

1. Run `6` rounds, one per architecture decision.
2. The order should be:
   - boundary
   - driver
   - core objects
   - state evolution
   - cooperation
   - dominant constraint
3. Each round must keep the same envelope format and map fields as follows:
   - The `[d2a]` line must include `architecture-challenge-in-progress`, `challenge-dialogue`, and `N/6` progress.
   - The `[next]` line should point to the post-challenge continuation target.
4. Before challenge round `N`, call:

   `d2a skill-state d2a-challenge-architecture --status progress --stage architecture-challenge-in-progress --phase challenge-dialogue --question-index <N> --question-total 6 --next-step "Continue architecture challenge dialogue." --next-skill "d2a-mini-1-scope" --next-file ".d2a/docs/implementation/00_mini_scope.md" --summary "Architecture challenge round <N> is active."`

5. In each round:
   - present one architecture decision
   - ask the learner to challenge it
   - accept the learner objection
   - answer the objection
   - classify the objection as `strong`, `partial`, or `weak`
   - explain the classification briefly
6. The AI must not directly revise the architecture docs during challenge dialogue.
7. If a challenge is strong, mark it for later review rather than silently editing architecture output.
8. Continue even when the learner objection is weak.

## Phase 3: Challenge Wrap-Up

1. After round 6:
   - summarize which objections were strong, partial, or weak
   - list unresolved questions if any
   - give one recommendation:
     - `proceed`
     - `review`
     - `revisit architecture`
2. Keep the architecture docs unchanged during this wrap-up.
3. At the end of the challenge phase, call:

   `d2a skill-state d2a-challenge-architecture --status completed --stage architecture-challenge-complete --phase challenge-dialogue --question-index 6 --question-total 6 --next-step "Proceed to d2a-mini-1-scope unless a review is required." --next-skill "d2a-mini-1-scope" --next-file ".d2a/docs/implementation/00_mini_scope.md" --summary "Completed architecture challenge phase."`

## Turn-End Continuation Rule

1. Before ending each turn, call `d2a skill-state` to persist the current phase and progress.
2. End the reply with: `Continue with $d2a-step`.


## Output (Artifacts)

- No new `docs/architecture/*` artifact file is created in this phase.

## Persistence (.d2a)

- Challenge rounds, objection strength, unresolved questions, and recommendation -> `.d2a/challenge.json` and `.d2a/challenge_log.jsonl`.
