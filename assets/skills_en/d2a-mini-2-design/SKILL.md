---
name: d2a-mini-2-design
description: Built-in d2a skill for d2a-mini-2-design stage guidance and state updates.
---

# d2a-mini-2-design

## Goal

Turn the chosen mini scope into a concrete design for the implementation under `.d2a/src/`, then verify that the learner actually understands why this design is the right minimal form.

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

If the active repository is unknown, stop and ask the user which repository should be used.

## Human In Loop Marker Rule

When the current turn asks the user a question and waits for user input, the last line of the reply body must append:

`[human_in_loop]`

## Three Mini Fast-Track Gates

In a 1-hour talk setting, the mini stage must pass these gates before implementation detail:

1. Provider Gate (stack)
   - Detect the target stack and match a built-in provider first.
   - If matched, use the provider's minimal implementation plan by default.
2. Timebox Gate (cost)
   - Declare a mini-build budget (recommended: 20 minutes).
   - If it does not fit, reduce design complexity and module count.
3. Intent Gate (alignment)
   - Design must directly serve arch anchors, not start a new business-analysis branch.

## Phase 1: Analysis Generation

1. After context is confirmed, call:

   `d2a skill-state d2a-mini-2-design --status started --stage mini-design-in-progress --phase analysis-generation --next-step "Design the smallest useful runnable mini architecture." --next-skill "d2a-mini-3-build" --next-file "docs/implementation/01_mini_design.md" --summary "Started mini-design work."`

2. Execute the three gates first and output gate conclusions (provider match, timebox, intent anchors).
3. If implementation planning files have not yet been prepared, call `d2a derive-mini`.
4. Read:
   - `docs/implementation/00_mini_scope.md`
   - the architecture files it references
5. Write the result into `docs/implementation/01_mini_design.md`.
6. Keep `.d2a/src/ARCHITECTURE.md` aligned with the chosen design if the summary there is stale.
7. Answer these atomic questions:
   - What are the main modules of the mini version?
   - What interfaces or entry points are required?
   - What is the runtime flow of the mini version?
   - What state model must be preserved?
8. After the first draft, force three correction passes:
   - compression pass
   - de-jargon pass
   - conversational simplification pass
9. Keep the design small enough to support one first runnable slice.
10. Keep the design tied to the core architecture idea, not to broad feature coverage.
11. When the analysis draft is stable, call:

   `d2a skill-state d2a-mini-2-design --status progress --stage mini-design-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first mini-design confirmation question." --next-skill "d2a-mini-2-design" --next-file "docs/implementation/01_mini_design.md" --summary "Mini-design analysis complete; moving into confirmation questions."`

## Phase 2: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - provider-constrained design fitness
   - module/interface minimization under timebox
   - runtime flow alignment with core intent
   - state-model alignment with arch anchors
3. Ask one question per turn.
4. Before each question, keep the same envelope format and map fields as follows:
   - The `[d2a]` line must include current stage, phase, and `N/<total>` progress.
   - The `[next]` line should point to the post-question next skill/file.
5. Before asking question `N`, call:

   `d2a skill-state d2a-mini-2-design --status progress --stage mini-design-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue mini-design confirmation questions." --next-skill "d2a-mini-3-build" --next-file ".d2a/src/ARCHITECTURE.md" --summary "Mini-design confirmation question <N> is active."`

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

    `d2a skill-state d2a-mini-2-design --status completed --stage mini-design-complete --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-mini-3-build." --next-skill "d2a-mini-3-build" --next-file ".d2a/src/ARCHITECTURE.md" --summary "Completed mini-design confirmation questions."`

11. Confirmation-question prompts, learner answers, evaluations, and explanations must be written to `.d2a/qa/<skill>.jsonl`, and must not be written into `docs/implementation/*.md` or `src/*`.

## Turn-End Continuation Rule

1. Before ending each turn, call `d2a skill-state` to persist the current phase and progress.
2. End the reply with: `Continue with $d2a-step`.


## Output (Artifacts)

- Main modules.
- Main interfaces or entry points.
- Runtime flow.
- State model.
- Written to `docs/implementation/01_mini_design.md`.

## Persistence (.d2a)

- Provider/Timebox/Intent gate decisions -> `.d2a/mini_gate/d2a-mini-2-design.json`.
- Confirmation prompts, learner answers, evaluations, explanations -> `.d2a/qa/<skill>.jsonl`.
- Short comprehension summary and `Comprehension Score` -> `.d2a/qa/<skill>.jsonl`.
