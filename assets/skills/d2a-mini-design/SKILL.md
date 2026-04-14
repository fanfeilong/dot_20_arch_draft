# d2a-mini-design

## Goal

Turn the chosen mini scope into a concrete design for the implementation under `.d2a/src/`, then verify that the learner actually understands why this design is the right minimal form.

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

   `d2a skill-state d2a-mini-design --status started --stage mini-design-in-progress --phase analysis-generation --next-step "Design the smallest useful runnable mini architecture." --next-skill "d2a-mini-build" --next-file ".d2a/docs/implementation/01_mini_design.md" --summary "Started mini-design work."`

2. If implementation planning files have not yet been prepared, call `d2a derive-mini`.
3. Read:
   - `.d2a/docs/implementation/00_mini_scope.md`
   - the architecture files it references
4. Write the result into `.d2a/docs/implementation/01_mini_design.md`.
5. Keep `.d2a/src/ARCHITECTURE.md` aligned with the chosen design if the summary there is stale.
6. Answer these atomic questions:
   - What are the main modules of the mini version?
   - What interfaces or entry points are required?
   - What is the runtime flow of the mini version?
   - What state model must be preserved?
7. After the first draft, force three correction passes:
   - compression pass
   - de-jargon pass
   - conversational simplification pass
8. Keep the design small enough to support one first runnable slice.
9. Keep the design tied to the core architecture idea, not to broad feature coverage.
10. When the analysis draft is stable, call:

   `d2a skill-state d2a-mini-design --status progress --stage mini-design-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first mini-design confirmation question." --next-skill "d2a-mini-design" --next-file ".d2a/docs/implementation/01_mini_design.md" --summary "Mini-design analysis complete; moving into confirmation questions."`

## Phase 2: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - main modules
   - interfaces or entry points
   - runtime flow
   - preserved state model
3. Ask one question per turn.
4. Before each question, print a compact header:
   - `d2a repo: ...`
   - `d2a repo path: ...`
   - `d2a path: ...`
   - `d2a stage: mini-design-in-progress`
   - `d2a phase: confirmation-questions`
   - `d2a question progress: N/4`
   - `d2a next step after questions: continue to d2a-mini-build`
5. Before asking question `N`, call:

   `d2a skill-state d2a-mini-design --status progress --stage mini-design-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue mini-design confirmation questions." --next-skill "d2a-mini-build" --next-file ".d2a/src/ARCHITECTURE.md" --summary "Mini-design confirmation question <N> is active."`

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

    `d2a skill-state d2a-mini-design --status completed --stage mini-design-complete --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-mini-build." --next-skill "d2a-mini-build" --next-file ".d2a/src/ARCHITECTURE.md" --summary "Completed mini-design confirmation questions."`

## Output

- Main modules
- Main interfaces or entry points
- Runtime flow
- State model
- Four confirmation questions
- Short comprehension summary
- `理解度打分`
