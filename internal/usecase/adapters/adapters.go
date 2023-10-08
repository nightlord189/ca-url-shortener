package adapters

import (
	"context"

	"github.com/nightlord189/ca-url-shortener/internal/entity"
)

//go:generate mockgen -destination=../mock/istorage.go -package=mock github.com/nightlord189/ca-url-shortener/internal/usecase/adapters IStorage
type IStorage interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetLink(ctx context.Context, shortURL string) (string, error)
	UpdateUserLinks(ctx context.Context, user *entity.User) error
}

//go:generate mockgen -destination=../mock/icache.go -package=mock github.com/nightlord189/ca-url-shortener/internal/usecase/adapters ICache
type ICache interface {
	PutLink(ctx context.Context, shortURL, originalURL string) error
	GetLink(ctx context.Context, shortURL string) (string, error)
}
