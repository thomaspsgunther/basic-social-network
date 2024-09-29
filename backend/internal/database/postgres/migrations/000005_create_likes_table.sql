CREATE TABLE IF NOT EXISTS likes (
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    post_id uuid REFERENCES posts(id) ON DELETE CASCADE,

    created_at timestamp DEFAULT (NOW() AT TIME ZONE 'utc'),

    PRIMARY KEY (post_id, user_id)
);
CREATE TRIGGER update_like_count_trigger AFTER INSERT OR DELETE ON likes FOR EACH STATEMENT EXECUTE FUNCTION update_like_counts();
