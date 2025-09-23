package utils

import (
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../")
	envPath := filepath.Join(basePath, ".env")

	if err := godotenv.Load(envPath); err != nil {
		panic("Failed to load .env file: " + err.Error())
	}
}