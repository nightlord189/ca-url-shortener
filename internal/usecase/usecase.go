package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/nightlord189/ca-url-shortener/internal/config"
	"github.com/nightlord189/ca-url-shortener/internal/entity"
	"github.com/nightlord189/ca-url-shortener/internal/usecase/adapters"
	"github.com/nightlord189/ca-url-shortener/pkg"
	"github.com/nightlord189/ca-url-shortener/pkg/log"
)

type Usecase struct {
	Config  config.Config
	Storage adapters.IStorage
}

func New(cfg config.Config, storage adapters.IStorage) *Usecase {
	return &Usecase{
		Config:  cfg,
		Storage: storage,
	}
}

func (u *Usecase) Auth(ctx context.Context, username, password string) error {
	user, err := u.Storage.GetUser(ctx, username)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("get user error: %w", err)
	}

	passwordHash := pkg.GetHash(password)

	if user != nil {
		if user.PasswordHash != passwordHash {
			return ErrInvalidCredentials
		}
	} else {
		user = &entity.User{
			Username:     username,
			PasswordHash: passwordHash,
			Links:        nil,
		}
		if err := u.Storage.CreateUser(ctx, user); err != nil {
			return fmt.Errorf("create user error: %w", err)
		}
		log.Ctx(ctx).Infof("user %s created", username)
	}
	log.Ctx(ctx).Infof("user %s authorized", username)
	return nil
}

func (u *Usecase) PutLink(username, link string) (string, error) {
	panic("not implemented")
}

func (u *Usecase) GetOriginalLink(short string) (string, error) {
	panic("not implemented")
}
