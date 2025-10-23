DROP TRIGGER IF EXISTS set_updated_at_users ON users;
DROP TRIGGER IF EXISTS set_updated_at_region ON region;
DROP TRIGGER IF EXISTS set_updated_at_location ON location;
DROP TRIGGER IF EXISTS set_updated_at_metro_station ON metro_station;
DROP TRIGGER IF EXISTS set_updated_at_housing_complex ON housing_complex;
DROP TRIGGER IF EXISTS set_updated_at_profile ON profile;
DROP TRIGGER IF EXISTS set_updated_at_offer ON offer;
DROP TRIGGER IF EXISTS set_updated_at_offer_photo ON offer_photo;
DROP TRIGGER IF EXISTS set_updated_at_complex_photo ON complex_photo;

DELETE FROM offer_photo;
DELETE FROM complex_photo;
DELETE FROM location_metro;
DELETE FROM offer;
DELETE FROM housing_complex;
DELETE FROM profile;
DELETE FROM metro_station;
DELETE FROM location;
DELETE FROM users;
DELETE FROM region;

DROP TABLE IF EXISTS offer_photo;
DROP TABLE IF EXISTS complex_photo;
DROP TABLE IF EXISTS location_metro;
DROP TABLE IF EXISTS offer;
DROP TABLE IF EXISTS housing_complex;
DROP TABLE IF EXISTS profile;
DROP TABLE IF EXISTS metro_station;
DROP TABLE IF EXISTS location;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS region;

DROP TYPE IF EXISTS offer_type_enum;
DROP TYPE IF EXISTS offer_status_enum;
DROP TYPE IF EXISTS user_role_enum;
DROP TYPE IF EXISTS property_type_enum;

DROP FUNCTION IF EXISTS update_updated_at_column();