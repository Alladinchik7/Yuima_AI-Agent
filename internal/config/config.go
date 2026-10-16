package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AIAPIKey      string
	AIBaseURL     string
	AIModel       string

	WeatherAPIKey string

	SearxNGURL   string
	FormatSearxNG string
	LangSerxNG    string
	Safesearch   string
	
	UseMockAgent  bool
	UseMockSearch bool

	LogLevel      string
	LogFormat     string

	UseCloud	  bool

	UseMemory     bool
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		AIAPIKey:      getEnv("API_AI", ""),
		AIBaseURL:     getEnv("BASE_URL", ""),
		AIModel:       getEnv("AI_MODEL", ""),
		
		WeatherAPIKey: getEnv("WEATHER_API_KEY", ""),

		SearxNGURL: getEnv("SEARXNG_URL", "http://localhost:8888"),
		FormatSearxNG: getEnv("FORMAT_SEARXNG", "text"),
		LangSerxNG: getEnv("LANG_SEARXNG", "ru"),
		Safesearch: getEnv("SAFESEARXNG", "0"),
		
		UseMockAgent:  getEnv("USE_MOCK_AGENT", "false") == "true",
		UseMockSearch: getEnv("USE_MOCK_SEARCH", "false") == "true",

		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "text"),

		UseCloud: getEnv("USE_CLOUD", "false") == "true",

		UseMemory: getEnv("USE_MEMORY", "false") == "true",
	}
}

func getEnv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}