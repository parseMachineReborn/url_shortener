package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type config struct {
	DBConnectionString string
	Port               string
	ReadWriteTimeout   time.Duration
	IdleTimeout        time.Duration
	ShutDownPeriod     time.Duration
	SecretKey          string
}

func NewConfig() *config {
	if err := godotenv.Load(); err != nil {
		log.Println("Не было найдено .env файла")
		return &config{}
	}

	connStr := os.Getenv("DB_CONNECTION_STRING")

	port := os.Getenv("PORT")
	rwTimeout, err := strconv.Atoi(os.Getenv("READ_WRITE_TIMEOUT"))
	if err != nil {
		log.Fatal("Проблема инициализации, таймаут на чтение/запись отсутствует или не валиден")
	}
	rwTimeoutDuration := time.Duration(rwTimeout) * time.Second

	shutDownPeriod, err := strconv.Atoi(os.Getenv("SHUTDOWN_PERIOD"))
	if err != nil {
		log.Fatal("Проблема инициализации, период завершения приложения отсутствует или не валиден")
	}
	shutDownPeriodDuration := time.Duration(shutDownPeriod) * time.Second

	idleTimeout, err := strconv.Atoi(os.Getenv("IDLE_TIMEOUT"))
	if err != nil {
		log.Fatal("Проблема инициализации, таймаут на простой отсутствует или не валиден")
	}
	idleTimeoutDuration := time.Duration(idleTimeout) * time.Second

	secretKey := os.Getenv("JWT_SECRET_KEY")

	return &config{
		DBConnectionString: connStr,
		Port:               port,
		ReadWriteTimeout:   rwTimeoutDuration,
		IdleTimeout:        idleTimeoutDuration,
		ShutDownPeriod:     shutDownPeriodDuration,
		SecretKey:          secretKey,
	}
}
