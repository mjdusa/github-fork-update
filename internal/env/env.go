package env

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/mjdusa/github-fork-update/internal/version"
)

// GetParameters returns the command line parameters with basic go flags.
func GetParameters() (*string, *bool, *bool, error) {
	app := ""
	if len(os.Args) > 0 {
		app = os.Args[0]
	}

	args := []string{}
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	flagSet := flag.NewFlagSet(app, flag.ContinueOnError)

	flagSet.SetOutput(os.Stderr)

	var auth string
	var dbg bool
	var verbose bool

	// add flags
	flagSet.StringVar(&auth, "auth", "", "GitHub Auth Token")
	flagSet.BoolVar(&dbg, "debug", false, "Log Debug")
	flagSet.BoolVar(&verbose, "verbose", false, "Show Verbose Logging")

	// Parse the flags
	if err := flagSet.Parse(args); err != nil {
		return nil, nil, nil, fmt.Errorf("error parsing flags: %w", err)
	}

	if len(auth) == 0 {
		return nil, nil, nil, fmt.Errorf("error missing auth token")
	}

	return &auth, &dbg, &verbose, nil
}

func Report(verbose bool, dbg bool) string {
	rpt := ""

	if verbose {
		rpt += version.GetVersion()
	}

	if dbg {
		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			rpt += buildInfo.String()
		}
	}

	return rpt
}
