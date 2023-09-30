package http

type IUsecase interface {
	Auth(username, password string) (string, error)
	PutLink(originalURL, username string) (string, error)
	GetOriginalLink(short string) (string, error)
}
