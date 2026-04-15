---
name: d2a-mini-3-build
description: Built-in d2a skill for d2a-mini-3-build stage guidance and state updates.
---

# d2a-mini-3-build

## Goal

Implement the first runnable mini version in `.d2a/src/` while preserving the target project's core architecture idea, then verify that the learner actually understands why this implementation slice is sufficient.

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

## Three Mini Fast-Track Gates

In a 1-hour talk setting, the mini stage must pass these gates before implementation detail:

1. Provider Gate (stack)
   - Detect the target stack and match a built-in provider first.
   - If matched, use the provider's minimal template and file layout directly.
2. Timebox Gate (cost)
   - Declare a mini-build budget (recommended: 20 minutes).
   - If it cannot fit, downgrade immediately to a single happy path.
3. Intent Gate (alignment)
   - Build only to prove architecture intent, not to expand business scope.
   - Code must stay on arch anchors: object/state/cooperation chain.

## Phase 1: Analysis Generation

1. After context is confirmed, call:

   `d2a skill-state d2a-mini-3-build --status started --stage mini-design-complete --phase analysis-generation --next-step "Implement the first runnable mini slice." --next-skill "d2a-mini-4-test" --next-file ".d2a/src/ARCHITECTURE.md" --summary "Started mini-build work."`

2. Execute the three gates first and output gate conclusions (provider match, timebox, intent anchors).
3. If implementation planning files have not yet been prepared, call `d2a derive-mini`.
4. Read:
   - `docs/implementation/00_mini_scope.md`
   - `docs/implementation/01_mini_design.md`
   - `docs/implementation/02_build_plan.md`
   - `.d2a/src/ARCHITECTURE.md`
5. Implement only the first runnable slice described in the build plan.
6. Prefer a small but executable result over broad coverage.
7. Keep the chosen stack aligned with the original project when practical.
8. Update `.d2a/src/ARCHITECTURE.md` if implementation reality forces a design correction.
9. Do not expand scope until the first runnable slice is working.
10. After the implementation is stable, write a brief note on what remains intentionally unimplemented.
11. When the implementation pass is stable, call:

   `d2a skill-state d2a-mini-3-build --status progress --stage mini-design-complete --phase confirmation-questions --question-index 0 --question-total 4 --next-step "Ask the first mini-build confirmation question." --next-skill "d2a-mini-3-build" --next-file ".d2a/src/ARCHITECTURE.md" --summary "Mini-build implementation complete; moving into confirmation questions."`

## Phase 2: Confirmation Questions

1. Generate `4` multiple-choice questions from the actual phase-1 output, not from generic examples.
2. The 4 questions should cover these angles:
   - provider template and stack execution
   - single-happy-path tradeoff under timebox
   - intent-anchor proof (object/state/chain)
   - intentional omissions and scope control
3. Ask one question per turn.
4. Before each question, keep the same envelope format and map fields as follows:
   - The `[d2a]` line must include current stage, phase, and `N/<total>` progress.
   - The `[next]` line should point to the post-question next skill/file.
5. Before asking question `N`, call:

   `d2a skill-state d2a-mini-3-build --status progress --stage mini-design-complete --phase confirmation-questions --question-index <N> --question-total 4 --next-step "Continue mini-build confirmation questions." --next-skill "d2a-mini-4-test" --next-file ".d2a/tests/README.md" --summary "Mini-build confirmation question <N> is active."`

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

    `d2a skill-state d2a-mini-3-build --status completed --stage mini-design-complete --phase confirmation-questions --question-index 4 --question-total 4 --next-step "Move to d2a-mini-4-test." --next-skill "d2a-mini-4-test" --next-file ".d2a/tests/README.md" --summary "Completed mini-build confirmation questions."`

13. Confirmation-question prompts, learner answers, evaluations, and explanations must be written to `.d2a/qa/<skill>.jsonl`, and must not be written into `docs/implementation/*.md` or `src/*`.

## Turn-End Continuation Rule

1. Before ending each turn, call `d2a skill-state` to persist the current phase and progress.
2. End the reply with: `Continue with $d2a-step`.


## Output (Artifacts)

- Runnable mini implementation (written to `src/*`).
- Updated `src/ARCHITECTURE.md`.
- Brief intentional-omission note (written to `docs/implementation/02_build_plan.md` or implementation notes in `src/ARCHITECTURE.md`).

## Persistence (.d2a)

- Provider/Timebox/Intent gate decisions -> `.d2a/mini_gate/d2a-mini-3-build.json`.
- Confirmation prompts, learner answers, evaluations, explanations -> `.d2a/qa/<skill>.jsonl`.
- Short comprehension summary and `Comprehension Score` -> `.d2a/qa/<skill>.jsonl`.
