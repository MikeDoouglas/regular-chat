package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/mikedoouglas/chat/internal/app/chat"
	"github.com/mikedoouglas/chat/internal/app/chat/config"
	namegenerator "github.com/mikedoouglas/chat/internal/pkg/generator"
	"go.uber.org/zap"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}
}

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := zap.NewDevelopment()
	if config.IsProd() {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatalf("fail to initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugaredLogger := logger.Sugar()
	room := chat.NewRoom()
	nameGenerator := namegenerator.New(sugaredLogger)
	handler := chat.NewHandler(room, nameGenerator, sugaredLogger)
	http.HandleFunc("/ws", handler.HandleWebsocket)

	sugaredLogger.Infof("starting server on port :%s", config.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil); err != nil {
		logger.Error("Server startup failed", zap.Error(err))
	}
}
