DROP TABLE IF EXISTS user_compatibility;
BEGIN;
CREATE TABLE IF NOT EXISTS session(
   id serial PRIMARY KEY,
   user_id INTEGER NOT NULL,
   session_id VARCHAR (255),
   expired BIGINT DEFAULT 0 NOT NULL
);
COMMIT;
