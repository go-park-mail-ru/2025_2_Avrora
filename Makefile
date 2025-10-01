.PHONY: test up down

# Поднять только тестовую БД
up:
	docker-compose up -d postgres_test

# Прогон тестов (локально, используя тестовую БД на 5433)
test: up
	# ждём пока база прогрузится
	sleep 5
	# запускаем go test с нужным env
	AVRORA_DB_HOST=localhost \
	AVRORA_DB_PORT=5433 \
	AVRORA_DB_NAME=avrora_test \
	AVRORA_DB_USER=avrora_user 
	AVRORA_DB_PASSWORD=avrora_password \
		go test ./... -v
	$(MAKE) down

# Остановить контейнеры
down:
	docker-compose down
