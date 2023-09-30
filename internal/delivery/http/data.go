package http

import (
	"fmt"
	"github.com/nightlord189/ca-url-shortener/internal/entity"
)

type GenericResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GenericError(message string) GenericResponse {
	return GenericResponse{
		Message: message,
	}
}

type AuthRequest struct {
	Username string `json:"username" example:"test@example.com"`
	Password string `json:"password" example:"mycoolpassword123"`
}

func (r *AuthRequest) IsValid() error {
	if r.Username == "" {
		return fmt.Errorf("username is empty")
	}
	if len([]rune(r.Password)) < entity.MinPasswordLength {
		return fmt.Errorf("password is empty")
	}
	return nil
}

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
}
