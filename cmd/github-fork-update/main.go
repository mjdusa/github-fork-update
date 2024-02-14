package main

import (
	"context"

	"github.com/mjdusa/github-fork-update/internal/run"
)

func main() {
	ctx := context.Background()

	run.Run(ctx)
}
