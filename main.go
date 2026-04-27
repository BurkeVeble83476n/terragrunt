// Package main is the entry point for Terragrunt, a thin wrapper for Terraform
// that provides extra tools for working with multiple Terraform modules.
package main

import (
	"fmt"
	"os"

	"github.com/gruntwork-io/terragrunt/cli"
	"github.com/gruntwork-io/terragrunt/errors"
	"github.com/gruntwork-io/terragrunt/shell"
	"github.com/gruntwork-io/terragrunt/util"
)

// VERSION is the current version of Terragrunt.
// This is set at build time via ldflags:
//
//	-ldflags "-X main.VERSION=x.y.z"
var VERSION = "dev"

func main() {
	// Configure the logger to write to stderr so that stdout can be used for
	// Terraform output without interference.
	logger := util.CreateLogger("")

	// Run the CLI and handle any errors.
	if err := runApp(os.Args, logger); err != nil {
		// Check if the error is a shell exit error, and if so, exit with the
		// same exit code as the underlying process.
		var exitErr *shell.ProcessExitError
		if ok := errors.As(err, &exitErr); ok {
			os.Exit(exitErr.ExitCode)
		}

		// For all other errors, print the error message and exit with code 1.
		// NOTE: printing to stderr ensures the error is visible even when stdout
		// is piped or redirected (e.g. in CI pipelines).
		fmt.Fprintf(os.Stderr, "\nError: %s\n", err.Error())
		logger.Debugf("%+v\n", err)
		os.Exit(1)
	}
}

// runApp initializes and runs the Terragrunt CLI application.
// It accepts the command-line arguments and a logger instance.
func runApp(args []string, logger *util.TerragruntLogger) error {
	// Build the CLI app with the current version.
	app := cli.CreateTerragruntCli(VERSION, os.Stdout, os.Stderr)
	return app.Run(args)
}
