BEGIN;
CREATE TABLE IF NOT EXISTS user_favorite(
   user_id INTEGER,
   favorite_id INTEGER,
   updated_at BIGINT DEFAULT 0 NOT NULL,
   created_at BIGINT DEFAULT 0 NOT NULL,
   deleted_at BIGINT DEFAULT 0 NOT NULL
);
COMMIT;
