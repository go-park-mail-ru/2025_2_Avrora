CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS offer (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    location_id INT NOT NULL,
    category_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image TEXT,
    price INT NOT NULL CHECK (price >= 0),
    area DECIMAL(10, 2) CHECK (area > 0),
    rooms INT CHECK (rooms >= 0),
    address VARCHAR(500) NOT NULL,
    offer_type VARCHAR(20) NOT NULL DEFAULT 'sale' CHECK (offer_type IN ('sale', 'rent')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE offer
    ADD CONSTRAINT fk_offer_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
