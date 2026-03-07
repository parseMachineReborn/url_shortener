package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/parseMachineReborn/url_shortener/internal/handler"
	"github.com/parseMachineReborn/url_shortener/internal/repository"
	"github.com/parseMachineReborn/url_shortener/internal/service"
)

const shutDownPeriod = 15 * time.Second

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println("Ошибка при запуске сервера")
				panic(err)
			}
		}
	}()

	<-ctx.Done()

	shutDownCtx, cancelFunc := context.WithTimeout(context.Background(), shutDownPeriod)
	defer cancelFunc()

	err := server.Shutdown(shutDownCtx)
	if err != nil {
		log.Println("Произошла ошибка при мягком завершении.")
	}

	log.Println("Произошло мягкое завершение сервера")
}
