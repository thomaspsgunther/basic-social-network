CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    username text NOT NULL UNIQUE,
    password text NOT NULL,
    email text UNIQUE,
    full_name text,
    description text,
    avatar text,
    follower_count int DEFAULT 0,

    created_at timestamp DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at timestamp DEFAULT (NOW() AT TIME ZONE 'utc')
);
CREATE INDEX idx_users_full_name ON users(full_name);
CREATE TRIGGER update_modified_time BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_modified_column();
