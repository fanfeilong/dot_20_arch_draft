# d2a-core-objects

## Goal

Extract the smallest useful set of core objects around which the system is organized, then verify that the learner actually understands the result.

## Required Start Header

Always print this header first:

- `d2a repo: ...`
- `d2a repo path: ...`
- `d2a path: ...`
- `d2a stage: ...`
- `d2a flow: initialized -> analysis-prepared -> architecture-in-progress -> architecture-complete -> architecture-challenge-prepared -> architecture-challenge-complete -> mini-derivation-prepared -> mini-design-complete -> test-plan-prepared -> testing-in-progress -> testing-complete -> report-ready`
- `d2a phase: ...`
- `d2a next step: ...`

If the active repository is unknown, stop and ask the user which repository should be used.

## Phase 1: Analysis Generation

1. After context is confirmed, call:

   `d2a skill-state d2a-core-objects --status started --stage architecture-in-progress --phase analysis-generation --next-step "Identify the core objects, their relations, and the state center." --next-skill "d2a-state-evolution" --next-file ".d2a/docs/architecture/03_core_objects.md" --summary "Started core-objects analysis."`

2. Work from the real repository and write the result into `.d2a/docs/architecture/03_core_objects.md`.
3. Answer these atomic questions:
   - What are the at most three core objects?
   - Who creates, consumes, persists, or drives them?
   - Where is the state center?
4. After the first draft, force three correction passes:
   - compression pass
   - de-jargon pass
   - conversational simplification pass
5. Avoid listing incidental structs, DTOs, or file formats unless they are architecture-critical.
6. When the analysis draft is stable, call:

   `d2a skill-state d2a-core-objects --status progress --stage architecture-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first core-objects confirmation question." --next-skill "d2a-core-objects" --next-file ".d2a/docs/architecture/03_core_objects.md" --summary "Core-objects analysis complete; moving into confirmation questions."`

## Phase 2: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - identity of the core objects
   - object relations
   - state center
   - why some large-looking types are not architecture core
3. Ask one question per turn.
4. Before each question, print a compact header:
   - `d2a repo: ...`
   - `d2a repo path: ...`
   - `d2a path: ...`
   - `d2a stage: architecture-in-progress`
   - `d2a phase: confirmation-questions`
   - `d2a question progress: N/4`
   - `d2a next step after questions: continue to d2a-state-evolution`
5. Before asking question `N`, call:

   `d2a skill-state d2a-core-objects --status progress --stage architecture-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue core-objects confirmation questions." --next-skill "d2a-state-evolution" --next-file ".d2a/docs/architecture/04_state_evolution.md" --summary "Core-objects confirmation question <N> is active."`

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

    `d2a skill-state d2a-core-objects --status completed --stage architecture-in-progress --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-state-evolution." --next-skill "d2a-state-evolution" --next-file ".d2a/docs/architecture/04_state_evolution.md" --summary "Completed core-objects confirmation questions."`

## Output

- Core objects
- Object relations
- State center
- Four confirmation questions
- Short comprehension summary
- `理解度打分`
