-- create user profile table
CREATE TABLE IF NOT EXISTS user_profile (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(200),
    user_id INT NOT NULL UNIQUE REFERENCES users(id),
    current_job VARCHAR(255),
    experience_level VARCHAR(100),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- create foreign key constraint to users table
ALTER TABLE user_profile
ADD CONSTRAINT fk_user
FOREIGN KEY (user_id) REFERENCES users(id)
ON DELETE CASCADE;