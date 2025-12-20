CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    company_name VARCHAR(150),
    notes TEXT,
    location VARCHAR(255),
    platform VARCHAR(100),
    date_applied TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'applied',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);