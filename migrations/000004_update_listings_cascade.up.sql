ALTER TABLE listings 
    DROP CONSTRAINT IF EXISTS listings_user_id_fkey;

ALTER TABLE listings 
    ADD CONSTRAINT listings_user_id_fkey 
    FOREIGN KEY (user_id) 
    REFERENCES users(id) 
    ON DELETE CASCADE;
