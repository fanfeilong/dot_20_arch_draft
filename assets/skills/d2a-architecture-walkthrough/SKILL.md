---
name: d2a-architecture-walkthrough
description: Built-in d2a skill for d2a-architecture-walkthrough stage guidance and state updates.
---

# d2a-architecture-walkthrough

## Goal

Use the six architecture skills to compress a project into a minimal architecture skeleton, then verify that the learner understands the full architecture set as a whole.

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

   `d2a skill-state d2a-architecture-walkthrough --status started --stage architecture-in-progress --phase analysis-generation --next-step "Check the six architecture files and produce the overview summary." --next-skill "d2a-project-scope" --next-file ".d2a/docs/architecture/00_overview.md" --summary "Started architecture walkthrough."`

2. Treat this skill as the user-facing architecture-stage entry inside Codex.
3. If architecture task files have not yet been prepared for the current target, call `d2a analyze <target-repo>` before deeper analysis.
4. Work in this order:
   - `d2a-project-scope`
   - `d2a-runtime-view`
   - `d2a-core-objects`
   - `d2a-state-evolution`
   - `d2a-module-view`
   - `d2a-tradeoff-view`
5. Treat the six architecture elements as mandatory coverage:
   - boundary
   - driver
   - core objects
   - state evolution
   - cooperation
   - constraints
6. Ensure the following files are completed or updated as needed:
   - `.d2a/docs/architecture/01_boundary.md`
   - `.d2a/docs/architecture/02_driver.md`
   - `.d2a/docs/architecture/03_core_objects.md`
   - `.d2a/docs/architecture/04_state_evolution.md`
   - `.d2a/docs/architecture/05_cooperation.md`
   - `.d2a/docs/architecture/06_constraints.md`
7. Then produce the cross-cutting summary files:
   - `.d2a/docs/architecture/00_overview.md`
   - `.d2a/docs/architecture/99_code_map.md`
8. In `00_overview.md`, ensure the learner can answer:
   - what the system is
   - where its boundary is
   - what drives it
   - what its core objects are
   - how its key state evolves
   - how its main modules cooperate
   - what constraint most shapes it
9. In `99_code_map.md`, map each major claim back to directories or files.
10. After the first draft, force three correction passes:
   - compression pass
   - de-jargon pass
   - conversational simplification pass
11. When the architecture set is stable, call:

   `d2a skill-state d2a-architecture-walkthrough --status progress --stage architecture-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first architecture walkthrough confirmation question." --next-skill "d2a-architecture-walkthrough" --next-file ".d2a/docs/architecture/00_overview.md" --summary "Architecture set complete; moving into walkthrough confirmation questions."`

## Phase 2: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual architecture set, not from generic examples.
2. The 4 questions should cover these angles:
   - one-sentence system definition and boundary
   - dominant driver and engine
   - core objects and state evolution
   - cooperation plus main constraint / tradeoff
3. Ask one question per turn.
4. Before each question, print a compact header:
   - `d2a repo: ...`
   - `d2a repo path: ...`
   - `d2a path: ...`
   - `d2a stage: architecture-in-progress`
   - `d2a phase: confirmation-questions`
   - `d2a question progress: N/4`
   - `d2a next step after questions: prepare the architecture challenge phase`
5. Before asking question `N`, call:

   `d2a skill-state d2a-architecture-walkthrough --status progress --stage architecture-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue architecture walkthrough confirmation questions." --next-skill "d2a-challenge-architecture" --next-file ".d2a/docs/architecture/00_overview.md" --summary "Architecture walkthrough confirmation question <N> is active."`

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

    `d2a skill-state d2a-architecture-walkthrough --status completed --stage architecture-complete --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Enter the architecture challenge phase before mini derivation." --next-skill "d2a-challenge-architecture" --next-file ".d2a/docs/architecture/00_overview.md" --summary "Completed architecture walkthrough confirmation questions."`

## Turn-End Continuation Rule

1. At the end of every learner-facing turn, call `d2a skill-state` to persist the latest phase and progress before finishing the reply.
2. If the current phase is still in progress, record `--status progress` with updated question or challenge index.
3. End the reply with: `继续请使用 $d2a-step`.
4. If the current phase just completed, still emit the same continuation line so the learner can resume from the next state via `d2a-step`.

## Output

- Minimal architecture skeleton
- `00_overview.md`
- `99_code_map.md`
- Four confirmation questions
- Short comprehension summary
- `理解度打分`
