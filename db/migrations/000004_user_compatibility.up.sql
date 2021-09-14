BEGIN;
CREATE TABLE IF NOT EXISTS users_compatibility(
   id serial PRIMARY KEY,
   user_from_id VARCHAR (255) NOT NULL,
   user_to_id VARCHAR (255) NOT NULL,
   star VARCHAR (255) UNIQUE NOT NULL,
   updated_at BIGINT DEFAULT 0 NOT NULL,
   created_at BIGINT DEFAULT 0 NOT NULL,
   deleted_at BIGINT DEFAULT 0 NOT NULL
);
COMMIT;
