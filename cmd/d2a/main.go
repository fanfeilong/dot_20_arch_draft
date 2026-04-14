package main

import (
	"fmt"
	"os"

	"github.com/fanfeilong/dot_20_arch_draft/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "d2a:", err)
		os.Exit(1)
	}
}
