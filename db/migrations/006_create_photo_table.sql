-- Создание таблицы фотографий
CREATE TABLE IF NOT EXISTS photo (
    id SERIAL PRIMARY KEY,
    offer_id INT NOT NULL REFERENCES offer(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    position INT NOT NULL DEFAULT 0,
    uploaded_at TIMESTAMPTZ DEFAULT NOW()
);

-- Уникальные индексы: позиция и URL уникальны в рамках объявления
CREATE UNIQUE INDEX IF NOT EXISTS idx_photo_offer_position ON photo (offer_id, position);
CREATE UNIQUE INDEX IF NOT EXISTS idx_photo_offer_url ON photo (offer_id, url);

-- Индекс для быстрой выборки фото по объявлению
CREATE INDEX IF NOT EXISTS idx_photo_offer_id ON photo (offer_id);