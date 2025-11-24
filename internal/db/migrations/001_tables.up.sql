CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Enums
CREATE TYPE offer_type_enum AS ENUM ('sale', 'rent');
CREATE TYPE offer_status_enum AS ENUM ('active', 'sold', 'archived');
CREATE TYPE user_role_enum AS ENUM ('user', 'owner', 'realtor');
CREATE TYPE property_type_enum AS ENUM ('house', 'apartment');

-- users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE
        CHECK (LENGTH(email) <= 255)
        CHECK (email ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'),
    password_hash TEXT NOT NULL CHECK (LENGTH(password_hash) <= 255),
    role user_role_enum NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_users
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- region
CREATE TABLE region (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL CHECK (LENGTH(name) <= 255),
    parent_id UUID REFERENCES region(id) ON DELETE SET NULL,
    level INT NOT NULL CHECK (level >= 0),
    slug TEXT NOT NULL UNIQUE CHECK (LENGTH(slug) <= 255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_region
    BEFORE UPDATE ON region
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- location
CREATE TABLE location (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    region_id UUID NOT NULL REFERENCES region(id) ON DELETE CASCADE,
    latitude DECIMAL(10,8) CHECK (latitude BETWEEN -90 AND 90),
    longitude DECIMAL(11,8) CHECK (longitude BETWEEN -180 AND 180),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_location
    BEFORE UPDATE ON location
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- metro_station (depends on location)
CREATE TABLE metro_station (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL CHECK (LENGTH(name) <= 100),
    location_id UUID NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_metro_station
    BEFORE UPDATE ON metro_station
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- location_metro (many-to-many)
CREATE TABLE location_metro (
    location_id UUID NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    metro_station_id UUID NOT NULL REFERENCES metro_station(id) ON DELETE CASCADE,
    distance_meters INT NOT NULL CHECK (distance_meters >= 0),
    PRIMARY KEY (location_id, metro_station_id)
);

-- housing_complex (depends on location)
CREATE TABLE housing_complex (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL CHECK (LENGTH(name) <= 255),
    description TEXT CHECK (LENGTH(description) <= 2000),
    year_built INT CHECK (year_built BETWEEN 1800 AND 2100),
    location_id UUID NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    developer TEXT CHECK (LENGTH(developer) <= 255),
    address TEXT CHECK (LENGTH(address) <= 255),
    starting_price BIGINT CHECK (starting_price >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_housing_complex
    BEFORE UPDATE ON housing_complex
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- profile (depends on users)
CREATE TABLE profile (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    first_name TEXT CHECK (LENGTH(first_name) <= 100),
    last_name TEXT CHECK (LENGTH(last_name) <= 100),
    phone TEXT CHECK (LENGTH(phone) <= 20),
    avatar_url TEXT
        CHECK (LENGTH(avatar_url) <= 1024),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_profile
    BEFORE UPDATE ON profile
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- offer (depends on users, location, housing_complex)
CREATE TABLE offer (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    housing_complex_id UUID REFERENCES housing_complex(id) ON DELETE SET NULL,
    title TEXT NOT NULL CHECK (LENGTH(title) <= 255),
    description TEXT CHECK (LENGTH(description) <= 5000),
    price BIGINT NOT NULL CHECK (price >= 0),
    area DECIMAL(10,2) NOT NULL CHECK (area > 0),
    address TEXT NOT NULL CHECK (LENGTH(address) <= 255),
    rooms INT NOT NULL CHECK (rooms >= 0),
    property_type property_type_enum NOT NULL,
    offer_type offer_type_enum NOT NULL,
    status offer_status_enum NOT NULL DEFAULT 'active',
    floor INT CHECK (floor >= 0),
    total_floors INT CHECK (total_floors >= 0),
    deposit BIGINT CHECK (deposit >= 0),
    commission BIGINT CHECK (commission >= 0),
    rental_period TEXT CHECK (LENGTH(rental_period) <= 100),
    living_area DECIMAL(10,2) CHECK (living_area >= 0),
    kitchen_area DECIMAL(10,2) CHECK (kitchen_area >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_offer
    BEFORE UPDATE ON offer
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Photos
CREATE TABLE offer_photo (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    offer_id UUID NOT NULL REFERENCES offer(id) ON DELETE CASCADE,
    url TEXT NOT NULL
        CHECK (LENGTH(url) <= 1024),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_offer_photo
    BEFORE UPDATE ON offer_photo
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE complex_photo (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    complex_id UUID NOT NULL REFERENCES housing_complex(id) ON DELETE CASCADE,
    url TEXT NOT NULL
        CHECK (LENGTH(url) <= 1024),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_complex_photo
    BEFORE UPDATE ON complex_photo
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
--Добавление лайков
CREATE TABLE offer_like (
                            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                            offer_id UUID NOT NULL REFERENCES offer(id) ON DELETE CASCADE,
                            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                            PRIMARY KEY (user_id, offer_id)
);
ALTER TABLE offer ADD COLUMN likes_count INT NOT NULL DEFAULT 0;
CREATE OR REPLACE FUNCTION update_likes_count()
    RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE offer SET likes_count = likes_count + 1 WHERE id = NEW.offer_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE offer SET likes_count = likes_count - 1 WHERE id = OLD.offer_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_likes_count
    AFTER INSERT OR DELETE ON offer_like
    FOR EACH ROW EXECUTE FUNCTION update_likes_count();


