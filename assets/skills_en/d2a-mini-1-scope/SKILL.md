---
name: d2a-mini-1-scope
description: Built-in d2a skill for d2a-mini-1-scope stage guidance and state updates.
---

# d2a-mini-1-scope

## Goal

Choose the smallest runnable slice that preserves the target project's core architecture idea, then verify that the learner actually understands the scope choice.

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

## Three Mini Fast-Track Gates

In a 1-hour talk setting, the mini stage must pass these gates before implementation detail:

1. Provider Gate (stack)
   - Detect the target stack and match a built-in provider first.
   - If matched, use the provider's minimal implementation plan by default.
   - If not matched, use a generic minimal fallback (single happy path, low dependency, runnable).
2. Timebox Gate (cost)
   - Declare a mini-build budget (recommended: 20 minutes).
   - If the scope cannot fit, reduce scope immediately; do not expand features.
3. Intent Gate (alignment)
   - Mini is for intent realization, not another business-domain analysis pass.
   - Must stay aligned with arch anchors: core object, state progression, cooperation chain.

## Phase 1: Analysis Generation

1. After context is confirmed, call:

   `d2a skill-state d2a-mini-1-scope --status started --stage mini-derivation-prepared --phase analysis-generation --next-step "Choose the single architecture idea to preserve." --next-skill "d2a-mini-2-design" --next-file ".d2a/docs/implementation/00_mini_scope.md" --summary "Started mini-scope derivation."`

2. Treat this skill as the user-facing entry for the mini-implementation stage inside Codex.
3. Execute the three gates first and output gate conclusions (provider match, timebox, intent anchors).
4. If implementation planning files have not yet been prepared, call `d2a derive-mini` before producing content.
5. Read:
   - `.d2a/docs/architecture/00_overview.md`
   - `.d2a/docs/architecture/02_driver.md`
   - `.d2a/docs/architecture/03_core_objects.md`
   - `.d2a/docs/architecture/04_state_evolution.md`
   - `.d2a/docs/architecture/05_cooperation.md`
   - `.d2a/docs/architecture/06_constraints.md`
6. Write the result into `.d2a/docs/implementation/00_mini_scope.md`.
7. Answer these atomic questions:
   - What single architecture idea must be preserved?
   - Which runnable 20 percent slice is enough to demonstrate it?
   - What will the mini version intentionally omit?
   - Which stack should stay aligned with the original project?
8. After the first draft, force three correction passes:
   - compression pass
   - de-jargon pass
   - conversational simplification pass
9. Keep the scope small enough to support one first runnable slice.
10. When the analysis draft is stable, call:

   `d2a skill-state d2a-mini-1-scope --status progress --stage mini-derivation-prepared --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first mini-scope confirmation question." --next-skill "d2a-mini-1-scope" --next-file ".d2a/docs/implementation/00_mini_scope.md" --summary "Mini-scope analysis complete; moving into confirmation questions."`

## Phase 2: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - provider choice and stack fitness
   - minimal runnable slice under timebox
   - intent-anchor alignment (object/state/chain)
   - intentional omissions and scope control
3. Ask one question per turn.
4. Before each question, keep the same envelope format and map fields as follows:
   - The `[d2a]` line must include current stage, phase, and `N/<total>` progress.
   - The `[next]` line should point to the post-question next skill/file.
5. Before asking question `N`, call:

   `d2a skill-state d2a-mini-1-scope --status progress --stage mini-derivation-prepared --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue mini-scope confirmation questions." --next-skill "d2a-mini-2-design" --next-file ".d2a/docs/implementation/01_mini_design.md" --summary "Mini-scope confirmation question <N> is active."`

6. Present one question with multiple choices.
7. Wait for the learner answer.
8. After the learner answer:
   - say whether the answer is correct, partially correct, or incorrect
   - give one short explanation
   - continue to the next question even if the answer is wrong
9. After question 4 is evaluated:
   - output a short recap
   - output a `Comprehension Score`
   - keep the `Comprehension Score` under 100 Chinese characters
10. At the end of the confirmation phase, call:

    `d2a skill-state d2a-mini-1-scope --status completed --stage mini-derivation-prepared --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-mini-2-design." --next-skill "d2a-mini-2-design" --next-file ".d2a/docs/implementation/01_mini_design.md" --summary "Completed mini-scope confirmation questions."`

11. Confirmation-question prompts, learner answers, evaluations, and explanations must be written to `.d2a/qa/<skill>.jsonl`, and must not be written into `docs/implementation/*.md` or `src/*`.

## Turn-End Continuation Rule

1. Before ending each turn, call `d2a skill-state` to persist the current phase and progress.
2. End the reply with: `Continue with $d2a-step`.


## Output (Artifacts)

- Preserved architecture intent.
- Runnable 20% slice definition.
- Intentional omissions.
- Target stack.
- Written to `docs/implementation/00_mini_scope.md`.

## Persistence (.d2a)

- Provider/Timebox/Intent gate decisions -> `.d2a/mini_gate/d2a-mini-1-scope.json`.
- Confirmation prompts, learner answers, evaluations, explanations -> `.d2a/qa/<skill>.jsonl`.
- Short comprehension summary and `Comprehension Score` -> `.d2a/qa/<skill>.jsonl`.
