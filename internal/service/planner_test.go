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

func TestPlannerReleasesSessionForTheNextOrder(t *testing.T) {
	repository := store.NewMemoryRepository()
	planner := NewPlanner(repository, domain.Stock{"book": 5})

	for _, orderID := range []string{"order-5a", "order-5b"} {
		_, err := planner.Plan(context.Background(), domain.ManifestRequest{
			OrderID:  orderID,
			Customer: "Ada",
			Packages: []domain.PackageRequest{{SKU: "book", Units: 1}},
		})
		if err != nil {
			t.Fatalf("Plan(%q) error = %v", orderID, err)
		}
	}
}
