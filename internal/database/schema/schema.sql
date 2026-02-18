-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(20),
    password_hash VARCHAR(255), -- Made nullable
    linkedin_id TEXT, -- Added new column
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on email for faster lookups
CREATE INDEX idx_users_email ON users(email);

-- Create index on active users
CREATE INDEX idx_users_active ON users(is_active) WHERE is_active = TRUE;

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create refresh tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on user_id for faster lookups
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- Create index on token for faster lookups
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);

-- Create index on expires_at for cleanup queries
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- create user profile table
CREATE TABLE IF NOT EXISTS user_profile (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE REFERENCES users(id),
    cv BYTEA,
    full_name VARCHAR(200),
    current_job VARCHAR(255),
    experience_level VARCHAR(100),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

------- JOB TABLE -------
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    company_name VARCHAR(150),
    notes TEXT,
    link TEXT,
    location VARCHAR(255),
    platform VARCHAR(100),
    date_applied TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'applied',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


-- ============================================================
-- Home Screen Features Migration
-- Focus Sessions, Streaks, Badges, Progress Stats
-- ============================================================

-- Focus sessions table (Lock-in sessions)
CREATE TABLE IF NOT EXISTS focus_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    duration_seconds INT NOT NULL DEFAULT 0,   -- total planned duration
    elapsed_seconds INT NOT NULL DEFAULT 0,    -- how much was completed
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active | paused | completed | abandoned
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_focus_sessions_user_id ON focus_sessions(user_id);
CREATE INDEX idx_focus_sessions_status ON focus_sessions(status);
CREATE INDEX idx_focus_sessions_started_at ON focus_sessions(started_at);

CREATE TRIGGER update_focus_sessions_updated_at BEFORE UPDATE ON focus_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Streaks table
CREATE TABLE IF NOT EXISTS streaks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    current_streak INT NOT NULL DEFAULT 0,
    longest_streak INT NOT NULL DEFAULT 0,
    last_session_date DATE,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_streaks_user_id ON streaks(user_id);

-- Badge definitions (static catalog)
CREATE TABLE IF NOT EXISTS badge_definitions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    icon_key VARCHAR(100),        -- frontend icon key / emoji slug
    criteria_type VARCHAR(50) NOT NULL, -- streak | sessions | focus_time | jobs
    criteria_value INT NOT NULL,        -- threshold value to unlock
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- User earned badges
CREATE TABLE IF NOT EXISTS user_badges (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    badge_id INT NOT NULL REFERENCES badge_definitions(id) ON DELETE CASCADE,
    earned_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, badge_id)
);

CREATE INDEX idx_user_badges_user_id ON user_badges(user_id);

-- Progress stats (aggregated per day for fast reads)
CREATE TABLE IF NOT EXISTS daily_progress (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    total_focus_seconds INT NOT NULL DEFAULT 0,
    sessions_completed INT NOT NULL DEFAULT 0,
    UNIQUE (user_id, date)
);

CREATE INDEX idx_daily_progress_user_id ON daily_progress(user_id);
CREATE INDEX idx_daily_progress_date ON daily_progress(date);

-- ============================================================
-- Seed badge definitions
-- ============================================================
INSERT INTO badge_definitions (name, description, icon_key, criteria_type, criteria_value) VALUES
    ('First Lock-in',   'Complete your first focus session',          'first_session',  'sessions',    1),
    ('3-Day Streak',    'Maintain a 3-day streak',                    'streak_3',       'streak',      3),
    ('7-Day Streak',    'Maintain a 7-day streak',                    'streak_7',       'streak',      7),
    ('30-Day Streak',   'Maintain a 30-day streak',                   'streak_30',      'streak',     30),
    ('Focus Rookie',    'Accumulate 1 hour of total focus time',      'focus_1h',       'focus_time', 3600),
    ('Focus Pro',       'Accumulate 10 hours of total focus time',    'focus_10h',      'focus_time', 36000),
    ('Job Hunter',      'Track 5 job applications',                   'job_hunter',     'jobs',        5),
    ('Consistency',     'Complete 10 focus sessions',                 'sessions_10',    'sessions',   10)
ON CONFLICT (name) DO NOTHING;