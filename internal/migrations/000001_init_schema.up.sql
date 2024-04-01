CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- Create the index only if it doesn't already exist
CREATE INDEX IF NOT EXISTS idx_users_id ON users (id);

-- Create a view for the users table
CREATE OR REPLACE VIEW users_view AS
SELECT id, name
FROM users;
