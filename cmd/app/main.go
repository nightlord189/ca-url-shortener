package main

import (
	"fmt"
	"github.com/nightlord189/ca-url-shortener/internal/config"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("start #1")

	cfg, err := config.LoadConfig("configs/config.yml")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(cfg)

	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)

	logger.Info("start #2")
}
