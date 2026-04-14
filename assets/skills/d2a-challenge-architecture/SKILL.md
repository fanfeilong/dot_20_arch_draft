# d2a-challenge-architecture

## Goal

Pressure-test the completed architecture through repeated learner objections, then record the challenge outcome without mutating the architecture docs by default.

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

## Inputs

Read these files before starting:

- `.d2a/docs/architecture/00_overview.md`
- `.d2a/docs/architecture/01_boundary.md`
- `.d2a/docs/architecture/02_driver.md`
- `.d2a/docs/architecture/03_core_objects.md`
- `.d2a/docs/architecture/04_state_evolution.md`
- `.d2a/docs/architecture/05_cooperation.md`
- `.d2a/docs/architecture/06_constraints.md`
- `.d2a/docs/architecture/99_code_map.md`

## Phase 1: Challenge Preparation

1. After context is confirmed, call:

   `d2a skill-state d2a-challenge-architecture --status started --stage architecture-challenge-prepared --phase challenge-preparation --next-step "Prepare the six architecture decisions for challenge dialogue." --next-skill "d2a-challenge-architecture" --next-file ".d2a/docs/architecture/00_overview.md" --summary "Started architecture challenge preparation."`

2. Extract these `6` architecture decisions to challenge:
   - system boundary
   - primary driver
   - core objects
   - state evolution
   - cooperation pattern
   - dominant constraint / tradeoff
3. Do not rewrite any architecture file in this phase.
4. If needed, create a temporary challenge checklist in the conversation, but keep the architecture docs unchanged.
5. When the challenge set is ready, call:

   `d2a skill-state d2a-challenge-architecture --status progress --stage architecture-challenge-in-progress --phase challenge-dialogue --question-index 0 --question-total 6 --next-step "Start the first architecture challenge round." --next-skill "d2a-challenge-architecture" --next-file ".d2a/docs/architecture/00_overview.md" --summary "Architecture challenge set prepared; moving into challenge dialogue."`

## Phase 2: Challenge Dialogue

1. Run `6` rounds, one per architecture decision.
2. The order should be:
   - boundary
   - driver
   - core objects
   - state evolution
   - cooperation
   - dominant constraint
3. Each round must begin with a compact header:
   - `d2a repo: ...`
   - `d2a repo path: ...`
   - `d2a path: ...`
   - `d2a stage: architecture-challenge-in-progress`
   - `d2a phase: challenge-dialogue`
   - `d2a challenge progress: N/6`
   - `d2a current decision: <decision>`
   - `d2a next step after challenge: finish challenge phase, then start mini scope selection`
4. Before challenge round `N`, call:

   `d2a skill-state d2a-challenge-architecture --status progress --stage architecture-challenge-in-progress --phase challenge-dialogue --question-index <N> --question-total 6 --next-step "Continue architecture challenge dialogue." --next-skill "d2a-mini-scope" --next-file ".d2a/docs/implementation/00_mini_scope.md" --summary "Architecture challenge round <N> is active."`

5. In each round:
   - present one architecture decision
   - ask the learner to challenge it
   - accept the learner objection
   - answer the objection
   - classify the objection as `strong`, `partial`, or `weak`
   - explain the classification briefly
6. The AI must not directly revise the architecture docs during challenge dialogue.
7. If a challenge is strong, mark it for later review rather than silently editing architecture output.
8. Continue even when the learner objection is weak.

## Phase 3: Challenge Wrap-Up

1. After round 6:
   - summarize which objections were strong, partial, or weak
   - list unresolved questions if any
   - give one recommendation:
     - `proceed`
     - `review`
     - `revisit architecture`
2. Keep the architecture docs unchanged during this wrap-up.
3. At the end of the challenge phase, call:

   `d2a skill-state d2a-challenge-architecture --status completed --stage architecture-challenge-complete --phase challenge-dialogue --question-index 6 --question-total 6 --next-step "Proceed to d2a-mini-scope unless a review is required." --next-skill "d2a-mini-scope" --next-file ".d2a/docs/implementation/00_mini_scope.md" --summary "Completed architecture challenge phase."`

## Output

- Challenge rounds
- Objection strength assessment
- Unresolved questions
- Recommendation: proceed, review, or revisit architecture
