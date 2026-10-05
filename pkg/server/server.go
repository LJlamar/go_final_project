package server

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/LJlamar/go_final_project/pkg/api"
)

type Application struct {
	Server *http.Server
	Logger *log.Logger
}

func NewServer(logger *log.Logger) *Application {

	//HTTP-router:
	mux := http.NewServeMux()
	api.Init(mux)

	//Port:
	port := os.Getenv("PORT")
	if port == "" {
		port = "7540"
	}
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 || len(strconv.Itoa(portNum)) != len(port) {
		logger.Printf("Некорректный порт '%s'. Используем значение по умолчанию: 7540", port)
		port = "7540"
	}

	addr := ":" + port

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	app := Application{
		Server: srv,
		Logger: logger,
	}

	return &app
}
