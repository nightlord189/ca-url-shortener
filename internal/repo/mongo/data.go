package mongo

import "github.com/nightlord189/ca-url-shortener/internal/entity"

type User struct {
	Username     string `bson:"username"`
	PasswordHash string `bson:"password_hash"`
	Links        []Link `bson:"links,omitempty"`
}

func (u *User) ToEntity() entity.User {
	result := entity.User{
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		Links:        nil,
	}
	if u.Links != nil {
		result.Links = make([]entity.Link, len(u.Links))
		for i := range u.Links {
			result.Links[i] = u.Links[i].ToEntity()
		}
	}
	return result
}

func UserFromEntity(item *entity.User) User {
	result := User{
		Username:     item.Username,
		PasswordHash: item.PasswordHash,
		Links:        nil,
	}
	if item.Links != nil {
		result.Links = make([]Link, len(item.Links))
		for i := range item.Links {
			result.Links[i] = LinkFromEntity(&item.Links[i])
		}
	}
	return result
}

type Link struct {
	OriginalURL string `bson:"original_url"`
	ShortURL    string `bson:"short_url"`
}

func (l *Link) ToEntity() entity.Link {
	return entity.Link{
		OriginalURL: l.OriginalURL,
		ShortURL:    l.ShortURL,
	}
}

func LinkFromEntity(item *entity.Link) Link {
	return Link{
		OriginalURL: item.OriginalURL,
		ShortURL:    item.ShortURL,
	}
}
