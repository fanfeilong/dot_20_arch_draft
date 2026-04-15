---
name: d2a-arch-6-tradeoff-view
description: Built-in d2a skill for d2a-arch-6-tradeoff-view stage guidance and state updates.
---

# d2a-arch-6-tradeoff-view

## Goal

State the strongest architectural constraint and the main tradeoff it forces, then verify that the learner actually understands the result.

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

## Phase 1: Atomic Question Alignment (One Supplement Chance)

1. After context is confirmed, call:

   `d2a skill-state d2a-arch-6-tradeoff-view --status started --stage architecture-in-progress --phase atomic-question-alignment --question-index 0 --question-total 1 --next-step "Show atomic analysis questions and ask for one-time supplements (yes/no)." --next-skill "d2a-arch-6-tradeoff-view" --next-file ".d2a/docs/architecture/06_constraints.md" --summary "Started atomic-question alignment for tradeoff-view."`

2. Before doing analysis, explicitly tell the user:

   `Next I will scan and answer these questions: <list this skill's atomic questions>; do you want to add any? (yes/no)`

3. Allow only one supplement interaction in this phase:
   - If user answers `yes`, collect extra questions, merge them with the base atomic questions, then echo the merged list.
   - If user answers `no`, proceed with the base atomic questions.
4. Do not write `.d2a/docs/architecture/06_constraints.md` before user confirmation.
5. After alignment is done, call:

   `d2a skill-state d2a-arch-6-tradeoff-view --status progress --stage architecture-in-progress --phase analysis-generation --question-index 1 --question-total 1 --next-step "Start analysis-generation with merged atomic questions." --next-skill "d2a-arch-6-tradeoff-view" --next-file ".d2a/docs/architecture/06_constraints.md" --summary "Atomic-question alignment completed for tradeoff-view."`

## Phase 2: Analysis Generation

1. After context is confirmed, call:

   `d2a skill-state d2a-arch-6-tradeoff-view --status progress --stage architecture-in-progress --phase analysis-generation --next-step "Identify the strongest constraints and architectural tradeoff." --next-skill "d2a-challenge-architecture" --next-file ".d2a/docs/architecture/06_constraints.md" --summary "Started tradeoff-view analysis."`

2. Work from the real repository and write the result into `.d2a/docs/architecture/06_constraints.md`.
3. Answer these merged atomic questions (base + optional user supplements):
   - What are the two to four hard constraints?
   - Which one is dominant?
   - What tradeoff does it force?
   - Which structures must remain if the system is rewritten?
   - Which large-looking parts are implementation detail rather than architecture core?
4. After the first draft, force three correction passes:
   - compression pass
   - de-jargon pass
   - conversational simplification pass
5. Prefer architectural pressure and design consequence over opinion.
6. When the analysis draft is stable, call:

   `d2a skill-state d2a-arch-6-tradeoff-view --status progress --stage architecture-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first tradeoff-view confirmation question." --next-skill "d2a-arch-6-tradeoff-view" --next-file ".d2a/docs/architecture/06_constraints.md" --summary "Tradeoff-view analysis complete; moving into confirmation questions."`

## Phase 3: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - hard constraints
   - dominant constraint
   - main tradeoff
   - must-keep structures vs implementation detail
3. Ask one question per turn.
4. Before each question, keep the same envelope format and map fields as follows:
   - The `[d2a]` line must include current stage, phase, and `N/<total>` progress.
   - The `[next]` line should point to the post-question next skill/file.
5. Before asking question `N`, call:

   `d2a skill-state d2a-arch-6-tradeoff-view --status progress --stage architecture-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue tradeoff-view confirmation questions." --next-skill "d2a-challenge-architecture" --next-file ".d2a/docs/architecture/06_constraints.md" --summary "Tradeoff-view confirmation question <N> is active."`

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

    `d2a skill-state d2a-arch-6-tradeoff-view --status completed --stage architecture-complete --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-challenge-architecture." --next-skill "d2a-challenge-architecture" --next-file ".d2a/docs/architecture/06_constraints.md" --summary "Completed tradeoff-view confirmation questions."`

11. Confirmation-question prompts, learner answers, evaluations, and explanations must be written to `.d2a/qa/<skill>.jsonl`, and must not be written into `docs/architecture/*.md`.

## Turn-End Continuation Rule

1. Before ending each turn, call `d2a skill-state` to persist the current phase and progress.
2. End the reply with: `Continue with $d2a-step`.


## Output (Artifacts)

- Stage architecture conclusions written into `docs/architecture/*`.

## Persistence (.d2a)

- Aligned atomic question set (base + user supplements) -> `.d2a/qa/<skill>.json`.
- Confirmation prompts, learner answers, evaluations, explanations -> `.d2a/qa/<skill>.jsonl`.
- Short comprehension summary and `Comprehension Score` -> `.d2a/qa/<skill>.jsonl`.
