---
name: d2a-arch-2-runtime-view
description: Built-in d2a skill for d2a-arch-2-runtime-view stage guidance and state updates.
---

# d2a-arch-2-runtime-view

## Goal

Describe what drives the system and which module acts as the engine, then verify that the learner actually understands the result.

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


## Body Formatting Hard Rules

1. The body must use `- ` bullet lines, minimum 2 lines and maximum 4 lines.
2. Each point must occupy exactly one line; do not write long paragraphs.
3. Keep each line under about 140 characters; split lines when longer.
4. Avoid Markdown emphasis in body text (for example `` `...` `` or `**...**`).

If the active repository is unknown, stop and ask the user which repository should be used.

## Human In Loop Marker Rule

When the current turn asks the user a question and waits for user input, the last line of the reply body must append:

`[human_in_loop]`

## Phase 1: Atomic Question Alignment (One Supplement Chance)

1. After context is confirmed, call:

   `d2a skill-state d2a-arch-2-runtime-view --status started --stage architecture-in-progress --phase atomic-question-alignment --question-index 0 --question-total 1 --next-step "Show atomic analysis questions and ask for one-time supplements (yes/no)." --next-skill "d2a-arch-2-runtime-view" --next-file ".d2a/docs/architecture/02_driver.md" --summary "Started atomic-question alignment for runtime-view."`

2. Before doing analysis, explicitly tell the user:

   `Next I will scan and answer these questions: <list this skill's atomic questions>; do you want to add any? (yes/no)`

3. Allow only one supplement interaction in this phase:
   - If user answers `yes`, collect extra questions, merge them with the base atomic questions, then echo the merged list.
   - If user answers `no`, proceed with the base atomic questions.
4. Do not write `.d2a/docs/architecture/02_driver.md` before user confirmation.
5. After alignment is done, call:

   `d2a skill-state d2a-arch-2-runtime-view --status progress --stage architecture-in-progress --phase analysis-generation --question-index 1 --question-total 1 --next-step "Start analysis-generation with merged atomic questions." --next-skill "d2a-arch-2-runtime-view" --next-file ".d2a/docs/architecture/02_driver.md" --summary "Atomic-question alignment completed for runtime-view."`

## Phase 2: Analysis Generation

1. After context is confirmed, call:

   `d2a skill-state d2a-arch-2-runtime-view --status progress --stage architecture-in-progress --phase analysis-generation --next-step "Identify the driver, core loop, and engine module." --next-skill "d2a-arch-3-core-objects" --next-file ".d2a/docs/architecture/02_driver.md" --summary "Started runtime-view analysis."`

2. Work from the real repository and write the result into `.d2a/docs/architecture/02_driver.md`.
3. Answer these merged atomic questions (base + optional user supplements):
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

   `d2a skill-state d2a-arch-2-runtime-view --status progress --stage architecture-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first runtime-view confirmation question." --next-skill "d2a-arch-2-runtime-view" --next-file ".d2a/docs/architecture/02_driver.md" --summary "Runtime-view analysis complete; moving into confirmation questions."`

## Phase 3: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - dominant driver
   - core runtime loop
   - engine module
   - role of supporting modules
3. Ask one question per turn.
4. Before each question, keep the same envelope format and map fields as follows:
   - The `[d2a]` line must include current stage, phase, and `N/<total>` progress.
   - The `[next]` line should point to the post-question next skill/file.
5. Before asking question `N`, call:

   `d2a skill-state d2a-arch-2-runtime-view --status progress --stage architecture-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue runtime-view confirmation questions." --next-skill "d2a-arch-3-core-objects" --next-file ".d2a/docs/architecture/03_core_objects.md" --summary "Runtime-view confirmation question <N> is active."`

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

    `d2a skill-state d2a-arch-2-runtime-view --status completed --stage architecture-in-progress --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-arch-3-core-objects." --next-skill "d2a-arch-3-core-objects" --next-file ".d2a/docs/architecture/03_core_objects.md" --summary "Completed runtime-view confirmation questions."`

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
