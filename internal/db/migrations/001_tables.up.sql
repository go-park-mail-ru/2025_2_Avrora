CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TYPE offer_type_enum AS ENUM ('sale', 'rent');
CREATE TYPE offer_status_enum AS ENUM ('active', 'sold', 'archived');
CREATE TYPE user_role_enum AS ENUM ('user', 'owner', 'realtor');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE 
        CHECK (LENGTH(email) <= 255)
        CHECK (email ~ '^[\p{L}\p{N}._%+-]+@[\p{L}\p{N}.-]+\.[\p{L}]{2,}$'), -- Регулярка совпадает с тем что на бэке и фронте
    password_hash TEXT NOT NULL CHECK (LENGTH(password_hash) <= 255),
    avatar_url TEXT 
        CHECK (LENGTH(avatar_url) <= 1024)
        CHECK (avatar_url IS NULL OR avatar_url ~ '^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$'),
    role user_role_enum NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_users 
    BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE category (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL CHECK (LENGTH(name) <= 100),
    slug TEXT NOT NULL UNIQUE CHECK (LENGTH(slug) <= 100),
    description TEXT CHECK (LENGTH(description) <= 1000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_category 
    BEFORE UPDATE ON category 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

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

CREATE TABLE metro_station (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL CHECK (LENGTH(name) <= 100),
    latitude DECIMAL(10,8) NOT NULL CHECK (latitude BETWEEN -90 AND 90),
    longitude DECIMAL(11,8) NOT NULL CHECK (longitude BETWEEN -180 AND 180),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_metro_station 
    BEFORE UPDATE ON metro_station 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE housing_complex (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL CHECK (LENGTH(name) <= 255),
    description TEXT CHECK (LENGTH(description) <= 2000),
    year_built INT CHECK (year_built BETWEEN 1800 AND 2100),
    region_id UUID NOT NULL REFERENCES region(id) ON DELETE CASCADE,
    latitude DECIMAL(10,8) CHECK (latitude BETWEEN -90 AND 90),
    longitude DECIMAL(11,8) CHECK (longitude BETWEEN -180 AND 180),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_housing_complex 
    BEFORE UPDATE ON housing_complex 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE location (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    region_id UUID NOT NULL REFERENCES region(id) ON DELETE CASCADE,
    housing_complex_id UUID REFERENCES housing_complex(id) ON DELETE SET NULL,
    street TEXT NOT NULL CHECK (LENGTH(street) <= 255),
    house_number TEXT NOT NULL CHECK (LENGTH(house_number) <= 50),
    latitude DECIMAL(10,8) CHECK (latitude BETWEEN -90 AND 90),
    longitude DECIMAL(11,8) CHECK (longitude BETWEEN -180 AND 180),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_location 
    BEFORE UPDATE ON location 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE location_metro (
    location_id UUID NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    metro_station_id UUID NOT NULL REFERENCES metro_station(id) ON DELETE CASCADE,
    distance_meters INT NOT NULL CHECK (distance_meters >= 0),
    PRIMARY KEY (location_id, metro_station_id)
);

CREATE TABLE offer (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES category(id) ON DELETE CASCADE,
    title TEXT NOT NULL CHECK (LENGTH(title) <= 255),
    description TEXT CHECK (LENGTH(description) <= 5000),
    price BIGINT NOT NULL CHECK (price >= 0),
    area DECIMAL(10,2) NOT NULL CHECK (area > 0),
    rooms INT NOT NULL CHECK (rooms >= 0),
    offer_type offer_type_enum NOT NULL,
    status offer_status_enum NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_offer 
    BEFORE UPDATE ON offer 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE photo (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    offer_id UUID NOT NULL REFERENCES offer(id) ON DELETE CASCADE,
    url TEXT NOT NULL 
        CHECK (LENGTH(url) <= 1024)
        CHECK (url ~ '^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_photo 
    BEFORE UPDATE ON photo 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

/*CREATE TABLE review (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    offer_id UUID NOT NULL REFERENCES offer(id) ON DELETE CASCADE,
    rating INT CHECK (rating BETWEEN 1 AND 5),
    comment TEXT CHECK (LENGTH(comment) <= 1000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_review 
    BEFORE UPDATE ON review 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
 */