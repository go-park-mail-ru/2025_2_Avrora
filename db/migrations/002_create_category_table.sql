-- Создание таблицы категорий недвижимости
CREATE TABLE IF NOT EXISTS category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Индекс по slug для URL и поиска
CREATE INDEX IF NOT EXISTS idx_category_slug ON category (slug);