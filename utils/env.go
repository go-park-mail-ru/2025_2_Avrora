// utils/env.go
package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// RequiredEnvVars — список обязательных переменных окружения
var RequiredEnvVars = []string{
	"DB_HOST",
	"DB_PORT",
	"DB_USER",
	"DB_PASSWORD",
	"DB_NAME",
	"JWT_SECRET",
}

func LoadEnv() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Не удалось определить текущую директорию:", err)
	}

	projectRoot := findProjectRoot(wd)
	if projectRoot == "" {
		log.Fatal("Не найден корень проекта (go.mod). Убедитесь, что запускаете из проекта.")
	}

	envPath := filepath.Join(projectRoot, ".env")

	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			log.Fatal("Ошибка загрузки .env:", err)
		}
		log.Println("✅ Переменные окружения загружены из .env")
	} else {
		examplePath := filepath.Join(projectRoot, ".env.example")
		if _, err := os.Stat(examplePath); err == nil {
			log.Println("⚠️  Файл .env не найден. Скопируйте .env.example в .env и заполните значения.")
		} else {
			log.Println("⚠️  Файл .env не найден и .env.example отсутствует.")
		}
	}

	validateEnvVars()
}

func findProjectRoot(startDir string) string {
	dir := startDir
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func validateEnvVars() {
	missing := []string{}
	for _, key := range RequiredEnvVars {
		if val := os.Getenv(key); val == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		log.Fatalf("❌ Отсутствуют обязательные переменные окружения: %v", missing)
	}

	log.Println("✅ Все обязательные переменные окружения заданы")
}

func GetPostgresDSN() string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )
}

func GetTestPostgresDSN() string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME_TEST"),
    )
}