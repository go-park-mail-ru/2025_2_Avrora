-- Создание таблицы регионов (иерархическая)
CREATE TABLE IF NOT EXISTS region (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INT REFERENCES region(id) ON DELETE SET NULL,
    level INT NOT NULL DEFAULT 0 CHECK (level >= 0),
    slug VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Индексы для быстрой навигации
CREATE INDEX IF NOT EXISTS idx_region_parent_id ON region (parent_id);
CREATE INDEX IF NOT EXISTS idx_region_slug ON region (slug);
CREATE INDEX IF NOT EXISTS idx_region_level ON region (level);