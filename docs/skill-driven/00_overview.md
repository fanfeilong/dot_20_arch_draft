# Skill-Driven d2a

## Goal

Define `d2a` as a skill-driven workflow inside Codex.

After installation and `d2a init`, the user should no longer drive the workflow mainly by typing `d2a` commands by hand.

Instead:

1. The user invokes `$d2a-step` in Codex.
2. The skill decides what stage the lab is in.
3. `d2a-step` routes to the correct sub-skill for that stage.
4. The sub-skill calls `d2a` commands as internal tools when needed.
5. The sub-skill writes or updates the right files in `docs/`, `src/`, `tests/`, and `report/`.

## Interaction Model

### Before Codex

Two shell entry points remain external:

1. install `d2a`
2. run `d2a init <lab-dir>`

These are bootstrap steps.

### Inside Codex

After the lab exists, the external entry should be skills, not commands.

The user should say things like:

- "Use `$d2a-step`."
- "继续请使用 `$d2a-step`."

The skill may then call:

- `d2a analyze`
- `d2a derive-mini`
- `d2a test-mini`
- `d2a report`

but those commands are internal stage tools, not the primary user-facing flow.

## Resume Rule

After an AI coding tool session is interrupted or restarted, the user should re-enter with `$d2a-step`.

`d2a-step` should recover from `.d2a/state.json` and continue from the exact pending stage/phase/question progress.

## Core Principle

`d2a` commands should manage state transitions and filesystem contracts.

`d2a-*` skills should remain the cognitive entry point and content-production layer.
