BEGIN;
CREATE TABLE IF NOT EXISTS user_compatibility(
   id serial PRIMARY KEY,
   user_from_id VARCHAR (255) NOT NULL,
   user_to_id VARCHAR (255) NOT NULL,
   star VARCHAR (255) UNIQUE NOT NULL
);
COMMIT;
