# Skill Flow

## End-To-End Flow In Codex

After `d2a init <repo-dir>`, the intended flow is:

1. User enters Codex in the target repository root.
2. User invokes `$d2a-step`.
3. `d2a-step` recovers current state and routes to the correct sub-skill.
4. The routed skill executes one turn and persists `d2a skill-state`.
5. The routed skill ends with `继续请使用 $d2a-step`.
6. User invokes `$d2a-step` again to continue.
7. The loop repeats until `report-ready`.

## Recommended Skill Entry Points

### Default

- `$d2a-step`

### Analysis (sub-skills routed by `d2a-step`)

- `$d2a-architecture-walkthrough`
- `$d2a-project-scope`
- `$d2a-runtime-view`
- `$d2a-core-objects`
- `$d2a-state-evolution`
- `$d2a-module-view`
- `$d2a-tradeoff-view`

### Implementation (sub-skills routed by `d2a-step`)

- `$d2a-mini-scope`
- `$d2a-mini-design`
- `$d2a-mini-build`

### Testing (sub-skills routed by `d2a-step`)

- `$d2a-mini-test`

### Reporting (sub-skills routed by `d2a-step`)

- `$d2a-report-build`

## File Ownership By Skill Stage

- Analysis skills own `docs/architecture/*`
- Implementation skills own `docs/implementation/*` and later `src/*`
- Testing skills own `tests/*`
- Reporting skills own `report/*`

## Success Condition

The user should feel that they are working through a set of skills, while `d2a` commands quietly keep repository state in the right structural stage underneath.

The user should not need to manually choose among many sub-skills in normal operation.
