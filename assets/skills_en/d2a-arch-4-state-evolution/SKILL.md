---
name: d2a-arch-4-state-evolution
description: Built-in d2a skill for d2a-arch-4-state-evolution stage guidance and state updates.
---

# d2a-arch-4-state-evolution

## Goal

Describe how the most important object or workflow changes state over time, then verify that the learner actually understands the result.

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
3. Keep each line under 80 characters; split lines when longer.
4. Avoid Markdown emphasis in body text (for example `` `...` `` or `**...**`).

5. For structured content (lists, numbered items, MCQ options A/B/C/D), do not merge lines just to satisfy width limits; keep one item per line.

If the active repository is unknown, stop and ask the user which repository should be used.

## Human In Loop Marker Rule

When the current turn asks the user a question and waits for user input, the last line of the reply body must append:

`[human_in_loop]`

## Phase 1: Atomic Question Alignment (One Supplement Chance)

1. After context is confirmed, call:

   `d2a skill-state d2a-arch-4-state-evolution --status started --stage architecture-in-progress --phase atomic-question-alignment --question-index 0 --question-total 1 --next-step "Show atomic analysis questions and ask for one-time supplements (yes/no)." --next-skill "d2a-arch-4-state-evolution" --next-file "docs/architecture/04_state_evolution.md" --summary "Started atomic-question alignment for state-evolution."`

2. Before doing analysis, explicitly tell the user:

   `Next I will scan and answer these questions: <list this skill's atomic questions>; do you want to add any? (yes/no)`

3. Allow only one supplement interaction in this phase:
   - If user answers `yes`, collect extra questions, merge them with the base atomic questions, then echo the merged list.
   - If user answers `no`, proceed with the base atomic questions.
4. Do not write `docs/architecture/04_state_evolution.md` before user confirmation.
5. After alignment is done, call:

   `d2a skill-state d2a-arch-4-state-evolution --status progress --stage architecture-in-progress --phase analysis-generation --question-index 1 --question-total 1 --next-step "Start analysis-generation with merged atomic questions." --next-skill "d2a-arch-4-state-evolution" --next-file "docs/architecture/04_state_evolution.md" --summary "Atomic-question alignment completed for state-evolution."`

## Phase 2: Analysis Generation

1. After context is confirmed, call:

   `d2a skill-state d2a-arch-4-state-evolution --status progress --stage architecture-in-progress --phase analysis-generation --next-step "Track the main state stages and transitions." --next-skill "d2a-arch-5-module-view" --next-file "docs/architecture/04_state_evolution.md" --summary "Started state-evolution analysis."`

2. Work from the real repository and write the result into `docs/architecture/04_state_evolution.md`.
3. Answer these merged atomic questions (base + optional user supplements):
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

   `d2a skill-state d2a-arch-4-state-evolution --status progress --stage architecture-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first state-evolution confirmation question." --next-skill "d2a-arch-4-state-evolution" --next-file "docs/architecture/04_state_evolution.md" --summary "State-evolution analysis complete; moving into confirmation questions."`

## Phase 3: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - tracked object or workflow
   - state stages
   - transition triggers
   - state storage / observation point
3. Ask one question per turn.
4. Before each question, keep the same envelope format and map fields as follows:
   - The `[d2a]` line must include current stage, phase, and `N/<total>` progress.
   - The `[next]` line should point to the post-question next skill/file.
5. Before asking question `N`, call:

   `d2a skill-state d2a-arch-4-state-evolution --status progress --stage architecture-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue state-evolution confirmation questions." --next-skill "d2a-arch-5-module-view" --next-file "docs/architecture/05_cooperation.md" --summary "State-evolution confirmation question <N> is active."`

6. Present one multiple-choice question.
7. Rendering format must preserve structure:
   - question stem on its own line
   - options `A.`, `B.`, `C.`, `D.` each on its own line
   - `[human_in_loop]` must be on its own line
8. Do not merge the stem and options into one line.
9. Wait for the learner answer.
10. After the learner answer:
   - say whether the answer is correct, partially correct, or incorrect
   - give one short explanation
   - continue to the next question even if the answer is wrong
11. After question 4 is evaluated:
   - output a short recap
   - output a `Comprehension Score`
   - keep the `Comprehension Score` under 100 Chinese characters
12. At the end of the confirmation phase, call:

    `d2a skill-state d2a-arch-4-state-evolution --status completed --stage architecture-in-progress --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-arch-5-module-view." --next-skill "d2a-arch-5-module-view" --next-file "docs/architecture/05_cooperation.md" --summary "Completed state-evolution confirmation questions."`

13. Confirmation-question prompts, learner answers, evaluations, and explanations must be written to `.d2a/qa/<skill>.jsonl`, and must not be written into `docs/architecture/*.md`.

## Turn-End Continuation Rule

1. Before ending each turn, call `d2a skill-state` to persist the current phase and progress.
2. End the reply with: `Continue with $d2a-step`.


## Output (Artifacts)

- A Mermaid `stateDiagram-v2` state-machine diagram only, written to `docs/architecture/04_state_evolution.md`.

## Persistence (.d2a)

- Aligned atomic question set (base + user supplements) -> `.d2a/qa/<skill>.json`.
- Confirmation prompts, learner answers, evaluations, explanations -> `.d2a/qa/<skill>.jsonl`.
- Short comprehension summary and `Comprehension Score` -> `.d2a/qa/<skill>.jsonl`.
