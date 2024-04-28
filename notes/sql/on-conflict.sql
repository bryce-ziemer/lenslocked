INSERT INTO sessions (user_id, token_hash)
VALUES (1, 'xyz-456') ON CONFLICT (user_id) DO
UPDATE
SET token_hash = 'xyz-456';