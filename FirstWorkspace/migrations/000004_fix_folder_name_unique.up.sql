CREATE UNIQUE INDEX idx_user_folder_name ON folders (name, user_id) WHERE deleted_at IS NULL;
