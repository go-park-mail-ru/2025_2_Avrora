-- +migrate Down

ALTER TABLE offer
    DROP COLUMN IF EXISTS image,
    DROP COLUMN IF EXISTS floor,
    DROP COLUMN IF EXISTS total_floors,
    DROP COLUMN IF EXISTS deposit,
    DROP COLUMN IF EXISTS commission,
    DROP COLUMN IF EXISTS rental_period,
    DROP COLUMN IF EXISTS living_area,
    DROP COLUMN IF EXISTS kitchen_area;

ALTER TABLE housing_complex
    DROP COLUMN IF EXISTS developer,
    DROP COLUMN IF EXISTS address,
    DROP COLUMN IF EXISTS starting_price,
    DROP COLUMN IF EXISTS image_urls;

ALTER TABLE users
    DROP COLUMN IF EXISTS first_name,
    DROP COLUMN IF EXISTS last_name,
    DROP COLUMN IF EXISTS phone,
    DROP COLUMN IF EXISTS photo_url;
