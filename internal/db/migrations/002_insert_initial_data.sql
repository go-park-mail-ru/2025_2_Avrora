-- Insert base regions (Country → City)
INSERT INTO region (id, name, parent_id, level, slug) VALUES
-- Russia (level 0)
('11111111-1111-1111-1111-111111111111', 'Russia', NULL, 0, 'russia'),
-- Moscow (level 1, child of Russia)
('22222222-2222-2222-2222-222222222222', 'Moscow', '11111111-1111-1111-1111-111111111111', 1, 'moscow'),
-- Central Administrative District (level 2)
('33333333-3333-3333-3333-333333333333', 'Central Administrative District', '22222222-2222-2222-2222-222222222222', 2, 'moscow-central'),
-- Tverskoy District (level 3)
('44444444-4444-4444-4444-444444444444', 'Tverskoy District', '33333333-3333-3333-3333-333333333333', 3, 'tverskoy');

-- Insert categories
INSERT INTO category (id, name, slug, description) VALUES
('55555555-5555-5555-5555-555555555555', 'Apartment', 'apartment', 'Residential apartment units'),
('66666666-6666-6666-6666-666666666666', 'House', 'house', 'Private residential houses'),
('77777777-7777-7777-7777-777777777777', 'Room', 'room', 'Single room in shared accommodation');

-- Insert a metro station in Tverskoy District
INSERT INTO metro_station (id, name, latitude, longitude, region_id) VALUES
('88888888-8888-8888-8888-888888888888', 'Tverskaya', 55.759774, 37.611538, '44444444-4444-4444-4444-444444444444');

INSERT INTO housing_complex (id, name, description, year_built, region_id, latitude, longitude) VALUES
('cccccccc-cccc-cccc-cccc-cccccccccccc', 
 'Tverskoy City Residence', 
 'Premium residential complex with 24/7 security and underground parking.',
 2018,
 '44444444-4444-4444-4444-444444444444',
 55.760500,
 37.612500);

-- Insert a location linked to the housing complex
INSERT INTO location (id, region_id, housing_complex_id, street, house_number, latitude, longitude) VALUES
('99999999-9999-9999-9999-999999999999', 
 '44444444-4444-4444-4444-444444444444',
 'cccccccc-cccc-cccc-cccc-cccccccccccc',
 'Tverskaya Street', 
 '15', 
 55.760123, 
 37.612045);
 
-- Insert a sample user (password_hash is bcrypt hash of "password123")
INSERT INTO users (id, email, password_hash, role) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'owner@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'owner');

-- Insert a sample offer
INSERT INTO offer (id, user_id, location_id, category_id, title, description, price, area, rooms, offer_type, status) VALUES
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 
 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
 '99999999-9999-9999-9999-999999999999',
 '55555555-5555-5555-5555-555555555555',
 'Luxury Apartment on Tverskaya',
 'Spacious 2-room apartment in the heart of Moscow, 5 min walk to Tverskaya metro.',
 12500000,
 65.5,
 2,
 'sale',
 'active');

-- Insert photos for the offer
INSERT INTO photo (offer_id, url) VALUES
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'https://example.com/photos/apt1_main.jpg'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'https://example.com/photos/apt1_kitchen.jpg'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'https://example.com/photos/apt1_bedroom.jpg');