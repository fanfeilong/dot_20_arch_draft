---
name: d2a-mini-test
description: Built-in d2a skill for d2a-mini-test stage guidance and state updates.
---

# d2a-mini-test

## Goal

Turn the mini implementation into a testable teaching artifact through a small set of integration tests, then verify that the learner actually understands why those tests are sufficient.

## Required Start Header

Always print this header first:

- `d2a repo: ...`
- `d2a repo path: ...`
- `d2a path: ...`
- `d2a stage: ...`
- `d2a flow: initialized -> analysis-prepared -> architecture-in-progress -> architecture-complete -> architecture-challenge-prepared -> architecture-challenge-in-progress -> architecture-challenge-complete -> mini-derivation-prepared -> mini-design-in-progress -> mini-design-complete -> test-plan-prepared -> testing-in-progress -> testing-complete -> report-ready`
- `d2a phase: ...`
- `d2a next step: ...`

If the active repository is unknown, stop and ask the user which repository should be used.

## Phase 1: Analysis Generation

1. After context is confirmed, call:

   `d2a skill-state d2a-mini-test --status started --stage testing-in-progress --phase analysis-generation --next-step "Create the first integration test." --next-skill "d2a-report-build" --next-file ".d2a/tests/README.md" --summary "Started mini-test work."`

2. Treat this skill as the user-facing entry for the testing stage inside Codex.
3. If test task files have not yet been prepared, call `d2a test-mini` before writing tests.
4. Read:
   - `.d2a/docs/implementation/03_test_plan.md`
   - `.d2a/tests/README.md`
   - `.d2a/tests/01_integration_tasks.md`
   - `.d2a/src/ARCHITECTURE.md`
5. Start with one end-to-end test for the first runnable slice.
6. Prefer observable behavior over internal unit details.
7. Add only the next most useful failure case after the first scenario is clear.
8. Keep tests aligned with the architecture idea being demonstrated.
9. Make explicit:
   - the observable success signal
   - the observable failure signal
10. When the testing pass is stable, call:

   `d2a skill-state d2a-mini-test --status progress --stage testing-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first mini-test confirmation question." --next-skill "d2a-mini-test" --next-file ".d2a/tests/01_integration_tasks.md" --summary "Mini-test work complete; moving into confirmation questions."`

## Phase 2: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - first end-to-end scenario
   - observable success signal
   - observable failure signal
   - why these tests are enough to validate the mini architecture idea
3. Ask one question per turn.
4. Before each question, print a compact header:
   - `d2a repo: ...`
   - `d2a repo path: ...`
   - `d2a path: ...`
   - `d2a stage: testing-in-progress`
   - `d2a phase: confirmation-questions`
   - `d2a question progress: N/4`
   - `d2a next step after questions: continue to d2a-report-build`
5. Before asking question `N`, call:

   `d2a skill-state d2a-mini-test --status progress --stage testing-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue mini-test confirmation questions." --next-skill "d2a-report-build" --next-file ".d2a/report/index.md" --summary "Mini-test confirmation question <N> is active."`

6. Present one question with multiple choices.
7. Wait for the learner answer.
8. After the learner answer:
   - say whether the answer is correct, partially correct, or incorrect
   - give one short explanation
   - continue to the next question even if the answer is wrong
9. After question 4 is evaluated:
   - output a short recap
   - output a `理解度打分`
   - keep the `理解度打分` under 100 Chinese characters
10. At the end of the confirmation phase, call:

    `d2a skill-state d2a-mini-test --status completed --stage testing-complete --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-report-build." --next-skill "d2a-report-build" --next-file ".d2a/report/index.md" --summary "Completed mini-test confirmation questions."`

## Output

- First integration test
- Next integration scenario
- Observable success signal
- Observable failure signal
- Four confirmation questions
- Short comprehension summary
- `理解度打分`
