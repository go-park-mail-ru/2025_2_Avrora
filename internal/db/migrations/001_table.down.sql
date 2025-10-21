DROP TRIGGER IF EXISTS set_updated_at_users ON users;
DROP TRIGGER IF EXISTS set_updated_at_category ON category;
DROP TRIGGER IF EXISTS set_updated_at_region ON region;
DROP TRIGGER IF EXISTS set_updated_at_metro_station ON metro_station;
DROP TRIGGER IF EXISTS set_updated_at_housing_complex ON housing_complex;
DROP TRIGGER IF EXISTS set_updated_at_location ON location;
DROP TRIGGER IF EXISTS set_updated_at_offer ON offer;
DROP TRIGGER IF EXISTS set_updated_at_photo ON photo;
DROP TRIGGER IF EXISTS set_updated_at_review ON review;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS review;
DROP TABLE IF EXISTS photo;
DROP TABLE IF EXISTS offer;
DROP TABLE IF EXISTS location_metro;
DROP TABLE IF EXISTS location;
DROP TABLE IF EXISTS housing_complex;
DROP TABLE IF EXISTS metro_station;
DROP TABLE IF EXISTS region;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS offer_status_enum;
DROP TYPE IF EXISTS offer_type_enum;
DROP TYPE IF EXISTS user_role_enum;