package adapters

type IRepo interface {
	PutLink(originalURL string) (string, error)
	GetOriginalLink(short string) (string, error)
}
