package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/nightlord189/ca-url-shortener/pkg/log"
	"io"
	"net/http"
	"time"
)

func (h *Handler) getToken(username string) (string, error) {
	payload := jwt.MapClaims{
		"username": username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Second * time.Duration(h.Config.AuthTokenExpTime)).Unix(),
		"iss":      "ca-url-shortener",
	}

	return createToken(payload, h.Config.AuthSecret)
}

func createToken(payload jwt.MapClaims, secret string) (string, error) {
	var err error
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func responseJSON(ctx context.Context, w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	marshaled, err := json.Marshal(data)
	if err != nil {
		log.Ctx(ctx).Errorf("marshal response error: %v", err.Error())
	}
	if _, err = w.Write(marshaled); err != nil {
		log.Ctx(ctx).Errorf("write response error: %v", err.Error())
	}
}

func parseBodyJSON(r *http.Request, out interface{}) error {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read request body error: %w", err)
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Ctx(r.Context()).Errorf("close request body error: %v", err.Error())
		}
	}()

	if err := json.Unmarshal(rawBody, out); err != nil {
		return fmt.Errorf("unmarshal json error: %w", err)
	}

	return nil
}
