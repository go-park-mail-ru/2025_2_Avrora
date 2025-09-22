CREATE TABLE photo (
    id SERIAL PRIMARY KEY,
    offer_id INT NOT NULL REFERENCES offer(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    position INT NOT NULL DEFAULT 0,
);

CREATE INDEX idx_photo_offer_id ON photo (offer_id);