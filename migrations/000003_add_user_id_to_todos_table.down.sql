ALTER TABLE todos_table DROP CONSTRAINT fk_todos_user;

ALTER TABLE todos_table DROP COLUMN IF EXISTS user_id;