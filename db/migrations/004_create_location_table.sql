CREATE TABLE location (
    id SERIAL PRIMARY KEY,
    region_id INT NOT NULL REFERENCES region(id) ON DELETE CASCADE,
    street VARCHAR(255) NOT NULL,
    house_number VARCHAR(50) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
);

CREATE INDEX idx_location_coords ON location (latitude, longitude);