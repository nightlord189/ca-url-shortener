package adapters

import (
	"context"
	"github.com/nightlord189/ca-url-shortener/internal/entity"
)

type IStorage interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUser(ctx context.Context, username string) (*entity.User, error)
	UpdateUserLinks(ctx context.Context, user *entity.User) error
}

type ICache interface {
	PutLink(ctx context.Context, shortURL, originalURL string) error
	GetLink(ctx context.Context, shortURL string) (string, error)
}
