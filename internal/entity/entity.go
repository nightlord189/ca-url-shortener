package entity

const MinPasswordLength = 6

type User struct {
	Username     string
	PasswordHash string
	Links        []Link
}

type Link struct {
	Original string
	Short    string
}
