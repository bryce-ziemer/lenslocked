CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users (id) ON DELETE CASCADE, 
    token_hash TEXT UNIQUE NOT NULL
);


ALTER TABLE sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users
    (id);

DELETE FROM sessions
WHERE token_hash = $1;

SELECT users.id, users.email, users.password_hash
FROM sessions
    JOIN users on users.id = sessions.user_id
WHERE sessions.token_hash = $1;