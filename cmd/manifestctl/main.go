package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/domain"
	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/service"
	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/store"
)

type manifestPlanner interface {
	Plan(context.Context, domain.ManifestRequest) (domain.Manifest, error)
}

func main() {
	repository := store.NewMemoryRepository()
	planner := service.NewPlanner(repository, domain.Stock{
		"book":   50,
		"lamp":   20,
		"mug":    30,
		"poster": 40,
	})
	os.Exit(run(context.Background(), os.Stdin, os.Stdout, os.Stderr, planner))
}

func run(parent context.Context, input io.Reader, output, errorOutput io.Writer, planner manifestPlanner) int {
	ctx, cancel := context.WithTimeout(parent, 2*time.Second)
	defer cancel()

	var request domain.ManifestRequest
	if err := json.NewDecoder(input).Decode(&request); err != nil {
		fmt.Fprintf(errorOutput, "invalid manifest request: %v\n", err)
		return 1
	}

	manifest, err := planner.Plan(ctx, request)
	if err != nil {
		fmt.Fprintln(errorOutput, userMessage(err))
		return exitCodeFor(err)
	}

	if err := json.NewEncoder(output).Encode(manifest); err != nil {
		fmt.Fprintf(errorOutput, "write manifest: %v\n", err)
		return 1
	}
	return 0
}

func exitCodeFor(err error) int {
	if errors.Is(err, store.ErrDuplicateOrder) {
		return 2
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return 3
	}
	return 1
}

func userMessage(err error) string {
	switch exitCodeFor(err) {
	case 2:
		return "order already has a manifest"
	case 3:
		return "manifest planning canceled"
	default:
		return fmt.Sprintf("manifest planning failed: %v", err)
	}
}
