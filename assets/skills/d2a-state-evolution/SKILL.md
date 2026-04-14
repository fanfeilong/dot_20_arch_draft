---
name: d2a-state-evolution
description: Built-in d2a skill for d2a-state-evolution stage guidance and state updates.
---

# d2a-state-evolution

## Goal

Describe how the most important object or workflow changes state over time, then verify that the learner actually understands the result.

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

   `d2a skill-state d2a-state-evolution --status started --stage architecture-in-progress --phase analysis-generation --next-step "Track the main state stages and transitions." --next-skill "d2a-module-view" --next-file ".d2a/docs/architecture/04_state_evolution.md" --summary "Started state-evolution analysis."`

2. Work from the real repository and write the result into `.d2a/docs/architecture/04_state_evolution.md`.
3. Answer these atomic questions:
   - What single object or workflow should be tracked?
   - What are its three to six state stages?
   - What triggers the main state transitions?
   - Where is state stored, observed, or reconstructed?
4. After the first draft, force three correction passes:
   - compression pass
   - de-jargon pass
   - conversational simplification pass
5. Prefer state progression over implementation detail.
6. When the analysis draft is stable, call:

   `d2a skill-state d2a-state-evolution --status progress --stage architecture-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first state-evolution confirmation question." --next-skill "d2a-state-evolution" --next-file ".d2a/docs/architecture/04_state_evolution.md" --summary "State-evolution analysis complete; moving into confirmation questions."`

## Phase 2: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - tracked object or workflow
   - state stages
   - transition triggers
   - state storage / observation point
3. Ask one question per turn.
4. Before each question, print a compact header:
   - `d2a repo: ...`
   - `d2a repo path: ...`
   - `d2a path: ...`
   - `d2a stage: architecture-in-progress`
   - `d2a phase: confirmation-questions`
   - `d2a question progress: N/4`
   - `d2a next step after questions: continue to d2a-module-view`
5. Before asking question `N`, call:

   `d2a skill-state d2a-state-evolution --status progress --stage architecture-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue state-evolution confirmation questions." --next-skill "d2a-module-view" --next-file ".d2a/docs/architecture/05_cooperation.md" --summary "State-evolution confirmation question <N> is active."`

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

    `d2a skill-state d2a-state-evolution --status completed --stage architecture-in-progress --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-module-view." --next-skill "d2a-module-view" --next-file ".d2a/docs/architecture/05_cooperation.md" --summary "Completed state-evolution confirmation questions."`

## Output

- Tracked object or workflow
- State stages
- Transition triggers
- State storage or observation point
- Four confirmation questions
- Short comprehension summary
- `理解度打分`
