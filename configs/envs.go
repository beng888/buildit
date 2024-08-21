package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost              string
	Port                    string
	DBUser                  string
	DBPassword              string
	DBAddress               string
	DBName                  string
	JWTSecret               string
	JWTExpirationInSeconds  int
	CookiesAuthSecret       string
	CookiesAuthAgeInSeconds int
	CookiesAuthIsSecure     bool
	CookiesAuthIsHttpOnly   bool
	DiscordClientID         string
	DiscordClientSecret     string
	GithubClientID          string
	GithubClientSecret      string
	GoogleClientID          string
	GoogleClientSecret      string
}

const (
	twoDaysInSeconds   = 60 * 60 * 24 * 2
	sevenDaysInSeconds = 60 * 60 * 24 * 7
)

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:              getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                    getEnv("PORT", "8080"),
		DBUser:                  getEnv("DB_USER", "root"),
		DBPassword:              getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:               fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                  getEnv("DB_NAME", "ecom"),
		JWTSecret:               getEnv("JWT_SECRET", "not-so-secret-now-is-it?"),
		JWTExpirationInSeconds:  getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", sevenDaysInSeconds),
		CookiesAuthSecret:       getEnv("COOKIES_AUTH_SECRET", "some-very-secret-key"),
		CookiesAuthAgeInSeconds: getEnvAsInt("COOKIES_AUTH_AGE_IN_SECONDS", twoDaysInSeconds),
		CookiesAuthIsSecure:     getEnvAsBool("COOKIES_AUTH_IS_SECURE", false),
		CookiesAuthIsHttpOnly:   getEnvAsBool("COOKIES_AUTH_IS_HTTP_ONLY", false),
		DiscordClientID:         getEnvOrError("DISCORD_CLIENT_ID"),
		DiscordClientSecret:     getEnvOrError("DISCORD_CLIENT_SECRET"),
		GithubClientID:          getEnvOrError("GITHUB_CLIENT_ID"),
		GithubClientSecret:      getEnvOrError("GITHUB_CLIENT_SECRET"),
		GoogleClientID:          getEnvOrError("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:      getEnvOrError("GOOGLE_CLIENT_SECRET"),
	}
}

// Gets the env by key or fallbacks
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvOrError(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	panic(fmt.Sprintf("Environment variable %s is not set", key))

}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}

		return b
	}

	return fallback
}
