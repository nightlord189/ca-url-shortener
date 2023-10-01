package mongo

import (
	"github.com/nightlord189/ca-url-shortener/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username"`
	PasswordHash string             `bson:"password_hash"`
	Links        map[string]string  `bson:"links,omitempty"`
}

func (u *User) ToEntity() entity.User {
	result := entity.User{
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		Links:        u.Links,
	}
	return result
}

func UserFromEntity(item *entity.User) User {
	result := User{
		Username:     item.Username,
		PasswordHash: item.PasswordHash,
		Links:        item.Links,
	}
	return result
}
