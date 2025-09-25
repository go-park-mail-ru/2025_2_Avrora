CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE region (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INT REFERENCES region(id) ON DELETE SET NULL,
    level INT NOT NULL DEFAULT 0 CHECK (level >= 0),
    slug VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE location (
    id SERIAL PRIMARY KEY,
    region_id INT NOT NULL REFERENCES region(id) ON DELETE CASCADE,
    street VARCHAR(255) NOT NULL,
    house_number VARCHAR(50) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8)
);

CREATE TABLE offer (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    location_id INT NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES category(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price INT NOT NULL CHECK (price >= 0),
    area DECIMAL(10, 2) CHECK (area > 0),  
    rooms INT CHECK (rooms >= 0),
    offer_type VARCHAR(20) NOT NULL DEFAULT 'sale' CHECK (offer_type IN ('sale', 'rent')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE photo (
    id SERIAL PRIMARY KEY,
    offer_id INT NOT NULL REFERENCES offer(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    position INT NOT NULL DEFAULT 0
);