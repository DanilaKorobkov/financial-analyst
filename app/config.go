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
//
// Состав карточки компании конфигом не управляется — он вшит в
// исполняемый файл. В env живут только операционные параметры:
// адреса/токены провайдеров и параметры HTTP-сервера.
type Config struct {
	// FinanceMarker — параметры доступа к FinanceMarker REST API.
	FinanceMarker FinanceMarkerConfig
	// Moex — параметры доступа к MOEX ISS REST API.
	Moex MoexConfig
	// Server — параметры HTTP/Connect-сервера приложения.
	Server ServerConfig
}

// MoexConfig — параметры доступа к MOEX ISS REST API.
type MoexConfig struct {
	// BaseURL — корень MOEX ISS без завершающего слэша.
	// Пример: `https://iss.moex.com/iss`.
	BaseURL string `env:"MOEX_BASE_URL,required,notEmpty"`

	// Timeout — таймаут на один HTTP-запрос к MOEX ISS.
	// Пример: `10s`.
	Timeout time.Duration `env:"MOEX_TIMEOUT,required,notEmpty"`
}

// FinanceMarkerConfig — параметры доступа к FinanceMarker REST API.
type FinanceMarkerConfig struct {
	// BaseURL — корень FinanceMarker без завершающего слэша.
	// Пример: `https://financemarker.ru/api/fm/v2`.
	BaseURL string `env:"FINANCEMARKER_BASE_URL,required,notEmpty"`

	// Token — API-токен FinanceMarker. Передаётся query-параметром
	// `api_token` во всех запросах.
	Token string `env:"FINANCEMARKER_TOKEN,required,notEmpty"`

	// CacheRootDir — корневой каталог файлового кеша FinanceMarker.
	// Внутри провайдер раскладывает данные по подкаталогам
	// `<CacheRootDir>/<provider>/<bundle>`. Пример: `./.cache`.
	CacheRootDir string `env:"FINANCEMARKER_CACHE_ROOT_DIR,required,notEmpty"`

	// Timeout — таймаут на один HTTP-запрос.
	Timeout time.Duration `env:"FINANCEMARKER_TIMEOUT,required,notEmpty"`
}

// ServerConfig — параметры Connect-сервера financial-analyst.
type ServerConfig struct {
	// Port — TCP-порт, который слушает сервер. Bind по всем интерфейсам.
	Port uint16 `env:"SERVER_PORT,required,notEmpty"`

	// ReadHeaderTimeout — лимит на чтение заголовков HTTP-запроса.
	// Защищает от Slowloris-атак; см. http.Server.ReadHeaderTimeout.
	ReadHeaderTimeout time.Duration `env:"SERVER_READ_HEADER_TIMEOUT,required,notEmpty"`
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
