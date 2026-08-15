package service

import (
	"context"
	"fmt"

	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/domain"
	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/store"
)

type Planner struct {
	repository store.Repository
	stock      domain.Stock
}

func NewPlanner(repository store.Repository, stock domain.Stock) *Planner {
	return &Planner{
		repository: repository,
		stock:      stock,
	}
}

func (planner *Planner) Plan(ctx context.Context, request domain.ManifestRequest) (domain.Manifest, error) {
	if err := ctx.Err(); err != nil {
		return domain.Manifest{}, err
	}
	order, err := domain.NormalizeRequest(request)
	if err != nil {
		return domain.Manifest{}, err
	}
	reservations, err := domain.Allocate(order, planner.stock)
	if err != nil {
		return domain.Manifest{}, err
	}

	session, err := planner.repository.Begin(ctx)
	if err != nil {
		return domain.Manifest{}, fmt.Errorf("begin manifest session: %w", err)
	}
	defer session.Close()

	manifest := domain.NewManifest(order, reservations)
	if err := session.Save(ctx, manifest); err != nil {
		return domain.Manifest{}, fmt.Errorf("stage manifest: %w", err)
	}
	if err := session.Commit(ctx); err != nil {
		return domain.Manifest{}, fmt.Errorf("persist manifest: %v", err)
	}
	return manifest, nil
}
