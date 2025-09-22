CREATE TABLE offer (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
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

CREATE INDEX idx_offer_user_id ON offer (user_id);
CREATE INDEX idx_offer_location_id ON offer (location_id);
CREATE INDEX idx_offer_category_id ON offer (category_id);
CREATE INDEX idx_offer_price ON offer (price);
CREATE INDEX idx_offer_created_at ON offer (created_at DESC);