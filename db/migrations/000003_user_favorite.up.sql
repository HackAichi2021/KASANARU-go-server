BEGIN;
CREATE TABLE IF NOT EXISTS user_favorite(
   user_id INTEGER,
   favorite_id INTEGER
);
COMMIT;
