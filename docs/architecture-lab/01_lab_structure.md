# Repository Structure

## Target Layout

`d2a init <target-repo-git-url>` should create this structure:

```text
<repo-dir>/
  .codex/skills/
  .claude/skills/
  .cursor/skills/
  .opencode/skills/
  .trae/skills/
  .neocode/skills/
  docs/
    architecture/
      00_overview.md
      01_boundary.md
      02_driver.md
      03_core_objects.md
      04_state_evolution.md
      05_cooperation.md
      06_constraints.md
      99_code_map.md
    implementation/
      00_mini_scope.md
      01_mini_design.md
      02_build_plan.md
      03_test_plan.md
    report/
      00_report_outline.md
  src/
    README.md
  tests/
    README.md
  report/
    README.md
  LAB.md
```

## Directory Responsibilities

- `docs/architecture/`
  Holds the six-essentials architecture decomposition and code map.
- `docs/implementation/`
  Holds the mini-clone scope, design, build plan, and test plan.
- `docs/report/`
  Holds the intended report structure and presentation notes.
- `src/`
  Holds the runnable mini implementation.
- `tests/`
  Holds integration tests that validate the mini implementation.
- `report/`
  Holds the future Vue app and report assets.
- `LAB.md`
  Acts as the local workflow entry document under `.d2a/`.

## Initialization Rule

The generated files should be templates, not final content.

Each template should:

1. State its purpose.
2. State which skill or future command should fill it.
3. Keep sections short and numbered.
