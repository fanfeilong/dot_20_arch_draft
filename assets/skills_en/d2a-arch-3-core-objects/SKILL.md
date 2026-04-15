---
name: d2a-arch-3-core-objects
description: Built-in d2a skill for d2a-arch-3-core-objects stage guidance and state updates.
---

# d2a-arch-3-core-objects

## Goal

Extract the smallest useful set of core objects around which the system is organized, then verify that the learner actually understands the result.

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

## Human In Loop Marker Rule

When the current turn asks the user a question and waits for user input, the last line of the reply body must append:

`[human_in_loop]`

## Phase 1: Atomic Question Alignment (One Supplement Chance)

1. After context is confirmed, call:

   `d2a skill-state d2a-arch-3-core-objects --status started --stage architecture-in-progress --phase atomic-question-alignment --question-index 0 --question-total 1 --next-step "Show atomic analysis questions and ask for one-time supplements (yes/no)." --next-skill "d2a-arch-3-core-objects" --next-file ".d2a/docs/architecture/03_core_objects.md" --summary "Started atomic-question alignment for core-objects."`

2. Before doing analysis, explicitly tell the user:

   `Next I will scan and answer these questions: <list this skill's atomic questions>; do you want to add any? (yes/no)`

3. Allow only one supplement interaction in this phase:
   - If user answers `yes`, collect extra questions, merge them with the base atomic questions, then echo the merged list.
   - If user answers `no`, proceed with the base atomic questions.
4. Do not write `.d2a/docs/architecture/03_core_objects.md` before user confirmation.
5. After alignment is done, call:

   `d2a skill-state d2a-arch-3-core-objects --status progress --stage architecture-in-progress --phase analysis-generation --question-index 1 --question-total 1 --next-step "Start analysis-generation with merged atomic questions." --next-skill "d2a-arch-3-core-objects" --next-file ".d2a/docs/architecture/03_core_objects.md" --summary "Atomic-question alignment completed for core-objects."`

## Phase 2: Analysis Generation

1. After context is confirmed, call:

   `d2a skill-state d2a-arch-3-core-objects --status progress --stage architecture-in-progress --phase analysis-generation --next-step "Identify the core objects, their relations, and the state center." --next-skill "d2a-arch-4-state-evolution" --next-file ".d2a/docs/architecture/03_core_objects.md" --summary "Started core-objects analysis."`

2. Work from the real repository and write the result into `.d2a/docs/architecture/03_core_objects.md`.
3. Answer these merged atomic questions (base + optional user supplements):
   - What are the at most three core objects?
   - Who creates, consumes, persists, or drives them?
   - Where is the state center?
4. After the first draft, force three correction passes:
   - compression pass
   - de-jargon pass
   - conversational simplification pass
5. Avoid listing incidental structs, DTOs, or file formats unless they are architecture-critical.
6. When the analysis draft is stable, call:

   `d2a skill-state d2a-arch-3-core-objects --status progress --stage architecture-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first core-objects confirmation question." --next-skill "d2a-arch-3-core-objects" --next-file ".d2a/docs/architecture/03_core_objects.md" --summary "Core-objects analysis complete; moving into confirmation questions."`

## Phase 3: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - identity of the core objects
   - object relations
   - state center
   - why some large-looking types are not architecture core
3. Ask one question per turn.
4. Before each question, keep the same envelope format and map fields as follows:
   - The `[d2a]` line must include current stage, phase, and `N/<total>` progress.
   - The `[next]` line should point to the post-question next skill/file.
5. Before asking question `N`, call:

   `d2a skill-state d2a-arch-3-core-objects --status progress --stage architecture-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue core-objects confirmation questions." --next-skill "d2a-arch-4-state-evolution" --next-file ".d2a/docs/architecture/04_state_evolution.md" --summary "Core-objects confirmation question <N> is active."`

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

    `d2a skill-state d2a-arch-3-core-objects --status completed --stage architecture-in-progress --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-arch-4-state-evolution." --next-skill "d2a-arch-4-state-evolution" --next-file ".d2a/docs/architecture/04_state_evolution.md" --summary "Completed core-objects confirmation questions."`

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
