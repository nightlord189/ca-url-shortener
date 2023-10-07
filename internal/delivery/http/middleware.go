package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type ContextKey string

const ContextUsername ContextKey = "username"

func (h *Handler) AuthMdw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			responseJSON(r.Context(), w, http.StatusUnauthorized, GenericError("invalid auth header"))
			return
		}
		authSplitted := strings.Split(authHeader, " ")
		if len(authSplitted) != 2 {
			responseJSON(r.Context(), w, http.StatusUnauthorized, GenericError("invalid auth header"))
			return
		}

		token := authSplitted[1]

		parsedToken, err := h.isJWTTokenValid(token)
		if err != nil {
			responseJSON(r.Context(), w, http.StatusUnauthorized, GenericError(err.Error()))
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			responseJSON(r.Context(), w, http.StatusUnauthorized, GenericError("invalid jwt claims"))
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			responseJSON(r.Context(), w, http.StatusUnauthorized, GenericError("invalid username claim"))
			return
		}

		ctx := context.WithValue(r.Context(), ContextUsername, username)

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
