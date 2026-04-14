package assets

import "embed"

// Files contains the built-in d2a skill templates and lab starter files
// distributed by the CLI.
//
//go:embed skills/*/SKILL.md lab
var Files embed.FS
