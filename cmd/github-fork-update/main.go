package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"runtime/pprof"

	"github.com/google/go-github/v53/github"
	"github.com/mjdusa/github-fork-update/internal/githubapi"
	"github.com/mjdusa/github-fork-update/internal/version"
)

const (
	OsExistCode int = 1
)

var _panicOnExit = false //nolint:gochecknoglobals // Set to true to panic during unit tests

func Exit(code int) {
	if _panicOnExit {
		panic(fmt.Sprintf("_panicOnExit is true, code=%d", code))
	}

	os.Exit(code)
}

func GetParameters() (string, bool, bool) {
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
		Exit(OsExistCode)
	}

	if len(auth) == 0 {
		flagSet.Usage()
		Exit(OsExistCode)
	}

	if verbose {
		fmt.Println(version.GetVersion())
	}

	if dbg {
		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			fmt.Println(buildInfo.String())
		}
	}

	return auth, dbg, verbose
}

func main() {
	ctx := context.Background()
	auth, debugFlag, verboseFlag := GetParameters()

	var cpuFile *os.File
	var memFile *os.File
	var cerr error
	var merr error

	if debugFlag {
		cpuFile, cerr = os.Create("cpu-profile.pprof")
		if cerr != nil {
			panic(cerr)
		}
		defer cpuFile.Close()

		memFile, merr = os.Create("mem-profile.pprof")
		if merr != nil {
			panic(merr)
		}
		defer memFile.Close()

		serr := pprof.StartCPUProfile(cpuFile)
		if serr != nil {
			panic(serr)
		}
		defer pprof.StopCPUProfile()
	}

	client := github.NewTokenClient(ctx, auth)

	serr := githubapi.SyncForks(ctx, client, "", verboseFlag, debugFlag)
	if debugFlag {
		werr := pprof.WriteHeapProfile(memFile)
		if werr != nil {
			panic(werr)
		}
	}
	if serr != nil {
		panic(serr)
	}
}
