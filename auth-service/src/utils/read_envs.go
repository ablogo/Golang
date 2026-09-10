package utils

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Settings struct {
	DB_HOST     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
	DB_PORT     string
	SERVER_PORT string
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

func getEnvVariableNumber(name string, defaultValue int) string {
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return strconv.Itoa(defaultValue)
	}
	return strconv.Itoa(value)
}

func (s *Settings) GetValues() {
	godotenv.Load("../.env")

	s.SERVER_PORT = getEnvVariableNumber("SERVER_PORT", 8000)
	s.DB_HOST = getEnvVariable("DB_HOST", "localhost")
	s.DB_NAME = getEnvVariable("DB_NAME", "auth-service")
	s.DB_USER = getEnvVariable("DB_USER", "")
	s.DB_PASSWORD = getEnvVariable("DB_PASSWORD", "")
	s.DB_PORT = getEnvVariableNumber("DB_PORT", 5432)
}
