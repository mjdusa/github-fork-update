package run

import (
	"context"
	"fmt"
	"os"

	"github.com/mjdusa/github-fork-update/internal/environment"
	"github.com/mjdusa/github-fork-update/internal/githubapi"
	"github.com/mjdusa/github-fork-update/internal/profile"
)

func Run(ctx context.Context) error {
	env, eerr := environment.NewEnvironment()
	if eerr != nil {
		return fmt.Errorf("NewEnvironment error: %w", eerr)
	}

	auth, debugFlag, verboseFlag, perr := env.GetParameters()
	if perr != nil {
		return fmt.Errorf("GetParameters error: %w", perr)
	}

	if *debugFlag {
		pro, merr := profile.NewProfile(ctx, "cpu-profile.pprof", "mem-profile.pprof")
		if merr != nil {
			return fmt.Errorf("NewProfile error: %w", merr)
		}

		serr := pro.StartCPUProfile()
		if serr != nil {
			fmt.Fprintf(os.Stderr, "profile StartCPUProfile error: %v\n", serr)
		}

		defer func() {
			pro.StopCPUProfile()

			werr := pro.WriteHeapProfile()
			if werr != nil {
				fmt.Fprintf(os.Stderr, "profile WriteHeapProfile error: %v\n", werr)
			}

			cerr := pro.Close()
			if cerr != nil {
				fmt.Fprintf(os.Stderr, "profile Close error: %v\n", cerr)
			}
		}()
	}

	merr := Process(ctx, auth, verboseFlag, debugFlag)
	if merr != nil {
		return fmt.Errorf("Process error: %w", merr)
	}

	return nil
}

func Process(ctx context.Context, auth *string, verboseFlag *bool, debugFlag *bool) error {
	if auth == nil {
		return fmt.Errorf("empty token error")
	}

	gapi, aerr := githubapi.NewGitHubAPI(ctx, *auth)
	if aerr != nil {
		return fmt.Errorf("NewGitHubAPI error: %w", aerr)
	}

	serr := gapi.SyncForks(ctx, "", *verboseFlag, *debugFlag)
	if serr != nil {
		return fmt.Errorf("SyncForks error: %w", serr)
	}
	return nil
}
