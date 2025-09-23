-- db/seed/mocks.sql
CREATE UNIQUE INDEX IF NOT EXISTS idx_location_unique_address ON location (region_id, street, house_number);
-- Категории
INSERT INTO category (name, slug, description) VALUES
    ('Квартира', 'kvartira', 'Квартиры на продажу и в аренду');

-- Регионы
INSERT INTO region (name, parent_id, level, slug) VALUES
    ('Россия', NULL, 0, 'rossiya'),
    ('Москва', (SELECT id FROM region WHERE slug = 'rossiya'), 1, 'moskva');

-- Локации
INSERT INTO location (region_id, street, house_number, latitude, longitude) VALUES
    (
        (SELECT id FROM region WHERE slug = 'moskva'),
        'Тверская улица',
        '15',
        55.760597,
        37.617870
    );

-- Пользователи
INSERT INTO users (email, password) VALUES
    ('user@example.com', 'hashed_pass');

-- Объявления
INSERT INTO offer (user_id, location_id, category_id, title, description, price, area, rooms, address, offer_type) VALUES
    (
        (SELECT id FROM users WHERE email = 'user@example.com'),
        (SELECT id FROM location WHERE street = 'Тверская улица' AND house_number = '15'),
        (SELECT id FROM category WHERE slug = 'kvartira'),
        'Продам квартиру на Тверской',
        'Отличный ремонт',
        100000000,
        40.0,
        1,
        'Тверская улица, 15',
        'sale'
    );