-- Включаем расширение для статистики
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Сбрасываем статистику перед тестом
SELECT pg_stat_statements_reset();
