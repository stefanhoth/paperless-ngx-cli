// Command paperless is the Paperless-NGX CLI entry point.
package main

import (
	_ "embed"

	"github.com/stefanhoth/paperless-ngx-cli/cmd"
)

//go:embed SKILL.md
var skillMD string

func main() {
	cmd.SkillMD = skillMD
	cmd.Execute()
}
