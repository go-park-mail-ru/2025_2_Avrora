INSERT INTO category (name, slug, description) VALUES
('Квартира', 'kvartira', 'Квартиры на продажу и в аренду'),
('Дом', 'dom', 'Частные дома и коттеджи'),
('Комната', 'komnata', 'Отдельные комнаты в квартирах и общежитиях'),
('Коммерческая недвижимость', 'commercial', 'Офисы, магазины, склады');

INSERT INTO region (name, parent_id, level, slug) VALUES
('Россия', NULL, 0, 'rossiya'),
('Москва', 1, 1, 'moskva'),
('Центральный административный округ', 2, 2, 'cao'),
('Северный административный округ', 2, 2, 'sao');

INSERT INTO location (region_id, street, house_number, latitude, longitude) VALUES
(3, 'Тверская улица', '15', 55.760597, 37.617870),
(3, 'Никольская улица', '10', 55.754430, 37.622940);

INSERT INTO users (email, password) VALUES
('user@example.com', 'hashed_password_here');

INSERT INTO offer (
    user_id, location_id, category_id, title, description, price, area, rooms, offer_type
) VALUES (
    1,
    1,
    1,
    'Продам 2-комнатную квартиру на Тверской',
    'Отличный ремонт, вид на Кремль, консьерж, парковка.',
    2500000,
    54.5,
    2,
    'sale'
);

INSERT INTO photo (offer_id, url, position) VALUES
(1, 'https://example.com/photo1.jpg', 1),
(1, 'https://example.com/photo2.jpg', 2);