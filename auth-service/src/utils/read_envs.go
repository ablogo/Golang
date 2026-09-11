package utils

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Settings struct {
	DB_HOST              string
	DB_NAME              string
	DB_USER              string
	DB_PASSWORD          string
	DB_PORT              string
	SERVER_PORT          string
	SERVER_READ_TIMEOUT  int
	SERVER_WRITE_TIMEOUT int
	SERVER_IDLE_TIMEOUT  int
	SERVER_MAX_HEADER    int
	JWT_SECRET_TOKEN     string
	JWT_ALGORITHM        string
	JWT_EXPIRE_MINUTES   string
	CORS_ALLOWED_HOSTS   string
}

var instance *Settings

// sync.once: Once is an object that will perform exactly an action only once.
var singleton sync.Once

func (s *Settings) GetInstance() *Settings {
	singleton.Do(func() {
		s.GetValues()
		instance = s
	})
	return instance
}

func getEnvVariable(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func getNumberEnvVariableAsString(name string, defaultValue int) string {
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return strconv.Itoa(defaultValue)
	}
	return strconv.Itoa(value)
}

func getNumberEnvVariable(name string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return defaultValue
	}
	return value
}

func (s *Settings) GetValues() {
	godotenv.Load("../.env")

	s.SERVER_PORT = getNumberEnvVariableAsString("SERVER_PORT", 8000)
	s.SERVER_READ_TIMEOUT = getNumberEnvVariable("SERVER_READ_TIMEOUT", 60)
	s.SERVER_WRITE_TIMEOUT = getNumberEnvVariable("SERVER_WRITE_TIMEOUT", 60)
	s.SERVER_IDLE_TIMEOUT = getNumberEnvVariable("SERVER_IDLE_TIMEOUT", 120)
	s.DB_HOST = getEnvVariable("DB_HOST", "localhost")
	s.DB_NAME = getEnvVariable("DB_NAME", "auth-service")
	s.DB_USER = getEnvVariable("DB_USER", "")
	s.DB_PASSWORD = getEnvVariable("DB_PASSWORD", "")
	s.DB_PORT = getNumberEnvVariableAsString("DB_PORT", 5432)
}
