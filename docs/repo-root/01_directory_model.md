# Directory Model

## Repository Root

The target repository root remains the main working directory.

`d2a init` should create two classes of content:

1. repo-local AI tool skill directories at the repository root
2. d2a-owned work content under a hidden `.d2a/` directory

Example:

```text
target-repo/
  .git/
  existing-project-files...
  .codex/skills/
  .claude/skills/
  .cursor/skills/
  .opencode/skills/
  .trae/skills/
  .neocode/skills/
  .d2a/
    LAB.md
    target.json
    docs/
      architecture/
      implementation/
      report/
    src/
    tests/
    report/
```

## Separation Rule

The root-level AI tool directories exist so Codex, Claude, Cursor, OpenCode, Trae, and NeoCode can discover the installed `d2a-*` skills normally.

The `.d2a/` directory exists only for d2a-owned state and work artifacts:

- architecture docs
- implementation docs
- mini source output
- testing output
- report output
- d2a metadata

## Ownership Rule

`d2a` should not create top-level `docs/`, `src/`, `tests/`, or `report/` in the repository root.

All d2a-generated analysis and derivation content should live under `.d2a/`.

The exception is skill installation:

- `.codex/skills/`
- `.claude/skills/`
- `.cursor/skills/`
- `.opencode/skills/`
- `.trae/skills/`
- `.neocode/skills/`

These belong at the repository root because they are consumed directly by the AI coding tools.

This keeps the target repository clean and reduces accidental pollution of real project files.

## Benefits

- safer initialization inside an existing source repository
- skills remain discoverable by AI coding tools without extra indirection
- clear separation between target project code and d2a output
- easy cleanup by removing `.d2a/`
- no confusion about whether generated files belong to the original project
