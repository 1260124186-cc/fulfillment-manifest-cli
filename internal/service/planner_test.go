package service

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/domain"
	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/store"
)

func TestPlannerCreatesManifest(t *testing.T) {
	repository := store.NewMemoryRepository()
	planner := NewPlanner(repository, domain.Stock{"book": 5})

	manifest, err := planner.Plan(context.Background(), domain.ManifestRequest{
		OrderID:  "order-1",
		Customer: "Ada",
		Packages: []domain.PackageRequest{{SKU: "book", Units: 2}},
	})
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}

	if diff := cmp.Diff([]domain.Reservation{{SKU: "book", Units: 2}}, manifest.Reservations); diff != "" {
		t.Fatalf("Plan() reservations mismatch (-want +got):\n%s", diff)
	}
	stored, ok := repository.Get("order-1")
	if !ok || stored.Status != "planned" {
		t.Fatalf("Get() = %#v, %t", stored, ok)
	}
}

func TestPlannerKeepsRequestAndRetrievedManifestIndependent(t *testing.T) {
	repository := store.NewMemoryRepository()
	planner := NewPlanner(repository, domain.Stock{"book": 5})
	request := domain.ManifestRequest{
		OrderID:  "order-2",
		Customer: "Ada",
		Packages: []domain.PackageRequest{{SKU: " BOOK ", Units: 1}},
	}

	if _, err := planner.Plan(context.Background(), request); err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if got := request.Packages[0].SKU; got != " BOOK " {
		t.Fatalf("Plan() changed the caller request package to %q", got)
	}

	first, ok := repository.Get("order-2")
	if !ok {
		t.Fatal("Get() did not return stored manifest")
	}
	first.Reservations[0].Units = 99
	second, ok := repository.Get("order-2")
	if !ok || second.Reservations[0].Units != 1 {
		t.Fatalf("Get() returned shared manifest data: %#v, %t", second, ok)
	}
}
