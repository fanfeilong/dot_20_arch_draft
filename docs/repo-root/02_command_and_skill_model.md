# Command And Skill Model

## External Entry

Before entering Codex, the user only does:

1. `git clone <repo>`
2. `cd <repo>`
3. `d2a init`

## Internal Entry

After that, the user enters Codex and uses skills as the main control surface.

Typical examples:

- `$d2a-step`
- `$d2a-mini-1-scope`
- `$d2a-mini-2-design`
- `$d2a-mini-3-build`
- `$d2a-mini-4-test`
- `$d2a-report-build`

Recommended default entry:

1. start with `$d2a-step`
2. after each answer turn, continue with `$d2a-step`

This keeps user control simple while state-based routing chooses the right sub-skill.

## Command Role

After `d2a init`, commands are internal stage tools used by skills.

They should operate on the repository-local `.d2a/` directory by default.

The user should not need to keep specifying a separate working-directory path.

## Skill Location

`d2a init` should install `d2a-*` skills into the repository-local AI tool directories at the repository root, for example:

- `.codex/skills/`
- `.claude/skills/`
- `.cursor/skills/`
- `.opencode/skills/`
- `.trae/skills/`
- `.neocode/skills/`

This keeps skill discovery aligned with how the AI tools already work.

The `.d2a/` directory is not a skill-discovery directory. It is only the d2a work and state directory.

## Stage Mapping

- analysis skills may call `d2a analyze`
- implementation skills may call `d2a derive-mini`
- testing skills may call `d2a test-mini`
- reporting skills may call `d2a report`
- serving/report-preview skills may call `d2a serve`
- status skills may call `d2a status`

## Step-Orchestrator Rule

`d2a-step` is the state-driven orchestrator skill.

It should:

1. recover stage/phase/question progress from `.d2a/state.json`
2. route to the correct next `d2a-*` skill
3. support session-resume continuation after AI tool restarts
4. require sub-skills to persist `d2a skill-state` on every learner-facing turn

## Default Resolution Rule

Commands and skills should treat the current repository root as the working d2a context.

They should discover it by looking for:

- `.git/`
- `.d2a/`

in the current directory.

If the user is not at the repository root, the tool may walk upward until it finds the repository root.

## State Visibility Rule

Before doing stage work, every command and every skill should expose the current state-machine position.

This means showing:

- repository name
- repository path
- d2a path
- current stage
- next step

The state should come from `.d2a/state.json` when available.
