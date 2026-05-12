// Package app — экспортируемая точка сборки приложения financial-analyst.
//
// Единственный пакет, который импортирует domain / infra / presentation
// одновременно. Используется из cmd/server и tests/integration.
package app

import (
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config — корневой конфиг приложения. Группируется по компонентам:
// каждый внешний источник или подсистема живёт в собственном sub-config.
//
// Источник значений — только переменные окружения; дефолтов в коде нет.
// Отсутствие любой `required` переменной — fatal на старте.
type Config struct {
	Moex   MoexConfig
	Server ServerConfig
}

// MoexConfig — параметры доступа к MOEX ISS REST API.
type MoexConfig struct {
	// BaseURL — корень MOEX ISS без завершающего слэша.
	// Пример: `https://iss.moex.com/iss`.
	BaseURL string `env:"MOEX_BASE_URL,required"`

	// Timeout — таймаут на один HTTP-запрос к MOEX ISS.
	// Пример: `10s`.
	Timeout time.Duration `env:"MOEX_TIMEOUT,required"`
}

// ServerConfig — параметры Connect-сервера financial-analyst.
type ServerConfig struct {
	// Port — TCP-порт, который слушает сервер. Bind по всем интерфейсам.
	Port uint16 `env:"SERVER_PORT,required"`

	// ReadHeaderTimeout — лимит на чтение заголовков HTTP-запроса.
	// Защищает от Slowloris-атак; см. http.Server.ReadHeaderTimeout.
	ReadHeaderTimeout time.Duration `env:"SERVER_READ_HEADER_TIMEOUT,required"`
}

// LoadConfig читает Config из переменных окружения. Перед разбором env
// подтягивает `./.env` (если файл существует) — это удобно для локальной
// разработки. В CI и production .env обычно нет; отсутствие файла — не ошибка.
func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("parse env config: %w", err)
	}
	return cfg, nil
}
