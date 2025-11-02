¡Write-Host "🚀 Поднимаем тестовую БД..."
docker-compose up -d postgres_test

# Ждём готовности контейнера
Write-Host "⏳ Ждём 10 секунд, пока БД стартует..."
Start-Sleep -Seconds 10

Write-Host "✅ Устанавливаем переменные окружения для тестов..."
$env:DB_HOST = "localhost"
$env:DB_PORT = "5433"
$env:DB_NAME = "avrora_test"
$env:DB_USER = "avrora_user"
$env:DB_PASSWORD = "avrora_password"
$env:JWT_SECRET = "super_secret_shhhhhh"
$env:PASSWORD_PEPPER = "bmstu_my!<>_super_secret_pepper_here_32_chars_min"
$env:CORS_ORIGIN = "http://localhost:3000"

Write-Host "✅ Запускаем тесты с покрытием..."
go test ./... -v -cover -coverprofile=coverage.out

Write-Host "📊 Выводим суммарное покрытие..."
go tool cover -func .\coverage.out

# Можно открыть HTML-отчёт (по желанию)
# go tool cover -html=coverage.out

Write-Host "🧹 Останавливаем контейнеры..."
docker-compose down
