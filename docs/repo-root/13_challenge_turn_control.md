# Challenge Turn Control

## Goal

Define how the challenge phase works as a multi-turn stateful interaction.

## Required Per-Turn Header

Each challenge turn should begin with:

- repository name
- repository path
- d2a path
- current stage
- current phase
- challenge progress
- current architecture decision under challenge
- next step after the challenge set

Example:

```text
d2a repo: n8n
d2a repo path: /abs/path/to/n8n
d2a path: /abs/path/to/n8n/.d2a
d2a stage: architecture-challenge-in-progress
d2a phase: challenge-dialogue
d2a challenge progress: 3/6
d2a current decision: primary driver
d2a next step: finish challenge phase, then start mini scope selection
```

## Turn Rule

Each round should contain:

1. one architecture decision
2. one learner objection
3. one AI answer
4. one evaluation of objection strength

Then move to the next architecture decision.

## Evaluation Rule

The AI should classify the learner objection as:

- strong
- partial
- weak

Then provide a short explanation.

It should not say only "correct" or "incorrect".

## Persistence Rule

The state machine should record:

- current challenge index
- total challenge decisions
- current decision name
- last learner objection
- last evaluation strength
- whether review is recommended later
