package cli

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fanfeilong/dot_20_arch_draft/internal/installer"
)

const usage = `d2a installs built-in architecture-analysis skills.

Usage:
  d2a help
  d2a init <target-dir>
`

func Run(args []string) error {
	return runWithIO(args, os.Stdout)
}

func runWithIO(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		printUsage(stdout)
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		printUsage(stdout)
		return nil
	case "init":
		if len(args) != 2 {
			return errors.New("init requires exactly one target directory")
		}

		target, err := installer.Install(args[1])
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(stdout, "initialized d2a skills in %s\n", target)
		return err
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func printUsage(stdout io.Writer) {
	fmt.Fprint(stdout, usage)
}
