package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/nightlord189/ca-url-shortener/docs"
	"github.com/nightlord189/ca-url-shortener/internal/config"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type Handler struct {
	Config  config.HTTPConfig
	Usecase IUsecase
}

func New(cfg config.HTTPConfig, uc IUsecase) *Handler {
	return &Handler{
		Config:  cfg,
		Usecase: uc,
	}
}

func (h *Handler) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/auth", h.Auth)

	r.Mount("/swagger", httpSwagger.WrapHandler)

	return http.ListenAndServe(fmt.Sprintf(":%d", h.Config.Port), r)
}
