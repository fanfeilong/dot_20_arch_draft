# Next Changes

## Why This Matters

The current implementation has a staged CLI, but the product intent is skill-first after initialization.

That means future changes should move toward:

1. skills that know when to call `d2a` commands
2. commands that only manage structural state
3. prompts and docs that teach the user to invoke skills, not raw commands

## Immediate Documentation Changes

Future top-level docs should explain:

1. install `d2a`
2. run `d2a init`
3. enter Codex
4. invoke the right `d2a-*` skill

They should not teach the post-init workflow as:

1. run `d2a analyze`
2. run `d2a derive-mini`
3. run `d2a test-mini`
4. run `d2a report`

because that would misplace the primary control surface.

## Current Product Status

Current implementation already includes:

1. implementation-focused skills (`d2a-mini-scope`, `d2a-mini-design`, `d2a-mini-build`)
2. testing-focused skill (`d2a-mini-test`)
3. reporting-focused skill (`d2a-report-build`)
4. explicit command-calling conventions in each skill template (`d2a skill-state` and stage commands)

## Next Changes

1. keep top-level docs skill-first after `d2a init`
2. keep commands stage/tool-oriented and avoid narrative claims
3. continue reducing any remaining lab-era wording in internal docs

## Design Rule

After initialization:

- skill is the user-facing entry
- command is the internal tool
- files are the stable state surface
