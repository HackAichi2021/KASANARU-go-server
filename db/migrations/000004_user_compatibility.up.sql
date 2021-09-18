BEGIN;
CREATE TABLE IF NOT EXISTS users_compatibility(
   id serial PRIMARY KEY,
   star INTEGER DEFAULT 0 NOT NULL,
   categories INTEGER[] NOT NULL,
);
COMMIT;
