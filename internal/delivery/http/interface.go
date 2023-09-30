package http

import "context"

type IUsecase interface {
	Auth(ctx context.Context, username, password string) error
	PutLink(originalURL, username string) (string, error)
	GetOriginalLink(short string) (string, error)
}
