# Command Role

## Principle

Inside Codex, commands should be treated as stage-transition tools used by skills.

They should not be treated as the main user workflow after lab initialization.

## Command Meaning Under Skill Control

### `d2a analyze`

Used by analysis skills to ensure the architecture docs are in the correct task-entry shape for a specific target repository.

The command should:

1. register target metadata
2. create or refresh architecture task files
3. avoid claiming that analysis is complete by itself

### `d2a derive-mini`

Used by implementation-oriented skills to ensure mini-implementation planning files exist and are aligned with the current target.

The command should:

1. validate analysis prerequisites
2. create or refresh implementation task files
3. create or refresh `src/ARCHITECTURE.md`

### `d2a test-mini`

Used by testing-oriented skills to ensure the test stage has a stable task surface.

The command should:

1. validate implementation-planning prerequisites
2. create or refresh test task files
3. create or refresh test manifests

### `d2a report`

Used by report-oriented skills to ensure the report stage has stable data files and a readable summary entry.

The command should:

1. validate report prerequisites
2. create or refresh report data
3. create or refresh `report/index.md`

## What Commands Must Not Do Alone

Without a skill or explicit higher-level logic, commands should not pretend to produce final architecture insight, final implementation code, or final teaching narrative.

That work belongs to skills operating through Codex.
