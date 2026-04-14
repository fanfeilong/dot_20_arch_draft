# Confirmation Questions

## Goal

Define the comprehension-confirmation phase used after each analysis phase.

## Question Generation Rule

Questions should be generated from the actual phase-1 result, not from generic templates alone.

They should cover the key angles of the current skill, such as:

- system type
- driver
- core objects
- state evolution
- module cooperation
- constraints
- mini implementation choices
- testing intent
- report structure

## Interaction Rule

The confirmation phase should ask one multiple-choice question per turn.

Each question should:

1. include the current state-machine progress indicator
2. include the current repository and d2a path metadata
3. present one question
4. offer multiple choices
5. wait for the learner's answer

## Evaluation Rule

After the learner answers:

1. evaluate the answer
2. say whether it is correct, partially correct, or incorrect
3. give a short explanation
4. move to the next question

The confirmation phase must continue even when the learner answers incorrectly.

The point is to reinforce understanding, not to block progress.

## Final Output Rule

After all questions are complete, the skill should output:

- a short recap of what was understood well
- the main misunderstanding if any
- a `理解度打分`

The `理解度打分` summary must stay under 100 Chinese characters.

## Coverage Rule

The generated questions should cover the major viewpoints of the phase-1 output, not only the easiest one.
