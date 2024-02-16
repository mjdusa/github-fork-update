package main

import (
	"context"

	"github.com/mjdusa/github-fork-update/internal/run"
)

func main() {
	ctx := context.Background()

	err := run.Run(ctx)
	if err != nil {
		panic(err)
	}
}
