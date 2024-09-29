CREATE TABLE IF NOT EXISTS followers (
    follower_id uuid REFERENCES users(id) ON DELETE CASCADE,
    followed_id uuid REFERENCES users(id) ON DELETE CASCADE,

    created_at timestamp DEFAULT (NOW() AT TIME ZONE 'utc'),
    
    PRIMARY KEY (follower_id, followed_id)
);
CREATE TRIGGER update_follower_count_trigger AFTER INSERT OR DELETE ON followers FOR EACH STATEMENT EXECUTE FUNCTION update_follower_counts();
