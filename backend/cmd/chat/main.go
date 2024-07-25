package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mikedoouglas/chat/internal/chat"
	"go.uber.org/zap"
)

func main() {
	environment := os.Getenv("ENV")
	if environment == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("error loading .env file")
		}
	}

	logger, err := zap.NewDevelopment()
	if environment == "production" {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatalf("fail to initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugaredLogger := logger.Sugar()
	handler := chat.NewHandler(sugaredLogger)
	http.HandleFunc("/ws", handler.HandleWebsocket)

	port := os.Getenv("PORT")
	sugaredLogger.Infof("Server started successfully on port :%s", port)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Error("Server startup failed", zap.Error(err))
	}
}
