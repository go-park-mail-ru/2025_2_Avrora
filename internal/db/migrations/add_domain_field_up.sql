ALTER TABLE users
    ADD COLUMN first_name TEXT CHECK (LENGTH(first_name) <= 100),
    ADD COLUMN last_name TEXT CHECK (LENGTH(last_name) <= 100),
    ADD COLUMN phone TEXT CHECK (LENGTH(phone) <= 20),
    ADD COLUMN photo_url TEXT;

ALTER TABLE housing_complex
    ADD COLUMN developer TEXT CHECK (LENGTH(developer) <= 255),
    ADD COLUMN address TEXT CHECK (LENGTH(address) <= 255),
    ADD COLUMN starting_price BIGINT CHECK (starting_price >= 0),
    ADD COLUMN image_urls TEXT[];


ALTER TABLE offer
    ADD COLUMN image TEXT,
    ADD COLUMN floor INT CHECK (floor >= 0),
    ADD COLUMN total_floors INT CHECK (total_floors >= 0),
    ADD COLUMN deposit BIGINT CHECK (deposit >= 0),
    ADD COLUMN commission BIGINT CHECK (commission >= 0),
    ADD COLUMN rental_period TEXT CHECK (LENGTH(rental_period) <= 100),
    ADD COLUMN living_area  CHECK (living_area >= 0),
    ADD COLUMN kitchen_area  CHECK (kitchen_area >= 0);
