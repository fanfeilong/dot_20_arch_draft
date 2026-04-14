package assets

import "embed"

// Files contains the built-in d2a skill templates distributed by the CLI.
//
//go:embed skills/*/SKILL.md
var Files embed.FS
