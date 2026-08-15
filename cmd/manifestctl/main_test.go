package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/domain"
	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/service"
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

func TestRunWritesDefaultDeliveryWindowForFlexibleDelivery(t *testing.T) {
	input := bytes.NewBufferString(`{"order_id":"order-2","customer":"Ada","packages":[{"sku":"book","units":1}]}`)
	output := &bytes.Buffer{}
	errorsOut := &bytes.Buffer{}
	planner := service.NewPlanner(store.NewMemoryRepository(), domain.Stock{"book": 1})

	if code := run(context.Background(), input, output, errorsOut, planner); code != 0 {
		t.Fatalf("run() code = %d, stderr = %s", code, errorsOut.String())
	}

	var manifest domain.Manifest
	if err := json.NewDecoder(output).Decode(&manifest); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if manifest.DeliveryWindow == nil ||
		manifest.DeliveryWindow.Start != "09:00" ||
		manifest.DeliveryWindow.End != "18:00" {
		t.Fatalf("run() delivery window = %#v", manifest.DeliveryWindow)
	}
}

func TestExitCodeForDuplicateOrder(t *testing.T) {
	err := errors.Join(errors.New("persist manifest"), store.ErrDuplicateOrder)
	if got := exitCodeFor(err); got != 2 {
		t.Fatalf("exitCodeFor() = %d, want 2", got)
	}
}
