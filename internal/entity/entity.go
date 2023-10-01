package entity

const MinPasswordLength = 6

type User struct {
	Username     string
	PasswordHash string
	Links        map[string]string
}
