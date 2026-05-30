CREATE INDEX IF NOT EXISTS idx_bookmarks_user_created
    ON bookmarks (user_id, created_at DESC)
    WHERE deleted_at IS NULL;
