ALTER TABLE offer
    ADD CONSTRAINT fk_offer_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
