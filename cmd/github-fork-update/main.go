package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-github/v53/github"
	"github.com/mjdusa/github-fork-update/internal/tools"
	"github.com/mjdusa/github-fork-update/internal/version"
)

func GetUsage() string {
	msg := "usage:\n"
	msg += fmt.Sprintf("\t%s -auth='github-auth-token' [-verbose]\n\n", os.Args[0])

	return msg
}

func GetParameters() (string, bool) {
	// Define flags
	token := flag.String("auth", "", "GitHub Auth Token")
	verbose := flag.Bool("verbose", false, "Verbose")

	// Parse the flags
	flag.Parse()

	return *token, *verbose
}

func main() {
	ctx := context.Background()
	token, verbose := GetParameters()

	if verbose {
		fmt.Println(version.GetVersion())
	}

	if len(token) == 0 {
		fmt.Println(GetUsage())
		return
	}

	client := github.NewTokenClient(ctx, token)

	err := tools.SyncForks(ctx, client, "", verbose)
	if err != nil {
		panic(err)
	}
}
