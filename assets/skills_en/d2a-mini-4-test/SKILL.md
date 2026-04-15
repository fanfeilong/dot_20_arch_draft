---
name: d2a-mini-4-test
description: Built-in d2a skill for d2a-mini-4-test stage guidance and state updates.
---

# d2a-mini-4-test

## Goal

Turn the mini implementation into a testable teaching artifact through a small set of integration tests, then verify that the learner actually understands why those tests are sufficient.

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
   - Test design should prioritize the provider's minimal verification contract.
2. Timebox Gate (cost)
   - Declare test-stage budget (recommended: first end-to-end scenario within 10 minutes).
   - If it cannot fit, keep only one success signal and one failure signal.
3. Intent Gate (alignment)
   - Tests validate mini intent realization, not broad business coverage.
   - Verify observability of object/state/cooperation-chain anchors.

## Phase 1: Analysis Generation

1. After context is confirmed, call:

   `d2a skill-state d2a-mini-4-test --status started --stage testing-in-progress --phase analysis-generation --next-step "Create the first integration test." --next-skill "d2a-report-build" --next-file ".d2a/tests/README.md" --summary "Started mini-test work."`

2. Execute the three gates first and output gate conclusions (provider match, timebox, intent anchors).
3. Treat this skill as the user-facing entry for the testing stage inside Codex.
4. If test task files have not yet been prepared, call `d2a test-mini` before writing tests.
5. Read:
   - `.d2a/docs/implementation/03_test_plan.md`
   - `.d2a/tests/README.md`
   - `.d2a/tests/01_integration_tasks.md`
   - `.d2a/src/ARCHITECTURE.md`
6. Start with one end-to-end test for the first runnable slice.
7. Prefer observable behavior over internal unit details.
8. Add only the next most useful failure case after the first scenario is clear.
9. Keep tests aligned with the architecture idea being demonstrated.
10. Make explicit:
   - the observable success signal
   - the observable failure signal
11. When the testing pass is stable, call:

   `d2a skill-state d2a-mini-4-test --status progress --stage testing-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first mini-test confirmation question." --next-skill "d2a-mini-4-test" --next-file ".d2a/tests/01_integration_tasks.md" --summary "Mini-test work complete; moving into confirmation questions."`

## Phase 2: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - provider test-contract adherence
   - minimal test set under timebox
   - observable success/failure signals
   - intent-anchor validation quality
3. Ask one question per turn.
4. Before each question, keep the same envelope format and map fields as follows:
   - The `[d2a]` line must include current stage, phase, and `N/<total>` progress.
   - The `[next]` line should point to the post-question next skill/file.
5. Before asking question `N`, call:

   `d2a skill-state d2a-mini-4-test --status progress --stage testing-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue mini-test confirmation questions." --next-skill "d2a-report-build" --next-file ".d2a/report/index.md" --summary "Mini-test confirmation question <N> is active."`

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

    `d2a skill-state d2a-mini-4-test --status completed --stage testing-complete --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-report-build." --next-skill "d2a-report-build" --next-file ".d2a/report/index.md" --summary "Completed mini-test confirmation questions."`

11. Confirmation-question prompts, learner answers, evaluations, and explanations must be written to `.d2a/qa/<skill>.jsonl`, and must not be written into `tests/*` or `report/*`.

## Turn-End Continuation Rule

1. Before ending each turn, call `d2a skill-state` to persist the current phase and progress.
2. End the reply with: `Continue with $d2a-step`.


## Output (Artifacts)

- First integration test (written to `tests/*`).
- Next integration scenario (written to `tests/*`).
- Observable success/failure signals (written to `tests/01_integration_tasks.md`).

## Persistence (.d2a)

- Provider/Timebox/Intent gate decisions -> `.d2a/mini_gate/d2a-mini-4-test.json`.
- Confirmation prompts, learner answers, evaluations, explanations -> `.d2a/qa/<skill>.jsonl`.
- Short comprehension summary and `Comprehension Score` -> `.d2a/qa/<skill>.jsonl`.
