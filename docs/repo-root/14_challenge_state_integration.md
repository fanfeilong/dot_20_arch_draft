# Challenge State Integration

## Goal

Integrate the challenge phase into the persistent d2a state machine.

## New Stages

Insert these stages after architecture completion:

1. `architecture-challenge-prepared`
2. `architecture-challenge-in-progress`
3. `architecture-challenge-complete`

Then continue to:

4. `mini-derivation-prepared`

## State Meaning

### `architecture-challenge-prepared`

The system is ready to begin the challenge phase.

### `architecture-challenge-in-progress`

The challenge dialogue is underway and challenge-turn progress must be persisted.

### `architecture-challenge-complete`

All required architecture decisions have been challenged and recorded.

## Required Files

Implemented files:

```text
.d2a/
  challenge.json
  challenge_log.jsonl
```

## Recommended Content

`challenge.json` should contain the current live challenge state:

- current decision index
- total decisions
- current decision label
- recommendation after the current round

`challenge_log.jsonl` should append completed challenge turns.

## Next Step Rule

The state machine should not allow mini derivation to begin until the challenge phase is either:

- completed
- explicitly skipped with a recorded reason

## Status

Implemented:

1. `d2a-challenge-architecture`
2. challenge-state persistence
3. challenge logs
4. status output of challenge progress during challenge stages
