---
name: d2a-mini-build
description: Built-in d2a skill for d2a-mini-build stage guidance and state updates.
---

# d2a-mini-build

## Goal

Implement the first runnable mini version in `.d2a/src/` while preserving the target project's core architecture idea, then verify that the learner actually understands why this implementation slice is sufficient.

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

   `d2a skill-state d2a-mini-build --status started --stage mini-design-complete --phase analysis-generation --next-step "Implement the first runnable mini slice." --next-skill "d2a-mini-test" --next-file ".d2a/src/ARCHITECTURE.md" --summary "Started mini-build work."`

2. If implementation planning files have not yet been prepared, call `d2a derive-mini`.
3. Read:
   - `.d2a/docs/implementation/00_mini_scope.md`
   - `.d2a/docs/implementation/01_mini_design.md`
   - `.d2a/docs/implementation/02_build_plan.md`
   - `.d2a/src/ARCHITECTURE.md`
4. Implement only the first runnable slice described in the build plan.
5. Prefer a small but executable result over broad coverage.
6. Keep the chosen stack aligned with the original project when practical.
7. Update `.d2a/src/ARCHITECTURE.md` if implementation reality forces a design correction.
8. Do not expand scope until the first runnable slice is working.
9. After the implementation is stable, write a brief note on what remains intentionally unimplemented.
10. When the implementation pass is stable, call:

   `d2a skill-state d2a-mini-build --status progress --stage mini-design-complete --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first mini-build confirmation question." --next-skill "d2a-mini-build" --next-file ".d2a/src/ARCHITECTURE.md" --summary "Mini-build implementation complete; moving into confirmation questions."`

## Phase 2: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - what runnable slice was actually implemented
   - why this slice is enough to demonstrate the core architecture idea
   - what was intentionally left unimplemented
   - how `src/ARCHITECTURE.md` matches the code reality
3. Ask one question per turn.
4. Before each question, print a compact header:
   - `d2a repo: ...`
   - `d2a repo path: ...`
   - `d2a path: ...`
   - `d2a stage: mini-design-complete`
   - `d2a phase: confirmation-questions`
   - `d2a question progress: N/4`
   - `d2a next step after questions: continue to d2a-mini-test`
5. Before asking question `N`, call:

   `d2a skill-state d2a-mini-build --status progress --stage mini-design-complete --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue mini-build confirmation questions." --next-skill "d2a-mini-test" --next-file ".d2a/tests/README.md" --summary "Mini-build confirmation question <N> is active."`

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

    `d2a skill-state d2a-mini-build --status completed --stage mini-design-complete --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-mini-test." --next-skill "d2a-mini-test" --next-file ".d2a/tests/README.md" --summary "Completed mini-build confirmation questions."`

## Output

- Runnable mini implementation under `.d2a/src/`
- Updated `.d2a/src/ARCHITECTURE.md`
- Brief note on what remains intentionally unimplemented
- Four confirmation questions
- Short comprehension summary
- `理解度打分`
