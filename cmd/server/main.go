// Command corven-server runs the Corven platform on a server.
//
// It is a placeholder for now: the deploy engine, API and panel arrive in later tasks.
package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.com/corvenc/corven/internal/version"
)

const binary = "corven-server"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	fs := flag.NewFlagSet(binary, flag.ContinueOnError)
	showVersion := fs.Bool("version", false, "print version information and exit")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *showVersion {
		fmt.Println(version.String(binary))
		return 0
	}

	fmt.Fprintln(os.Stderr, binary+": nothing to run yet, see openspec/changes/add-v0-1-core/tasks.md")
	return 1
}
