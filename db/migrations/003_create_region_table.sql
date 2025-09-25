CREATE TABLE region (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INT REFERENCES region(id) ON DELETE SET NULL,
    level INT NOT NULL DEFAULT 0 CHECK (level >= 0),
    slug VARCHAR(255) UNIQUE NOT NULL,
);

CREATE INDEX idx_region_parent_id ON region (parent_id);
CREATE INDEX idx_region_slug ON region (slug);
CREATE INDEX idx_region_level ON region (level);