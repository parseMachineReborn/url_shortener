package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parseMachineReborn/url_shortener/internal/config"
	"github.com/parseMachineReborn/url_shortener/internal/handler/url"
	"github.com/parseMachineReborn/url_shortener/internal/handler/user"
	urlRepo "github.com/parseMachineReborn/url_shortener/internal/repository/postgres/url"
	userRepo "github.com/parseMachineReborn/url_shortener/internal/repository/postgres/user"
	urlSrv "github.com/parseMachineReborn/url_shortener/internal/service/url"
	userSrv "github.com/parseMachineReborn/url_shortener/internal/service/user"
)

func main() {
	config := config.NewConfig()

	pool := connectDB(config.DBConnectionString)
	defer pool.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	urlRepository := urlRepo.NewRepository(pool)
	userRepository := userRepo.NewRepository(pool)

	urlService := urlSrv.NewService(urlRepository)
	userService := userSrv.NewService(userRepository)

	urlHandler := url.NewHandler(urlService, config.SecretKey)
	userHandler := user.NewHandler(userService, config.SecretKey)

	mux := http.NewServeMux()
	urlHandler.RegisterRoutes(mux)
	userHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         config.Port,
		Handler:      mux,
		ReadTimeout:  config.ReadWriteTimeout,
		WriteTimeout: config.ReadWriteTimeout,
		IdleTimeout:  config.IdleTimeout,
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

	shutDownCtx, cancelFunc := context.WithTimeout(context.Background(), config.ShutDownPeriod)
	defer cancelFunc()

	err := server.Shutdown(shutDownCtx)
	if err != nil {
		log.Println("Произошла ошибка при мягком завершении.")
	}

	log.Println("Произошло мягкое завершение сервера")
}

func connectDB(connStr string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		log.Fatal("Ошибка при попытке создать пул подключений к БД")
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Нет ответа от БД")
	}

	fmt.Println("БД подключена")

	return pool
}
