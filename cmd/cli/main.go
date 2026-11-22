// Command corven is the Corven command-line client.
//
// Only the version command exists for now; deploy, logs, env and rollback arrive in later tasks.
package main

import (
	"fmt"
	"io"
	"os"

	"gitlab.com/corvenc/corven/internal/version"
)

const binary = "corven"

const usage = `Usage: corven <command>

Commands:
  version   Print version information
  help      Show this help
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}

	switch args[0] {
	case "version", "--version", "-v":
		fmt.Fprintln(stdout, version.String(binary))
		return 0
	case "help", "--help", "-h":
		fmt.Fprint(stdout, usage)
		return 0
	default:
		fmt.Fprintf(stderr, "corven: unknown command %q\n\n%s", args[0], usage)
		return 2
	}
}
