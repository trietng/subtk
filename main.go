package main

import (
	"os"
	"strings"
	"trietng/subtk/cli"
	"trietng/subtk/cli/flags"
	"trietng/subtk/cli/module"
)

func main() {
	lowered := strings.ToLower(os.Args[1])
	if len(os.Args) < 2 || strings.HasPrefix(lowered, "-") || lowered == "console" {
		flags.SetModuleFlags(module.Console)
		cli.Run(module.Console)
	} else {
		os.Args = os.Args[1:]
		flags.SetModuleFlags(lowered)
		cli.Run(lowered)
	}
}