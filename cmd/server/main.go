package main

import (
	"log"
	"net/http"
	"time"

	"github.com/parseMachineReborn/url_shortener/internal/handler"
	"github.com/parseMachineReborn/url_shortener/internal/repository"
	"github.com/parseMachineReborn/url_shortener/internal/service"
)

func main() {
	repository := repository.NewDefaultRepository()
	service := service.NewURLService(repository)
	handler := handler.NewHandler(service)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
