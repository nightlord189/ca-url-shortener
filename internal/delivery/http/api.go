package http

import (
	"errors"
	"github.com/nightlord189/ca-url-shortener/internal/usecase"
	"github.com/nightlord189/ca-url-shortener/pkg/log"
	"net/http"
)

// Auth godoc
// @Summary Request to issue access token
// @Tags auth
// @Accept  json
// @Produce json
// @Param data body AuthRequest true "Input model"
// @Success 200 {object} AuthResponse
// @Failure 401 {object} GenericResponse
// @Failure 400 {object} GenericResponse
// @Failure 500 {object} GenericResponse
// @Router /api/auth [Post]
// @BasePath /
func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req AuthRequest
	if err := parseBodyJSON(r, &req); err != nil {
		log.Ctx(ctx).Errorf("parse request error: %v", err.Error())
		responseJSON(ctx, w, http.StatusBadRequest, GenericError("invalid request"))
		return
	}

	if err := req.IsValid(); err != nil {
		responseJSON(ctx, w, http.StatusBadRequest, GenericError("invalid request: "+err.Error()))
		return
	}

	err := h.Usecase.Auth(ctx, req.Username, req.Password)
	switch {
	case errors.Is(err, usecase.ErrInvalidCredentials):
		responseJSON(ctx, w, http.StatusUnauthorized, GenericError(err.Error()))
		return
	case err != nil:
		responseJSON(ctx, w, http.StatusInternalServerError, GenericError("auth error: "+err.Error()))
		return
	}

	authToken, err := h.getToken(req.Username)
	if err != nil {
		responseJSON(ctx, w, http.StatusInternalServerError, GenericError("create token error: "+err.Error()))
		return
	}

	responseJSON(ctx, w, http.StatusOK, AuthResponse{AccessToken: authToken})
}
