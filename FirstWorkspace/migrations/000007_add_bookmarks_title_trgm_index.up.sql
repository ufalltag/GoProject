CREATE INDEX IF NOT EXISTS idx_bookmarks_title_btree ON bookmarks (user_id, lower(title)) WHERE deleted_at IS NULL;
