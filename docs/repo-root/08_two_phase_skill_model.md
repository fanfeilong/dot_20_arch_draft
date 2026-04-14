# Two-Phase Skill Model

## Goal

Redefine every `d2a-*` skill as a two-phase skill.

The core product goal is not only to analyze a repository architecture, but also to ensure that the learner actually understands each analysis result.

## Two Required Phases

Every `d2a-*` skill should contain:

1. analysis-generation phase
2. confirmation-question phase

## Phase 1: Analysis Generation

This phase produces the stage result.

It should use:

- atomic questions
- analysis
- de-jargonization
- conversational simplification
- compression

The output should be the smallest useful explanation for the current stage.

## Phase 2: Confirmation Questions

This phase verifies whether the learner actually understands the phase-1 output.

It should:

1. derive multiple-choice questions from phase-1 output
2. cover multiple perspectives of the analysis
3. ask one question at a time
4. evaluate the answer after each question
5. continue regardless of whether the answer is correct or incorrect
6. finish with a short comprehension summary

## Design Rule

A skill is not complete when it only writes analysis.

A skill is complete only when it has:

1. produced the analysis
2. checked understanding
3. updated workflow state

## Pre-Mini Challenge Exception

Before mini implementation begins, there should be an additional challenge phase.

This challenge phase is separate from the normal confirmation-question phase.

Its purpose is not to check simple recall.

Its purpose is to pressure-test the architecture by letting the learner question it without mutating the architecture output by default.
