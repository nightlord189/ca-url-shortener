package usecase

import "github.com/nightlord189/ca-url-shortener/internal/config"

type Usecase struct {
	Config config.Config
}

func New(cfg config.Config) *Usecase {
	return &Usecase{Config: cfg}
}

func (u *Usecase) Auth(username, password string) (string, error) {
	panic("not implemented")
}

func (u *Usecase) PutLink(originalURL, username string) (string, error) {
	panic("not implemented")
}

func (u *Usecase) GetOriginalLink(short string) (string, error) {
	panic("not implemented")
}
