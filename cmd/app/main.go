package main

import (
	"fmt"
	"github.com/nightlord189/ca-url-shortener/internal/config"
	"github.com/nightlord189/ca-url-shortener/internal/delivery/http"
	"github.com/nightlord189/ca-url-shortener/internal/usecase"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("start #1")

	cfg, err := config.LoadConfig("configs/config.yml")
	if err != nil {
		panic(err.Error())
	}

	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)

	logger.Info("start #2")

	usecaseInst := usecase.New(cfg)

	handler := http.New(cfg.HTTP, usecaseInst)

	logger.Info("running handler", zap.Int("port", cfg.HTTP.Port))

	if err := handler.Run(); err != nil {
		logger.Error("run handler error: ", zap.Error(err))
	}
}
