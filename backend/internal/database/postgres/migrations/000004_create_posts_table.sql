CREATE TABLE IF NOT EXISTS posts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    image text NOT NULL,
    description text,
    like_count int DEFAULT 0,
    comment_count int DEFAULT 0,

    created_at timestamp DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at timestamp DEFAULT (NOW() AT TIME ZONE 'utc')
);
CREATE INDEX IF NOT EXISTS idx_posts_pagination ON posts (created_at, id);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE TRIGGER update_post_count_trigger AFTER INSERT OR DELETE ON posts FOR EACH STATEMENT EXECUTE FUNCTION update_post_counts();
CREATE TRIGGER update_modified_time BEFORE UPDATE ON posts FOR EACH ROW EXECUTE FUNCTION update_modified_column();
