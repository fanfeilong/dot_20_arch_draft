---
name: d2a-runtime-view
description: Built-in d2a skill for d2a-runtime-view stage guidance and state updates.
---

# d2a-runtime-view

## Goal

Describe what drives the system and which module acts as the engine, then verify that the learner actually understands the result.

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

   `d2a skill-state d2a-runtime-view --status started --stage architecture-in-progress --phase analysis-generation --next-step "Identify the driver, core loop, and engine module." --next-skill "d2a-core-objects" --next-file ".d2a/docs/architecture/02_driver.md" --summary "Started runtime-view analysis."`

2. Work from the real repository and write the result into `.d2a/docs/architecture/02_driver.md`.
3. Answer these atomic questions:
   - What is the dominant runtime driver?
   - What is the core runtime loop?
   - Which module is the engine?
   - Which up to two supporting modules are required to understand the engine?
4. After the first draft, force three correction passes:
   - compression pass
   - de-jargon pass
   - conversational simplification pass
5. Prefer a compact structural view over long explanation.
6. When the analysis draft is stable, call:

   `d2a skill-state d2a-runtime-view --status progress --stage architecture-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first runtime-view confirmation question." --next-skill "d2a-runtime-view" --next-file ".d2a/docs/architecture/02_driver.md" --summary "Runtime-view analysis complete; moving into confirmation questions."`

## Phase 2: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - dominant driver
   - core runtime loop
   - engine module
   - role of supporting modules
3. Ask one question per turn.
4. Before each question, print a compact header:
   - `d2a repo: ...`
   - `d2a repo path: ...`
   - `d2a path: ...`
   - `d2a stage: architecture-in-progress`
   - `d2a phase: confirmation-questions`
   - `d2a question progress: N/4`
   - `d2a next step after questions: continue to d2a-core-objects`
5. Before asking question `N`, call:

   `d2a skill-state d2a-runtime-view --status progress --stage architecture-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue runtime-view confirmation questions." --next-skill "d2a-core-objects" --next-file ".d2a/docs/architecture/03_core_objects.md" --summary "Runtime-view confirmation question <N> is active."`

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

    `d2a skill-state d2a-runtime-view --status completed --stage architecture-in-progress --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-core-objects." --next-skill "d2a-core-objects" --next-file ".d2a/docs/architecture/03_core_objects.md" --summary "Completed runtime-view confirmation questions."`

## Turn-End Continuation Rule

1. At the end of every learner-facing turn, call `d2a skill-state` to persist the latest phase and progress before finishing the reply.
2. If the current phase is still in progress, record `--status progress` with updated question or challenge index.
3. End the reply with: `继续请使用 $d2a-step`.
4. If the current phase just completed, still emit the same continuation line so the learner can resume from the next state via `d2a-step`.

## Output

- Primary driver
- Core loop
- Engine module
- Supporting modules
- Four confirmation questions
- Short comprehension summary
- `理解度打分`
