package store

import (
	"context"
	"errors"

	"github.com/1260124186-cc/fulfillment-manifest-cli/internal/domain"
)

var (
	ErrDuplicateOrder    = errors.New("manifest already exists for order")
	ErrTransactionActive = errors.New("another manifest session is active")
)

type Repository interface {
	Begin(context.Context) (Session, error)
	Get(string) (domain.Manifest, bool)
}

type Session interface {
	Save(context.Context, domain.Manifest) error
	Commit(context.Context) error
	Close()
}
