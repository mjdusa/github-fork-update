package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/google/go-github/v53/github"
	"github.com/mjdusa/github-fork-update/internal/githubapi"
	"github.com/mjdusa/github-fork-update/internal/version"
)

const (
	OsExistCode int = 1
)

var PanicOnExit bool = false // Set to true to tell Exit() to Panic rather than call os.Exit() - should ONLY be used for testing

func Exit(code int) {
	if PanicOnExit {
		panic(fmt.Sprintf("PanicOnExit is true, code=%d", code))
	}

	os.Exit(code)
}

func GetParameters() (string, bool, bool) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flagSet.SetOutput(os.Stderr)

	var auth string
	var debug bool
	var verbose bool

	// add flags
	flagSet.StringVar(&auth, "auth", "", "GitHub Auth Token")
	flagSet.BoolVar(&debug, "debug", false, "Log Debug")
	flagSet.BoolVar(&verbose, "verbose", false, "Show Verbose Logging")

	// Parse the flags
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		flagSet.Usage()
		Exit(OsExistCode)
	}

	if len(auth) == 0 {
		flagSet.Usage()
		Exit(OsExistCode)
	}

	return auth, debug, verbose
}

func main() {
	ctx := context.Background()
	auth, debugFlag, verboseFlag := GetParameters()

	if verboseFlag {
		fmt.Println(version.GetVersion())
	}

	if debugFlag {
		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			fmt.Println(buildInfo.String())
		}
	}

	client := github.NewTokenClient(ctx, auth)

	err := githubapi.SyncForks(ctx, client, "", verboseFlag, debugFlag)
	if err != nil {
		panic(err)
	}
}