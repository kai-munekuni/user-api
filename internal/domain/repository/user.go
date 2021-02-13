package repository

import (
	"context"

	"github.com/kai-munekuni/user-api/internal/domain/model"
)

// User user_repository
type User interface {
	CreateUser(ctx context.Context, id, password string) error
	Authorize(ctx context.Context, id, password string) error
	Find(ctx context.Context, ID string) (*model.User, error)
	UpdateField(ctx context.Context, ID string, nickname, comment *string) (*model.User, error)
	Delete(ctx context.Context, ID string) error
}
