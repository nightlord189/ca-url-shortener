package http

import (
	"errors"
	"github.com/nightlord189/ca-url-shortener/internal/usecase"
	"github.com/nightlord189/ca-url-shortener/pkg/log"
	"net/http"
	"strings"
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

	authToken, err := h.getJWTToken(req.Username)
	if err != nil {
		responseJSON(ctx, w, http.StatusInternalServerError, GenericError("create token error: "+err.Error()))
		return
	}

	responseJSON(ctx, w, http.StatusOK, AuthResponse{AccessToken: authToken})
}

// PutLink godoc
// @Summary Create new short link
// @Tags link
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param data body PutLinkRequest true "Input model"
// @Success 200 {object} PutLinkResponse
// @Failure 401 {object} GenericResponse
// @Failure 400 {object} GenericResponse
// @Failure 500 {object} GenericResponse
// @Router /api/link [Put]
// @BasePath /
func (h *Handler) PutLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req PutLinkRequest
	if err := parseBodyJSON(r, &req); err != nil {
		log.Ctx(ctx).Errorf("parse request error: %v", err.Error())
		responseJSON(ctx, w, http.StatusBadRequest, GenericError("invalid request"))
		return
	}

	if err := req.IsValid(); err != nil {
		responseJSON(ctx, w, http.StatusBadRequest, GenericError("invalid request: "+err.Error()))
		return
	}

	username := ctx.Value("username").(string)

	shortLink, err := h.Usecase.PutLink(ctx, username, req.OriginalURL)
	if err != nil {
		responseJSON(ctx, w, http.StatusInternalServerError, GenericError(err.Error()))
		return
	}

	responseJSON(ctx, w, http.StatusOK, PutLinkResponse{ShortURL: shortLink})
}

// GetLink godoc
// @Summary Go to original url
// @Tags link
// @Accept  json
// @Produce json
// @Param short path string true "short relative url"
// @Success 200
// @Failure 404
// @Failure 500
// @Router /{short} [Get]
// @BasePath /
func (h *Handler) GetLink(w http.ResponseWriter, r *http.Request) {
	relativeURL := strings.Replace(r.URL.RequestURI(), "/", "", 1)

	originalURL, err := h.Usecase.GetOriginalLink(r.Context(), relativeURL)
	if err != nil {
		responseString(r.Context(), w, http.StatusNotFound, "not found")
		return
	}

	log.Ctx(r.Context()).Debugf("redirecting %s to %s", relativeURL, originalURL)

	http.Redirect(w, r, originalURL, http.StatusFound)
}
