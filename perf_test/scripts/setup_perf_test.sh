#!/bin/bash

set -e

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}╔════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                                                    ║${NC}"
echo -e "${GREEN}║   Performance Testing Directory Setup Script      ║${NC}"
echo -e "${GREEN}║                                                    ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════╝${NC}"
echo ""

# Создаем основную структуру директорий
echo -e "${BLUE}📁 Creating directory structure...${NC}"

mkdir -p perf_test/{scripts,configs,sql,results,payloads,profiles}
mkdir -p perf_test/results/{iteration_1,iteration_2,iteration_3}

echo -e "${GREEN}✅ Directory structure created${NC}"
echo ""

# ============================================================================
# 1. README.md
# ============================================================================
echo -e "${BLUE}📝 Creating README.md...${NC}"

cat > perf_test/README.md << 'EOF'
# Отчет о нагрузочном тестировании

## Описание приложения

**Основная сущность:** Объявления о недвижимости (`offer`)

**Обоснование выбора:** Объявления (offer) являются центральной бизнес-сущностью приложения недвижимости. Они связывают пользователей, локации, жилые комплексы и содержат основную информацию для пользователей.

**Тестируемые API:**
- `POST /api/v1/offers/create` - создание объявления
- `GET /api/v1/offers` - получение списка объявлений (с фильтрами)
- `GET /api/v1/offers/{id}` - получение конкретного объявления
- `PUT /api/v1/offers/update/{id}` - обновление объявления
- `DELETE /api/v1/offers/delete/{id}` - удаление объявления
- `GET /api/v1/profile/myoffers/{user_id}` - получение объявлений пользователя

---

## Инструменты

- **Нагрузочное тестирование:** `vegeta` (https://github.com/tsenart/vegeta)
- **Мониторинг:** `htop`, `pg
