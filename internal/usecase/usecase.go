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
	"go.uber.org/zap"
)

type Usecase struct {
	Config  config.Config
	Storage adapters.IStorage
	Cache   adapters.ICache
}

func New(cfg config.Config, storage adapters.IStorage, cache adapters.ICache) *Usecase {
	return &Usecase{
		Config:  cfg,
		Storage: storage,
		Cache:   cache,
	}
}

func (u *Usecase) Auth(ctx context.Context, username, password string) error {
	user, err := u.Storage.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("get user error: %w", err)
	}

	passwordHash := pkg.GetSHA256Hash(password)

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

func (u *Usecase) PutLink(ctx context.Context, username, originalURL string) (string, error) {
	user, err := u.Storage.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("get user error: %w", err)
	}

	if user.Links == nil {
		user.Links = make(map[string]string, 1)
	}

	shortURL := pkg.GetFNVHash(username + originalURL)
	user.Links[shortURL] = originalURL

	if err := u.Storage.UpdateUserLinks(ctx, user); err != nil {
		return "", fmt.Errorf("save user error: %w", err)
	}

	if err := u.Cache.PutLink(ctx, shortURL, originalURL); err != nil {
		log.Ctx(ctx).Error("put new link to cache error", zap.Error(err))
	}

	return shortURL, nil
}

func (u *Usecase) GetOriginalLink(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := u.Cache.GetLink(ctx, shortURL)
	if err == nil {
		return originalURL, nil
	}

	originalURL, err = u.Storage.GetLink(ctx, shortURL)
	if err != nil {
		return "", fmt.Errorf("get originalURL error: %w", err)
	}

	if originalURL == "" {
		return "", fmt.Errorf("originalURL not found")
	}

	if err := u.Cache.PutLink(ctx, shortURL, originalURL); err != nil {
		log.Ctx(ctx).Error("put new originalURL to cache error", zap.Error(err))
	}

	return originalURL, nil
}
