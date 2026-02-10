ALTER TABLE todos_table ADD COLUMN user_id UUID NOT NULL;

ALTER TABLE todos_table ADD CONSTRAINT fk_todos_user FOREIGN KEY (user_id) REFERENCES users_table(id) ON DELETE CASCADE;