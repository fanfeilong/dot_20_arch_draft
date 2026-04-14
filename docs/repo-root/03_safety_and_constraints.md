# Safety And Constraints

## Safety Goal

Allow `d2a init` to be run inside a real cloned repository without polluting the repository root.

## Safety Rule

`d2a init` should create:

1. repository-root AI tool skill directories
2. `.d2a/` and its contents

It should not write:

- `docs/`
- `src/`
- `tests/`
- `report/`

at the repository root.

Skill installation at the repository root is allowed because those directories are the intended integration points for the AI tools.

## Context Rule

After initialization, every command and every skill should report:

- current repository root
- current d2a path

For this simplified design, repository name can be the repository basename.

Example:

```text
d2a repo: n8n
d2a repo path: /abs/path/to/n8n
d2a path: /abs/path/to/n8n/.d2a
```

The d2a path is the work directory.

The installed skill paths remain at the repository root under the tool-specific directories.

## Failure Rule

If the current repository root or `.d2a/` directory cannot be determined, stage work must stop.

This avoids writing d2a output to the wrong repository.

## Current Design Implication

This simplified model replaces:

- multi-target workspace design
- active-target switching logic
- explicit separate working-directory management

with:

- one repository
- one hidden `.d2a/`
- one d2a context

## Status

CLI arguments and output headers are repository-root based.

Remaining cleanup is internal naming only (for example, some legacy `LAB.md` and helper function names), which does not affect command contracts.
