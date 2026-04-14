---
name: d2a-mini-scope
description: Built-in d2a skill for d2a-mini-scope stage guidance and state updates.
---

# d2a-mini-scope

## Goal

Choose the smallest runnable slice that preserves the target project's core architecture idea, then verify that the learner actually understands the scope choice.

## Required Start Header

Always print this header first:

- `d2a repo: ...`
- `d2a repo path: ...`
- `d2a path: ...`
- `d2a stage: ...`
- `d2a flow: initialized -> analysis-prepared -> architecture-in-progress -> architecture-complete -> architecture-challenge-prepared -> architecture-challenge-in-progress -> architecture-challenge-complete -> mini-derivation-prepared -> mini-design-complete -> test-plan-prepared -> testing-in-progress -> testing-complete -> report-ready`
- `d2a phase: ...`
- `d2a next step: ...`

If the active repository is unknown, stop and ask the user which repository should be used.

## Phase 1: Analysis Generation

1. After context is confirmed, call:

   `d2a skill-state d2a-mini-scope --status started --stage mini-derivation-prepared --phase analysis-generation --next-step "Choose the single architecture idea to preserve." --next-skill "d2a-mini-design" --next-file ".d2a/docs/implementation/00_mini_scope.md" --summary "Started mini-scope derivation."`

2. Treat this skill as the user-facing entry for the mini-implementation stage inside Codex.
3. If implementation planning files have not yet been prepared, call `d2a derive-mini` before producing content.
4. Read:
   - `.d2a/docs/architecture/00_overview.md`
   - `.d2a/docs/architecture/02_driver.md`
   - `.d2a/docs/architecture/03_core_objects.md`
   - `.d2a/docs/architecture/04_state_evolution.md`
   - `.d2a/docs/architecture/05_cooperation.md`
   - `.d2a/docs/architecture/06_constraints.md`
5. Write the result into `.d2a/docs/implementation/00_mini_scope.md`.
6. Answer these atomic questions:
   - What single architecture idea must be preserved?
   - Which runnable 20 percent slice is enough to demonstrate it?
   - What will the mini version intentionally omit?
   - Which stack should stay aligned with the original project?
7. After the first draft, force three correction passes:
   - compression pass
   - de-jargon pass
   - conversational simplification pass
8. Keep the scope small enough to support one first runnable slice.
9. When the analysis draft is stable, call:

   `d2a skill-state d2a-mini-scope --status progress --stage mini-derivation-prepared --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first mini-scope confirmation question." --next-skill "d2a-mini-scope" --next-file ".d2a/docs/implementation/00_mini_scope.md" --summary "Mini-scope analysis complete; moving into confirmation questions."`

## Phase 2: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - preserved architecture idea
   - runnable 20 percent slice
   - intentional omissions
   - chosen stack
3. Ask one question per turn.
4. Before each question, print a compact header:
   - `d2a repo: ...`
   - `d2a repo path: ...`
   - `d2a path: ...`
   - `d2a stage: mini-derivation-prepared`
   - `d2a phase: confirmation-questions`
   - `d2a question progress: N/4`
   - `d2a next step after questions: continue to d2a-mini-design`
5. Before asking question `N`, call:

   `d2a skill-state d2a-mini-scope --status progress --stage mini-derivation-prepared --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue mini-scope confirmation questions." --next-skill "d2a-mini-design" --next-file ".d2a/docs/implementation/01_mini_design.md" --summary "Mini-scope confirmation question <N> is active."`

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

    `d2a skill-state d2a-mini-scope --status completed --stage mini-derivation-prepared --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-mini-design." --next-skill "d2a-mini-design" --next-file ".d2a/docs/implementation/01_mini_design.md" --summary "Completed mini-scope confirmation questions."`

## Turn-End Continuation Rule

1. At the end of every learner-facing turn, call `d2a skill-state` to persist the latest phase and progress before finishing the reply.
2. If the current phase is still in progress, record `--status progress` with updated question or challenge index.
3. End the reply with: `继续请使用 $d2a-step`.
4. If the current phase just completed, still emit the same continuation line so the learner can resume from the next state via `d2a-step`.

## Output

- Preserved architecture idea
- Runnable 20 percent slice
- Intentional omissions
- Target stack
- Four confirmation questions
- Short comprehension summary
- `理解度打分`
