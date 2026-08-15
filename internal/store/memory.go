package store

import (
	"context"
	"sync"
	"time"

	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/domain"
)

type MemoryRepository struct {
	mu        sync.Mutex
	manifests map[string]domain.Manifest
	active    bool
	delay     time.Duration
}

func NewMemoryRepository() *MemoryRepository {
	return NewMemoryRepositoryWithDelay(0)
}

func NewMemoryRepositoryWithDelay(delay time.Duration) *MemoryRepository {
	return &MemoryRepository{
		manifests: make(map[string]domain.Manifest),
		delay:     delay,
	}
}

func (repository *MemoryRepository) Begin(ctx context.Context) (Session, error) {
	if err := wait(ctx, repository.delay); err != nil {
		return nil, err
	}

	repository.mu.Lock()
	defer repository.mu.Unlock()
	if repository.active {
		return nil, ErrTransactionActive
	}
	repository.active = true
	return &memorySession{repository: repository}, nil
}

func (repository *MemoryRepository) Get(orderID string) (domain.Manifest, bool) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	manifest, ok := repository.manifests[orderID]
	return cloneManifest(manifest), ok
}

type memorySession struct {
	repository *MemoryRepository
	pending    *domain.Manifest
	committed  bool
}

func (session *memorySession) Save(ctx context.Context, manifest domain.Manifest) error {
	if err := wait(ctx, session.repository.delay); err != nil {
		return err
	}
	session.pending = pointerTo(cloneManifest(manifest))
	return nil
}

func (session *memorySession) Commit(ctx context.Context) error {
	if err := wait(ctx, session.repository.delay); err != nil {
		return err
	}
	if session.pending == nil {
		return nil
	}

	session.repository.mu.Lock()
	defer session.repository.mu.Unlock()
	if _, exists := session.repository.manifests[session.pending.OrderID]; exists {
		return ErrDuplicateOrder
	}
	session.repository.manifests[session.pending.OrderID] = cloneManifest(*session.pending)
	session.committed = true
	return nil
}

func (session *memorySession) Close() {
	session.repository.mu.Lock()
	defer session.repository.mu.Unlock()
	if session.committed && session.pending != nil {
		delete(session.repository.manifests, session.pending.OrderID)
	}
	session.repository.active = false
}

func wait(ctx context.Context, delay time.Duration) error {
	if delay == 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func pointerTo(manifest domain.Manifest) *domain.Manifest {
	return &manifest
}

func cloneManifest(manifest domain.Manifest) domain.Manifest {
	clone := manifest
	if manifest.DeliveryWindow != nil {
		clone.DeliveryWindow = &domain.DeliveryWindow{
			Start: manifest.DeliveryWindow.Start,
			End:   manifest.DeliveryWindow.End,
		}
	}
	clone.Reservations = append([]domain.Reservation(nil), manifest.Reservations...)
	return clone
}
