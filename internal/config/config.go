package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config armazena as configurações da aplicação.
type Config struct {
	Port                 int
	APIKey               string
	DefaultVoice         string
	DefaultFormat        string
	DefaultSpeed         float64
	DefaultLanguage      string
	RequireAPIKey        bool
	RemoveFilter         bool
	DetailedErrorLogging bool
	Proxy                string
	CacheEnabled         bool
	CacheMaxMB           int
	CacheTTL             time.Duration
}

// Load carrega as variáveis de ambiente com fallback para os valores padrão.
func Load() *Config {
	_ = godotenv.Load() // Carrega .env caso exista

	cacheTTLHours := getEnvInt("CACHE_TTL_HOURS", 24)

	return &Config{
		Port:                 getEnvInt("PORT", 5050),
		APIKey:               os.Getenv("API_KEY"),
		DefaultVoice:         getEnvString("DEFAULT_VOICE", "pt-BR-ThalitaMultilingualNeural"),
		DefaultFormat:        getEnvString("DEFAULT_RESPONSE_FORMAT", "mp3"),
		DefaultSpeed:         getEnvFloat("DEFAULT_SPEED", 1.0),
		DefaultLanguage:      getEnvString("DEFAULT_LANGUAGE", "pt-BR"),
		RequireAPIKey:        getEnvBool("REQUIRE_API_KEY", true),
		RemoveFilter:         getEnvBool("REMOVE_FILTER", false),
		DetailedErrorLogging: getEnvBool("DETAILED_ERROR_LOGGING", true),
		Proxy:                os.Getenv("PROXY"),
		CacheEnabled:         getEnvBool("CACHE_ENABLED", true),
		CacheMaxMB:           getEnvInt("CACHE_MAX_MB", 50),
		CacheTTL:             time.Duration(cacheTTLHours) * time.Hour,
	}
}

func getEnvString(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func getEnvFloat(key string, defaultVal float64) float64 {
	if valStr := os.Getenv(key); valStr != "" {
		if val, err := strconv.ParseFloat(valStr, 64); err == nil {
			return val
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	lower := strings.ToLower(valStr)
	return lower == "true" || lower == "1" || lower == "yes" || lower == "t" || lower == "y"
}
