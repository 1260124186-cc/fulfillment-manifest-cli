package service

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestPlannerStopsWhenStorageContextExpires(t *testing.T) {
	repository := store.NewMemoryRepositoryWithDelay(50 * time.Millisecond)
	planner := NewPlanner(repository, domain.Stock{"book": 5})
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	_, err := planner.Plan(ctx, domain.ManifestRequest{
		OrderID:  "order-4",
		Customer: "Ada",
		Packages: []domain.PackageRequest{{SKU: "book", Units: 1}},
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Plan() error = %v, want deadline exceeded", err)
	}
	if _, ok := repository.Get("order-4"); ok {
		t.Fatal("Plan() stored a manifest after the context deadline")
	}
}

func TestPlannerDoesNotCommitAfterSaveCancelsContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	repository := &cancelOnSaveRepository{cancel: cancel}
	planner := NewPlanner(repository, domain.Stock{"book": 5})

	_, err := planner.Plan(ctx, domain.ManifestRequest{
		OrderID:  "order-5",
		Customer: "Ada",
		Packages: []domain.PackageRequest{{SKU: "book", Units: 1}},
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Plan() error = %v, want context canceled", err)
	}
	if repository.committed {
		t.Fatal("Plan() committed after Save canceled the context")
	}
}

type cancelOnSaveRepository struct {
	cancel    context.CancelFunc
	committed bool
}

func (repository *cancelOnSaveRepository) Begin(context.Context) (store.Session, error) {
	return cancelOnSaveSession{repository: repository}, nil
}

func (repository *cancelOnSaveRepository) Get(string) (domain.Manifest, bool) {
	return domain.Manifest{}, false
}

type cancelOnSaveSession struct {
	repository *cancelOnSaveRepository
}

func (session cancelOnSaveSession) Save(ctx context.Context, _ domain.Manifest) error {
	session.repository.cancel()
	return ctx.Err()
}

func (session cancelOnSaveSession) Commit(context.Context) error {
	session.repository.committed = true
	return nil
}

func (cancelOnSaveSession) Close() {}
