package assets

import "embed"

// Files contains the built-in d2a skill templates and lab starter files
// distributed by the CLI.
//
//go:embed skills_zh_cn/*/SKILL.md skills_en/*/SKILL.md lab_zh_cn lab_en
var Files embed.FS
