# MVP Plan

## Phase 1

Upgrade `d2a init` from a skill-only initializer into a lab initializer.

Deliverables:

1. Create the full lab directory tree.
2. Create starter templates in `docs/`, `src/`, `tests/`, `report/`, and `LAB.md`.
3. Preserve the current skill installation behavior.
4. Update tests to verify the new structure.

## Phase 2

Teach `d2a` to write architecture docs deterministically.

Deliverables:

1. Define file-writing conventions for `docs/architecture/`.
2. Add commands or scripted flows that populate those files.
3. Ensure each output is short, numbered, and code-anchored.

## Phase 3

Derive a runnable mini implementation in the target stack.

Deliverables:

1. Define the mini-clone input contract from `docs/architecture/`.
2. Add implementation-oriented skills and templates.
3. Add a first supported stack, likely Go or TypeScript.

## Phase 4

Build the report layer.

Deliverables:

1. Define report data files.
2. Add a Vue app skeleton under `report/`.
3. Add a local serve command.

## Current Status

The implementation now covers Phase 1 through a baseline of Phase 4:

- Phase 1: full `.d2a/` tree + built-in skills
- Phase 2: deterministic architecture task-file generation (`d2a analyze`)
- Phase 3: deterministic mini/test planning files + first runnable Go scaffold under `src/go-mini/`
- Phase 4: report data files + static report surface + Vue app skeleton under `report/vue-app/`

The next frontier is quality depth, not missing scaffolding.
