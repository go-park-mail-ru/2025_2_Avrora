-- uuid
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Generic updated_at trigger function (used by many tables)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Enums (must come before tables that use them)
CREATE TYPE offer_type_enum AS ENUM ('sale', 'rent');
CREATE TYPE offer_status_enum AS ENUM ('active', 'sold', 'archived');
CREATE TYPE user_role_enum AS ENUM ('user', 'owner', 'realtor');
CREATE TYPE property_type_enum AS ENUM ('house', 'apartment');

-- users (root dependency)
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

-- region (hierarchical)
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

-- location (depends on region)
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
    avatar_url TEXT CHECK (LENGTH(avatar_url) <= 1024),
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
    url TEXT NOT NULL CHECK (LENGTH(url) <= 1024),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_offer_photo
    BEFORE UPDATE ON offer_photo
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE complex_photo (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    complex_id UUID NOT NULL REFERENCES housing_complex(id) ON DELETE CASCADE,
    url TEXT NOT NULL CHECK (LENGTH(url) <= 1024),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_updated_at_complex_photo
    BEFORE UPDATE ON complex_photo
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ✅ CRITICAL: offer_price_history must come BEFORE functions/triggers that use it
CREATE TABLE offer_price_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    offer_id UUID NOT NULL REFERENCES offer(id) ON DELETE CASCADE,
    old_price BIGINT,
    new_price BIGINT NOT NULL CHECK (new_price >= 0),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reason TEXT CHECK (LENGTH(reason) <= 500) DEFAULT NULL,
    changed_by UUID REFERENCES users(id) ON DELETE SET NULL
);

CREATE OR REPLACE FUNCTION get_offer_price_history(offer_uuid UUID)
RETURNS JSON AS $$
    SELECT COALESCE(
        (
            SELECT json_agg(
                json_build_object('date', ts, 'price', price_val)
                ORDER BY ts
            )
            FROM (
                -- Initial price
                SELECT o.created_at AS ts, o.price AS price_val
                FROM offer o
                WHERE o.id = offer_uuid

                UNION ALL

                -- Subsequent updates
                SELECT h.changed_at AS ts, h.new_price AS price_val
                FROM offer_price_history h
                WHERE h.offer_id = offer_uuid
            ) AS history
        ),
        '[]'::json
    );
$$ LANGUAGE sql STABLE;

-- Trigger function for initial price (uses offer.created_at)
CREATE OR REPLACE FUNCTION log_initial_offer_price()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO offer_price_history (
        offer_id,
        old_price,
        new_price,
        changed_at,
        reason,
        changed_by
    ) VALUES (
        NEW.id,
        NULL,
        NEW.price,
        NEW.created_at,  -- ✅ consistent with offer creation time
        'listed',
        NEW.user_id
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for price updates (uses offer.updated_at for consistency)
CREATE OR REPLACE FUNCTION log_offer_price_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.price IS DISTINCT FROM NEW.price THEN
        INSERT INTO offer_price_history (
            offer_id,
            old_price,
            new_price,
            changed_at,
            reason,
            changed_by
        ) VALUES (
            NEW.id,
            OLD.price,
            NEW.price,
            NEW.updated_at,  -- ✅ use updated_at (set by trigger) instead of NOW()
            'price_updated',
            NEW.user_id
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ✅ Triggers (now safe — table and functions exist)
CREATE TRIGGER trigger_log_initial_offer_price
    AFTER INSERT ON offer
    FOR EACH ROW
    EXECUTE FUNCTION log_initial_offer_price();

CREATE TRIGGER trigger_log_offer_price_change
    AFTER UPDATE OF price ON offer
    FOR EACH ROW
    EXECUTE FUNCTION log_offer_price_change();


-- likes, views

-- View tracking (one row per view event)
CREATE TABLE offer_view (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    offer_id UUID NOT NULL REFERENCES offer(id) ON DELETE CASCADE,
    viewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Like tracking (one row per user per offer)
CREATE TABLE offer_like (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    offer_id UUID NOT NULL REFERENCES offer(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    liked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (offer_id, user_id)  -- Prevent duplicate likes
);

-- Function to log a view (can be called from application)
CREATE OR REPLACE FUNCTION log_offer_view(
    offer_uuid UUID
)
RETURNS void AS $$
BEGIN
    INSERT INTO offer_view (offer_id)
    VALUES (offer_uuid);
END;
$$ LANGUAGE plpgsql;

-- Function to toggle like (add/remove)
CREATE OR REPLACE FUNCTION toggle_offer_like(
    offer_uuid UUID,
    user_uuid UUID
)
RETURNS void AS $$
BEGIN
    -- Check if like exists
    IF EXISTS (SELECT 1 FROM offer_like WHERE offer_id = offer_uuid AND user_id = user_uuid) THEN
        -- Remove like
        DELETE FROM offer_like WHERE offer_id = offer_uuid AND user_id = user_uuid;
    ELSE
        -- Add like
        INSERT INTO offer_like (offer_id, user_id)
        VALUES (offer_uuid, user_uuid);
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_offer_view_count(offer_uuid UUID)
RETURNS BIGINT AS $$
    SELECT COUNT(*) FROM offer_view WHERE offer_id = offer_uuid;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_offer_like_count(offer_uuid UUID)
RETURNS BIGINT AS $$
    SELECT COUNT(*) FROM offer_like WHERE offer_id = offer_uuid;
$$ LANGUAGE sql STABLE;

-- For getting user's like status
CREATE OR REPLACE FUNCTION is_offer_liked(offer_uuid UUID, user_uuid UUID)
RETURNS BOOLEAN AS $$
    SELECT EXISTS (SELECT 1 FROM offer_like WHERE offer_id = offer_uuid AND user_id = user_uuid);
$$ LANGUAGE sql STABLE;