---
name: d2a-arch-7-overview
description: Built-in d2a skill for d2a-arch-7-overview stage guidance and state updates.
---

# d2a-arch-7-overview

## Goal

Synthesize docs 01-06 into `00_overview`, then hand off to challenge phase.

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
3. Do not print the old multi-line start header list anymore.
4. If `next_file` is unknown, print `unknown` and keep the same shape.


## Body Formatting Hard Rules

1. The body must use `- ` bullet lines, minimum 2 lines and maximum 4 lines.
2. Each point must occupy exactly one line; do not write long paragraphs.
3. Keep each line under 80 characters; split lines when longer.
4. Avoid Markdown emphasis in body text.

5. For structured content (lists, numbered items, MCQ options A/B/C/D), do not merge lines just to satisfy width limits; keep one item per line.

If the active repository is unknown, stop and ask the user which repo to use.

## Human In Loop Marker Rule

When the current turn asks the user a question and waits for user input, the
last line of the reply body must append:

`[human_in_loop]`

## Inputs

Read these files before starting:

- `docs/architecture/01_boundary.md`
- `docs/architecture/02_driver.md`
- `docs/architecture/03_core_objects.md`
- `docs/architecture/04_state_evolution.md`
- `docs/architecture/05_cooperation.md`
- `docs/architecture/06_constraints.md`

## Phase 1: Atomic Question Alignment (One Supplement Chance)

1. After context is confirmed, call:

   `d2a skill-state d2a-arch-7-overview --status started --stage architecture-in-progress --phase atomic-question-alignment --question-index 0 --question-total 1 --next-step "Show overview synthesis questions and ask for one-time supplements (yes/no)." --next-skill "d2a-arch-7-overview" --next-file "docs/architecture/00_overview.md" --summary "Started overview synthesis question alignment."`

2. Before doing analysis, explicitly tell the user:

   `Next I will synthesize 01-06 and answer these questions: <list atomic questions>; do you want to add any? (yes/no)`

3. Allow only one supplement interaction in this phase.
4. Do not write `docs/architecture/00_overview.md` before user confirmation.
5. After alignment is done, call:

   `d2a skill-state d2a-arch-7-overview --status progress --stage architecture-in-progress --phase analysis-generation --question-index 1 --question-total 1 --next-step "Start overview synthesis with merged atomic questions." --next-skill "d2a-arch-7-overview" --next-file "docs/architecture/00_overview.md" --summary "Overview synthesis question alignment completed."`

## Phase 2: Analysis Generation

1. After context is confirmed, call:

   `d2a skill-state d2a-arch-7-overview --status progress --stage architecture-in-progress --phase analysis-generation --next-step "Synthesize 01-06 into overview conclusions." --next-skill "d2a-challenge-architecture" --next-file "docs/architecture/00_overview.md" --summary "Started overview synthesis analysis."`

2. Synthesize 01-06 and write the result into `docs/architecture/00_overview.md`.
3. Answer these merged atomic questions:
   - What is the one-sentence system definition?
   - What capability must remain if 80% of code is deleted?
   - What architecture idea is most worth preserving?
   - What are the first four things a reader should understand?
4. After the first draft, force three correction passes:
   - compression pass
   - de-jargon pass
   - conversational simplification pass
5. The overview must cite 01-06 conclusions and stay consistent with them.
6. When the draft is stable, call:

   `d2a skill-state d2a-arch-7-overview --status progress --stage architecture-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first overview confirmation question." --next-skill "d2a-arch-7-overview" --next-file "docs/architecture/00_overview.md" --summary "Overview synthesis complete; moving into confirmation questions."`

## Phase 3: Confirmation Questions

1. Generate `4` multiple-choice questions from phase-1 output.
   - Distractors must be plausible and close to real implementation paths or common misconceptions; avoid obviously absurd options.
   - The correct option and distractors should be lexically similar enough to be confusing; correctness should require understanding this project's conclusions, not keyword spotting.
   - At least `2` distractors per question must reuse real concepts, modules, or flow names from this project, but with semantically incorrect mapping.
2. Coverage should include:
   - one-sentence definition
   - non-removable capability
   - core architecture idea
   - first-four understanding points
3. Ask one question per turn.
4. Before each question, keep the same envelope format and map fields.
5. Before asking question `N`, call:

   `d2a skill-state d2a-arch-7-overview --status progress --stage architecture-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue overview confirmation questions." --next-skill "d2a-challenge-architecture" --next-file "docs/architecture/00_overview.md" --summary "Overview confirmation question <N> is active."`

6. Present one multiple-choice question.
7. Rendering format must preserve structure:
   - question stem on its own line
   - options `A.`, `B.`, `C.`, `D.` each on its own line
   - `[human_in_loop]` must be on its own line
8. Do not merge the stem and options into one line.
9. Wait for the learner answer.
7. Evaluate answer and continue even when wrong.
8. After question 4, output a short recap and `Comprehension Score`.
9. At phase end, call:

   `d2a skill-state d2a-arch-7-overview --status completed --stage architecture-complete --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-challenge-architecture." --next-skill "d2a-challenge-architecture" --next-file "docs/architecture/00_overview.md" --summary "Completed overview confirmation questions."`

10. Confirmation content must go to `.d2a/qa/<skill>.jsonl`.

## Turn-End Continuation Rule

1. Before ending each turn, call `d2a skill-state`.
2. End the reply with: `Continue with $d2a-step`.


## Output (Artifacts)

- Stage synthesis conclusion written into `docs/architecture/00_overview.md`.

## Persistence (.d2a)

- Aligned atomic questions -> `.d2a/qa/<skill>.json`.
- Confirmation records -> `.d2a/qa/<skill>.jsonl`.
