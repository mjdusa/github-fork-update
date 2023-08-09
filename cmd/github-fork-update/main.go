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

func GetParameters() (string, bool, bool) {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fs.SetOutput(os.Stderr)

	var auth string
	var debug bool
	var verbose bool

	// add flags
	fs.StringVar(&auth, "auth", "", "GitHub Auth Token")
	fs.BoolVar(&debug, "debug", false, "Log Debug")
	fs.BoolVar(&verbose, "verbose", false, "Show Verbose Logging")

	// Parse the flags
	if err := fs.Parse(os.Args[1:]); err != nil {
		fs.Usage()
		os.Exit(2)
	}

	if len(auth) <= 0 {
		fs.Usage()
		os.Exit(2)
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
