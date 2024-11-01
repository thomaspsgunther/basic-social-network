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

   
-- -- Update follower count
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


-- -- Update followed count
-- CREATE OR REPLACE FUNCTION update_followed_counts()
--     RETURNS TRIGGER AS $$
--     BEGIN
--         UPDATE users u
--         SET followed_count = (
--             SELECT COUNT(*)
--             FROM followers f
--             WHERE f.follower_id = u.id
--         )
--         WHERE EXISTS (SELECT 1 FROM users WHERE id = u.id);

--         RETURN NULL;
--     END;
--     $$ LANGUAGE 'plpgsql';


-- -- Update post count
-- CREATE OR REPLACE FUNCTION update_post_counts()
--     RETURNS TRIGGER AS $$
--     BEGIN
--         UPDATE users u
--         SET post_count = (
--             SELECT COUNT(*)
--             FROM posts p
--             WHERE p.user_id = u.id
--         )
--         WHERE EXISTS (SELECT 1 FROM users WHERE id = u.id);

--         RETURN NULL;
--     END;
--     $$ LANGUAGE 'plpgsql';


-- -- Update like count
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


-- -- Update comment count
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
