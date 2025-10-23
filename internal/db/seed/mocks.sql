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
  ('40000000-0000-0000-0000-000000000001', 'Alex', 'Petrov', '+79001234567', 'https://example.com/avatars/alex.jpg'),
  ('40000000-0000-0000-0000-000000000002', 'Maria', 'Ivanova', '+79007654321', 'https://example.com/avatars/maria.jpg');

-- Insert housing complexes
INSERT INTO housing_complex (id, name, description, year_built, location_id, developer, address, starting_price) VALUES
  ('50000000-0000-0000-0000-000000000001', 'Moscow City Residences', 'Luxury apartments in the business district', 2020,
   '20000000-0000-0000-0000-000000000001', 'Capital Development', 'Presnenskaya Embankment, 10', 15000000);

-- Insert complex photos
INSERT INTO complex_photo (complex_id, url) VALUES
  ('50000000-0000-0000-0000-000000000001', 'https://example.com/photos/complex1_1.jpg'),
  ('50000000-0000-0000-0000-000000000001', 'https://example.com/photos/complex1_2.jpg');

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
  ('60000000-0000-0000-0000-000000000001', 'https://example.com/photos/offer1_1.jpg'),
  ('60000000-0000-0000-0000-000000000001', 'https://example.com/photos/offer1_2.jpg'),
  ('60000000-0000-0000-0000-000000000002', 'https://example.com/photos/offer2_1.jpg');