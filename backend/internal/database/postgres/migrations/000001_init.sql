CREATE TABLE IF NOT EXISTS migrations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    file text,
    created_at timestamp DEFAULT (NOW() AT TIME ZONE 'utc')
);

-- -- Update modified column
-- CREATE OR REPLACE FUNCTION update_modified_column()
--    RETURNS TRIGGER AS $$
--    BEGIN
--        NEW.updated_at = (NOW() AT TIME ZONE 'utc');
--        RETURN NEW;
--    END;
--    $$ LANGUAGE 'plpgsql';

   
-- -- Update followers
-- CREATE OR REPLACE FUNCTION update_follower_counts()
--     RETURNS TRIGGER AS $$
--     BEGIN
--         UPDATE users u
--         SET follower_count = (
--             SELECT COUNT(*)
--             FROM followers f
--             WHERE f.followed_id = u.id
--         )
--         WHERE EXISTS (SELECT 1 FROM users WHERE id = u.id);

--         RETURN NULL;
--     END;
--     $$ LANGUAGE 'plpgsql';

-- -- Update likes
-- CREATE OR REPLACE FUNCTION update_like_counts()
--     RETURNS TRIGGER AS $$
--     BEGIN
--         UPDATE posts p
--         SET like_count = (
--             SELECT COUNT(*)
--             FROM likes l
--             WHERE l.post_id = p.id
--         )
--         WHERE EXISTS (SELECT 1 FROM posts WHERE id = p.id);

--         RETURN NULL;
--     END;
--     $$ LANGUAGE 'plpgsql';

-- -- Update comments
-- CREATE OR REPLACE FUNCTION update_comment_counts()
--     RETURNS TRIGGER AS $$
--     BEGIN
--         UPDATE posts p
--         SET comment_count = (
--             SELECT COUNT(*)
--             FROM comments c
--             WHERE c.post_id = p.id
--         )
--         WHERE EXISTS (SELECT 1 FROM posts WHERE id = p.id);

--         RETURN NULL;
--     END;
--     $$ LANGUAGE 'plpgsql';
