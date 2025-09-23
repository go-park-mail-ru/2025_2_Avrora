CREATE TABLE IF NOT EXISTS location (
    id SERIAL PRIMARY KEY,
    region_id INT NOT NULL REFERENCES region(id) ON DELETE CASCADE,
    street VARCHAR(255) NOT NULL,
    house_number VARCHAR(50) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_location_unique_address ON location (region_id, street, house_number);
CREATE INDEX IF NOT EXISTS idx_location_coords ON location (latitude, longitude);