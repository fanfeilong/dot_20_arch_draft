---
name: d2a-project-scope
description: Built-in d2a skill for d2a-project-scope stage guidance and state updates.
---

# d2a-project-scope

## Goal

Define the system boundary before deeper architecture analysis, then verify that the learner actually understands the result.

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

   `d2a skill-state d2a-project-scope --status started --stage architecture-in-progress --phase analysis-generation --next-step "Define the system boundary and entry points." --next-skill "d2a-runtime-view" --next-file ".d2a/docs/architecture/01_boundary.md" --summary "Started project-scope analysis."`

2. Work from the real repository and write the result into `.d2a/docs/architecture/01_boundary.md`.
3. Answer these atomic questions:
   - What kind of system is this?
   - What is the one-sentence definition?
   - What capability must remain if 80% of the code were deleted?
   - What are the one to three best entry points?
   - What is inside the system boundary?
   - What is outside the system boundary?
4. After the first draft, force three correction passes:
   - compression pass
   - de-jargon pass
   - conversational simplification pass
5. Keep the result short, structural, and free of project history.
6. When the analysis draft is stable, call:

   `d2a skill-state d2a-project-scope --status progress --stage architecture-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first project-scope confirmation question." --next-skill "d2a-project-scope" --next-file ".d2a/docs/architecture/01_boundary.md" --summary "Project-scope analysis complete; moving into confirmation questions."`

## Phase 2: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - system type
   - one-sentence definition
   - non-removable capability
   - boundary / in-scope vs out-of-scope
3. Ask one question per turn.
4. Before each question, print a compact header:
   - `d2a repo: ...`
   - `d2a repo path: ...`
   - `d2a path: ...`
   - `d2a stage: architecture-in-progress`
   - `d2a phase: confirmation-questions`
   - `d2a question progress: N/4`
   - `d2a next step after questions: continue to d2a-runtime-view`
5. Before asking question `N`, call:

   `d2a skill-state d2a-project-scope --status progress --stage architecture-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue project-scope confirmation questions." --next-skill "d2a-runtime-view" --next-file ".d2a/docs/architecture/02_driver.md" --summary "Project-scope confirmation question <N> is active."`

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

    `d2a skill-state d2a-project-scope --status completed --stage architecture-in-progress --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-runtime-view." --next-skill "d2a-runtime-view" --next-file ".d2a/docs/architecture/02_driver.md" --summary "Completed project-scope confirmation questions."`

## Output

- System type
- One-sentence definition
- Non-removable core capability
- Entry points
- Outer boundary
- In scope
- Out of scope
- Four confirmation questions
- Short comprehension summary
- `理解度打分`
