package store

import (
	"context"
	"testing"

	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/domain"
)

func TestMemoryRepositoryStoresIsolatedManifest(t *testing.T) {
	repository := NewMemoryRepository()
	session, err := repository.Begin(context.Background())
	if err != nil {
		t.Fatalf("Begin() error = %v", err)
	}
	defer session.Close()

	manifest := domain.Manifest{
		OrderID:      "order-1",
		Reservations: []domain.Reservation{{SKU: "book", Units: 1}},
	}
	if err := session.Save(context.Background(), manifest); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if err := session.Commit(context.Background()); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}

	manifest.Reservations[0].Units = 99
	stored, ok := repository.Get("order-1")
	if !ok || stored.Reservations[0].Units != 1 {
		t.Fatalf("Get() = %#v, %t", stored, ok)
	}
}

func TestMemoryRepositoryAllowsManifestWithoutDeliveryWindow(t *testing.T) {
	repository := NewMemoryRepository()
	session, err := repository.Begin(context.Background())
	if err != nil {
		t.Fatalf("Begin() error = %v", err)
	}
	defer session.Close()

	if err := session.Save(context.Background(), domain.Manifest{OrderID: "manual-order"}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if err := session.Commit(context.Background()); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}
	stored, ok := repository.Get("manual-order")
	if !ok || stored.DeliveryWindow != nil {
		t.Fatalf("Get() = %#v, %t", stored, ok)
	}
}
