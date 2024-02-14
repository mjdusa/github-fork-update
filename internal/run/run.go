package run

import (
	"context"

	"github.com/google/go-github/v53/github"
	"github.com/mjdusa/github-fork-update/internal/env"
	gapi "github.com/mjdusa/github-fork-update/internal/githubapi"
	"github.com/mjdusa/github-fork-update/internal/profile"
)

func Run(ctx context.Context) {
	auth, debugFlag, verboseFlag, err := env.GetParameters()
	if err != nil {
		panic(err)
	}

	if *debugFlag {
		pfl, err := profile.NewProfile("cpu-profile.pprof", "mem-profile.pprof")
		if err != nil {
			panic(err)
		}
		defer pfl.Close()

		serr := pfl.StartCPUProfile()
		if serr != nil {
			panic(serr)
		}
		defer pfl.StopCPUProfile()

		defer pfl.WriteHeapProfile()
	}

	client := github.NewTokenClient(ctx, *auth)

	serr := gapi.SyncForks(ctx, client, "", *verboseFlag, *debugFlag)
	if serr != nil {
		panic(serr)
	}
}
