package main

import (
	"fmt"

	"github.com/nightlord189/ca-url-shortener/internal/config"
	"github.com/nightlord189/ca-url-shortener/internal/delivery/http"
	"github.com/nightlord189/ca-url-shortener/internal/repo/mongo"
	"github.com/nightlord189/ca-url-shortener/internal/repo/redis"
	"github.com/nightlord189/ca-url-shortener/internal/usecase"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("start #1")

	cfg, err := config.LoadConfig("configs/config.yml")
	if err != nil {
		panic(err.Error())
	}

	var logger *zap.Logger

	if cfg.LogLevel == "debug" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	zap.ReplaceGlobals(logger)

	logger.Info("start #2")

	mongoRepo, err := mongo.New(cfg.Mongo)
	if err != nil {
		logger.Fatal("init mongo repo error: ", zap.Error(err))
	}

	redisRepo, err := redis.New(cfg.Redis)
	if err != nil {
		logger.Fatal("init redis repo error: ", zap.Error(err))
	}

	usecaseInst := usecase.New(mongoRepo, redisRepo)

	handler := http.New(cfg.HTTP, usecaseInst)

	logger.Info("running handler", zap.Int("port", cfg.HTTP.Port))

	if err := handler.Run(); err != nil {
		logger.Error("run handler error: ", zap.Error(err))
	}
}
