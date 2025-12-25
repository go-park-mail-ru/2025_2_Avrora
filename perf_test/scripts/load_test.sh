#!/bin/bash

set -e

# Конфигурация
TARGET_URL="${TARGET_URL:-http://localhost:8080}"
DURATION="${DURATION:-30s}"
RATE="${RATE:-100}"
OUTPUT_DIR="perf_test/results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULT_DIR="$OUTPUT_DIR/$TIMESTAMP"

# Цвета
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}╔════════════════════════════════════╗${NC}"
echo -e "${GREEN}║   Load Testing Script              ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════╝${NC}"
echo ""
echo -e "Target:   ${YELLOW}$TARGET_URL${NC}"
echo -e "Duration: ${YELLOW}$DURATION${NC}"
echo -e "Rate:     ${YELLOW}$RATE req/s${NC}"
echo ""

# Создаем директорию для результатов
mkdir -p "$RESULT_DIR"

# Функция для получения токена
get_token() {
    echo -e "${YELLOW}🔐 Getting JWT token...${NC}"
    
    TOKEN=$(curl -s -X POST "$TARGET_URL/api/v1/login" \
      -H "Content-Type: application/json" \
      -d '{"email":"user2@example.com","password":"Postgres1!"}' \
      | jq -r '.token')
    
    if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
      echo -e "${RED}❌ Failed to get token!${NC}"
      exit 1
    fi
    
    echo -e "${GREEN}✅ Token obtained${NC}"
    echo "$TOKEN" > "$RESULT_DIR/token.txt"
}

# Функция для сброса статистики БД
reset_db_stats() {
    echo -e "${YELLOW}📊 Resetting database statistics...${NC}"
    psql "$DATABASE_URL" -c "SELECT pg_stat_statements_reset();" > /dev/null 2>&1 || true
}

# Функция для сохранения статистики БД
save_db_stats() {
    local filename=$1
    echo -e "${YELLOW}💾 Saving database statistics to $filename...${NC}"
    psql "$DATABASE_URL" -f perf_test/sql/check_slow_queries.sql > "$RESULT_DIR/$filename"
}

# Функция запуска теста
run_test() {
    local test_name=$1
    local targets_file=$2
    
    echo -e "${YELLOW}🚀 Running test: $test_name${NC}"
    
    vegeta attack \
      -duration=$DURATION \
      -rate=$RATE \
      -targets="$targets_file" \
      > "$RESULT_DIR/${test_name}_results.bin"
    
    # Генерируем отчеты
    vegeta report "$RESULT_DIR/${test_name}_results.bin" > "$RESULT_DIR/${test_name}_report.txt"
    vegeta report -type=json "$RESULT_DIR/${test_name}_results.bin" > "$RESULT_DIR/${test_name}_report.json"
    vegeta plot "$RESULT_DIR/${test_name}_results.bin" > "$RESULT_DIR/${test_name}_plot.html"
    
    echo -e "${GREEN}✅ Test completed: $test_name${NC}"
    echo ""
    cat "$RESULT_DIR/${test_name}_report.txt"
    echo ""
}

# Основной процесс
get_token

# Создаем targets файлы с токеном
create_targets() {
    # Test 1: GET список объявлений
    cat > "$RESULT_DIR/targets_list.txt" <<EOF
GET $TARGET_URL/api/v1/offers
Authorization: Bearer $TOKEN

GET $TARGET_URL/api/v1/offers?offer_type=sale
Authorization: Bearer $TOKEN

GET $TARGET_URL/api/v1/offers?offer_type=sale&property_type=apartment
Authorization: Bearer $TOKEN

GET $TARGET_URL/api/v1/offers?offer_type=sale&property_type=apartment&rooms=2
Authorization: Bearer $TOKEN

GET $TARGET_URL/api/v1/offers?offer_type=rent&property_type=house
Authorization: Bearer $TOKEN
EOF

    # Test 2: GET конкретное объявление (нужны ID)
    # Получаем несколько ID из базы
    OFFER_IDS=$(psql "$DATABASE_URL" -t -c "SELECT id FROM offer WHERE status='active' LIMIT 10;" | tr -d ' ' | tr '\n' ' ')
    
    > "$RESULT_DIR/targets_single.txt"
    for offer_id in $OFFER_IDS; do
        echo "GET $TARGET_URL/api/v1/offers/$offer_id" >> "$RESULT_DIR/targets_single.txt"
        echo "Authorization: Bearer $TOKEN" >> "$RESULT_DIR/targets_single.txt"
        echo "" >> "$RESULT_DIR/targets_single.txt"
    done

    # Test 3: Смешанная нагрузка
    cat > "$RESULT_DIR/targets_mixed.txt" <<EOF
GET $TARGET_URL/api/v1/offers
Authorization: Bearer $TOKEN

GET $TARGET_URL/api/v1/offers?offer_type=sale
Authorization: Bearer $TOKEN

GET $TARGET_URL/api/v1/complexes/list
Authorization: Bearer $TOKEN

GET $TARGET_URL/api/v1/profile/$(psql "$DATABASE_URL" -t -c "SELECT user_id FROM profile LIMIT 1;" | tr -d ' ')
Authorization: Bearer $TOKEN
EOF
}

create_targets

echo -e "${GREEN}╔════════════════════════════════════╗${NC}"
echo -e "${GREEN}║   Starting Tests                   ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════╝${NC}"
echo ""

# Тест 1: Список объявлений
reset_db_stats
run_test "offers_list" "$RESULT_DIR/targets_list.txt"
save_db_stats "offers_list_db_stats.txt"

# Тест 2: Получение конкретных объявлений
reset_db_stats
run_test "offers_single" "$RESULT_DIR/targets_single.txt"
save_db_stats "offers_single_db_stats.txt"

# Тест 3: Смешанная нагрузка
reset_db_stats
run_test "mixed" "$RESULT_DIR/targets_mixed.txt"
save_db_stats "mixed_db_stats.txt"

echo -e "${GREEN}╔════════════════════════════════════╗${NC}"
echo -e "${GREEN}║   All Tests Completed!             ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════╝${NC}"
echo ""
echo -e "Results saved to: ${YELLOW}$RESULT_DIR${NC}"
echo ""
echo -e "View HTML reports:"
echo -e "  - open $RESULT_DIR/offers_list_plot.html"
echo -e "  - open $RESULT_DIR/offers_single_plot.html"
echo -e "  - open $RESULT_DIR/mixed_plot.html"
