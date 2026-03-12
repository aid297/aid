package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	os.Exit(runCLI(os.Args[1:], os.Stdout, os.Stderr))
}

func runCLI(args []string, stdout, stderr *os.File) int {
	if len(args) == 0 {
		return runServe(args, stdout, stderr)
	}

	switch args[0] {
	case "serve":
		return runServe(args[1:], stdout, stderr)
	case "print-config":
		return runPrintConfig(args[1:], stdout, stderr)
	case "init-config":
		return runInitConfig(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		printUsage(stdout)
		return 0
	default:
		return runServe(args, stdout, stderr)
	}
}

func printUsage(stdout *os.File) {
	_, _ = fmt.Fprintln(stdout, "simpleDB CLI")
	_, _ = fmt.Fprintln(stdout, "commands:")
	_, _ = fmt.Fprintln(stdout, "  serve        start HTTP service from config file")
	_, _ = fmt.Fprintln(stdout, "  print-config print merged config")
	_, _ = fmt.Fprintln(stdout, "  init-config  write default config file")
	_, _ = fmt.Fprintln(stdout, "")
	_, _ = fmt.Fprintf(stdout, "default config path: %s\n", defaultConfigPath())
}

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	return fs
}
