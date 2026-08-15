package main

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/domain"
	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/store"
)

type stubPlanner struct {
	plan func(context.Context, domain.ManifestRequest) (domain.Manifest, error)
}

func (s stubPlanner) Plan(ctx context.Context, request domain.ManifestRequest) (domain.Manifest, error) {
	return s.plan(ctx, request)
}

func TestRunWritesManifest(t *testing.T) {
	input := bytes.NewBufferString(`{"order_id":"order-1","customer":"Ada","packages":[{"sku":"book","units":1}]}`)
	output := &bytes.Buffer{}
	errorsOut := &bytes.Buffer{}
	planner := stubPlanner{plan: func(context.Context, domain.ManifestRequest) (domain.Manifest, error) {
		return domain.Manifest{OrderID: "order-1", Status: "planned"}, nil
	}}

	if code := run(context.Background(), input, output, errorsOut, planner); code != 0 {
		t.Fatalf("run() code = %d, stderr = %s", code, errorsOut.String())
	}
	if got := output.String(); got == "" {
		t.Fatal("run() did not write a manifest")
	}
}

func TestExitCodeForDuplicateOrder(t *testing.T) {
	err := errors.Join(errors.New("persist manifest"), store.ErrDuplicateOrder)
	if got := exitCodeFor(err); got != 2 {
		t.Fatalf("exitCodeFor() = %d, want 2", got)
	}
}

func TestRunReportsWrappedDuplicateOrder(t *testing.T) {
	input := bytes.NewBufferString(`{"order_id":"order-duplicate","customer":"Ada","packages":[{"sku":"book","units":1}]}`)
	output := &bytes.Buffer{}
	errorsOut := &bytes.Buffer{}
	planner := stubPlanner{plan: func(context.Context, domain.ManifestRequest) (domain.Manifest, error) {
		return domain.Manifest{}, errors.Join(errors.New("persist manifest"), store.ErrDuplicateOrder)
	}}

	if code := run(context.Background(), input, output, errorsOut, planner); code != 2 {
		t.Fatalf("run() code = %d, want 2", code)
	}
	if got, want := errorsOut.String(), "order already has a manifest\n"; got != want {
		t.Fatalf("run() stderr = %q, want %q", got, want)
	}
	if got := output.String(); got != "" {
		t.Fatalf("run() stdout = %q, want empty", got)
	}
}
