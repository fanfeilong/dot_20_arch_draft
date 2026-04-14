# d2a-challenge-architecture Skill

## Goal

Define a dedicated skill that runs the architecture challenge phase.

Suggested name:

- `d2a-challenge-architecture`

## Role

This skill should:

1. read the completed architecture analysis
2. expose each important architecture decision to challenge
3. let the learner repeatedly question those decisions
4. answer and record the objections
5. avoid mutating the architecture by default

## Coverage

The challenge skill should cover at least:

- system boundary
- primary driver
- core objects
- state evolution
- cooperation pattern
- dominant constraint

## Interaction Rule

The skill should run as repeated rounds:

1. present one architecture decision
2. ask the learner to challenge it
3. accept the learner's objection
4. answer the objection
5. record whether the objection is strong, partial, or weak
6. move to the next decision

## Non-Mutation Rule

The skill must not directly revise the architecture docs during the challenge phase.

It should only:

- answer
- evaluate
- record

If a challenge seems valid, it should be marked for later review instead of silently changing the architecture output.

## Output

The skill should produce:

- challenge log
- challenge strength assessment
- unresolved questions
- recommendation: proceed, review, or revisit architecture
