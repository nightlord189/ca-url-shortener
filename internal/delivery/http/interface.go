package http

import "context"

type IUsecase interface {
	Auth(ctx context.Context, username, password string) error
	PutLink(ctx context.Context, originalURL, username string) (string, error)
	GetOriginalLink(short string) (string, error)
}
