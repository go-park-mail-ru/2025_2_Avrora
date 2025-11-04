-- Enable UUID extension (already in your schema, but safe to repeat)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Insert regions (hierarchical: country → federal district → city)
INSERT INTO region (id, name, parent_id, level, slug) VALUES
  ('10000000-0000-0000-0000-000000000001', 'Russia', NULL, 0, 'russia'),
  ('10000000-0000-0000-0000-000000000002', 'Moscow', '10000000-0000-0000-0000-000000000001', 1, 'moscow');

-- Insert locations
INSERT INTO location (id, region_id, latitude, longitude) VALUES
  ('20000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000002', 55.7558, 37.6176), -- Moscow center
  ('20000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000002', 55.7512, 37.6184); -- Nearby

-- Insert metro stations
INSERT INTO metro_station (id, name, location_id) VALUES
  ('30000000-0000-0000-0000-000000000001', 'Teatralnaya', '20000000-0000-0000-0000-000000000001'),
  ('30000000-0000-0000-0000-000000000002', 'Okhotny Ryad', '20000000-0000-0000-0000-000000000002');

-- Link locations to metro stations
INSERT INTO location_metro (location_id, metro_station_id, distance_meters) VALUES
  ('20000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 200),
  ('20000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000002', 400);

-- Insert users
INSERT INTO users (id, email, password_hash, role) VALUES
  ('40000000-0000-0000-0000-000000000001', 'user@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user'),
  ('40000000-0000-0000-0000-000000000002', 'owner@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'owner'),
  ('40000000-0000-0000-0000-000000000003', 'realtor@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'realtor');

-- Insert profiles
INSERT INTO profile (user_id, first_name, last_name, phone, avatar_url) VALUES
  ('40000000-0000-0000-0000-000000000001', 'Alex', 'Petrov', '+79001234567', 'http://37.139.40.252:8080/api/v1/image/default_avatar.jpg'),
  ('40000000-0000-0000-0000-000000000002', 'Maria', 'Ivanova', '+79007654321', 'http://37.139.40.252:8080/api/v1/image/default_avatar.jpg');

-- Insert housing complexes
INSERT INTO housing_complex (id, name, description, year_built, location_id, developer, address, starting_price) VALUES
  ('50000000-0000-0000-0000-000000000001', 'Moscow City Residences', 'Luxury apartments in the business district', 2020,
   '20000000-0000-0000-0000-000000000001', 'Capital Development', 'Presnenskaya Embankment, 10', 15000000);

-- Insert complex photos
INSERT INTO complex_photo (complex_id, url) VALUES
  ('50000000-0000-0000-0000-000000000001', 'http://37.139.40.252:8080/api/v1/image/default_offer.jpg'),
  ('50000000-0000-0000-0000-000000000001', 'http://37.139.40.252:8080/api/v1/image/default_offer.jpg');

-- Insert offers
INSERT INTO offer (
  id, user_id, location_id, housing_complex_id, title, description, price, area,
  address, rooms, property_type, offer_type, status, floor, total_floors,
  deposit, commission, rental_period, living_area, kitchen_area
) VALUES
  ('60000000-0000-0000-0000-000000000001',
   '40000000-0000-0000-0000-000000000002',
   '20000000-0000-0000-0000-000000000001',
   '50000000-0000-0000-0000-000000000001',
   '2-Bedroom Apartment in Moscow City',
   'Modern apartment with panoramic views',
   25000000, 75.50,
   'Presnenskaya Embankment, 10, Apt 1205', 2, 'apartment', 'sale', 'active',
   12, 42, NULL, NULL, NULL, 50.00, 12.50),

  ('60000000-0000-0000-0000-000000000002',
   '40000000-0000-0000-0000-000000000003',
   '20000000-0000-0000-0000-000000000002',
   NULL,
   'Cozy Studio for Rent near Metro',
   'Fully furnished studio, 5 min walk to metro',
   65000, 32.00,
   'Tverskaya St, 15', 0, 'apartment', 'rent', 'active',
   3, 9, 65000, 10000, '12 months', 22.00, 8.00);

-- Insert offer photos
INSERT INTO offer_photo (offer_id, url) VALUES
  ('60000000-0000-0000-0000-000000000001', 'http://37.139.40.252:8080/api/v1/image/default_offer.jpg'),
  ('60000000-0000-0000-0000-000000000001', 'http://37.139.40.252:8080/api/v1/image/default_offer.jpg'),
  ('60000000-0000-0000-0000-000000000002', 'http://37.139.40.252:8080/api/v1/image/default_offer.jpg');

-- Дополнительные жилищные комплексы

INSERT INTO housing_complex (id, name, description, year_built, location_id, developer, address, starting_price) VALUES
  ('50000000-0000-0000-0000-000000000002', 'Garden Quarters', 'Экологичный ЖК с парковой зоной и детской инфраструктурой', 2022,
   '20000000-0000-0000-0000-000000000002', 'PIK Group', 'Kutuzovsky Prospekt, 44', 12000000),

  ('50000000-0000-0000-0000-000000000003', 'Neva Tower Residences', 'Элитные апартаменты в высотке бизнес-класса', 2021,
   '20000000-0000-0000-0000-000000000001', 'LSR Group', 'Krasnopresnenskaya Embankment, 12', 22000000);

-- Фото для новых жилищных комплексов

INSERT INTO complex_photo (complex_id, url) VALUES
  ('50000000-0000-0000-0000-000000000002', 'http://37.139.40.252:8080/api/v1/image/default_complex.jpg'),
  ('50000000-0000-0000-0000-000000000002', 'http://37.139.40.252:8080/api/v1/image/default_complex.jpg'),
  ('50000000-0000-0000-0000-000000000003', 'http://37.139.40.252:8080/api/v1/image/default_complex.jpg');

-- Дополнительные объявления (offers)

INSERT INTO offer (
  id, user_id, location_id, housing_complex_id, title, description, price, area,
  address, rooms, property_type, offer_type, status, floor, total_floors,
  deposit, commission, rental_period, living_area, kitchen_area
) VALUES
  -- Продажа в Garden Quarters
  ('60000000-0000-0000-0000-000000000003',
   '40000000-0000-0000-0000-000000000002',
   '20000000-0000-0000-0000-000000000002',
   '50000000-0000-0000-0000-000000000002',
   '3-Bedroom Apartment in Garden Quarters',
   'Светлая квартира с отделкой от застройщика, большая терраса',
   18500000, 92.30,
   'Kutuzovsky Prospekt, 44, Corp. 2, Apt 312', 3, 'apartment', 'sale', 'active',
   3, 18, NULL, NULL, NULL, 68.00, 14.00),

  -- Аренда студии без привязки к ЖК
  ('60000000-0000-0000-0000-000000000004',
   '40000000-0000-0000-0000-000000000003',
   '20000000-0000-0000-0000-000000000001',
   NULL,
   'Уютная студия рядом с Театральной',
   'Современная мебель, техника, консьерж, видеонаблюдение',
   85000, 28.50,
   'Bolshaya Dmitrovka, 7', 0, 'apartment', 'rent', 'active',
   5, 7, 85000, 12000, '11 months', 20.00, 6.50),

  -- Продажа в Neva Tower Residences
  ('60000000-0000-0000-0000-000000000005',
   '40000000-0000-0000-0000-000000000002',
   '20000000-0000-0000-0000-000000000001',
   '50000000-0000-0000-0000-000000000003',
   'Пентхаус с панорамным видом',
   '300 м², панорамные окна, высота 67 этажа, частный лифт',
   320000000, 300.00,
   'Krasnopresnenskaya Embankment, 12, Penthouse', 4, 'apartment', 'sale', 'active',
   67, 70, NULL, NULL, NULL, 220.00, 30.00),

  -- Аренда 2-комнатной квартиры в центре
  ('60000000-0000-0000-0000-000000000006',
   '40000000-0000-0000-0000-000000000003',
   '20000000-0000-0000-0000-000000000002',
   NULL,
   '2-комнатная квартира в Тверском районе',
   'Отличное состояние, рядом метро, парковка включена',
   130000, 58.00,
   'Tverskaya St, 22', 2, 'apartment', 'rent', 'active',
   4, 12, 130000, 15000, '12 months', 42.00, 10.00);

-- Фото для новых объявлений

INSERT INTO offer_photo (offer_id, url) VALUES
  ('60000000-0000-0000-0000-000000000003', 'http://37.139.40.252:8080/api/v1/image/default_offer.jpg'),
  ('60000000-0000-0000-0000-000000000003', 'http://37.139.40.252:8080/api/v1/image/default_offer1.jpg'),
  ('60000000-0000-0000-0000-000000000004', 'http://37.139.40.252:8080/api/v1/image/default_offer2.jpg'),
  ('60000000-0000-0000-0000-000000000005', 'http://37.139.40.252:8080/api/v1/image/default_offer3.jpg'),
  ('60000000-0000-0000-0000-000000000005', 'http://37.139.40.252:8080/api/v1/image/default_offer4.jpg'),
  ('60000000-0000-0000-0000-000000000005', 'http://37.139.40.252:8080/api/v1/image/default_offer5.jpg'),
  ('60000000-0000-0000-0000-000000000006', 'http://37.139.40.252:8080/api/v1/image/default_offer6.jpg');