CREATE TABLE IF NOT EXISTS comments (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    post_id uuid REFERENCES posts(id) ON DELETE CASCADE,
    description text NOT NULL,
    like_count int DEFAULT 0,

    created_at timestamp DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at timestamp DEFAULT (NOW() AT TIME ZONE 'utc')
);
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE TRIGGER update_modified_time BEFORE UPDATE ON comments FOR EACH ROW EXECUTE FUNCTION update_modified_column();
