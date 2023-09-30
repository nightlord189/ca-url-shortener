package adapters

import (
	"context"
	"github.com/nightlord189/ca-url-shortener/internal/entity"
)

type IStorage interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUser(ctx context.Context, username string) (*entity.User, error)
	PutLink(username, link string) error
}
