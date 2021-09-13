BEGIN;
CREATE TABLE IF NOT EXISTS favorite(
   user_id serial PRIMARY KEY,
   content VARCHAR (255) NOT NULL
);
COMMIT;
