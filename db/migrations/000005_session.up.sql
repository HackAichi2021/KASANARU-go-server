BEGIN;
CREATE TABLE IF NOT EXISTS sessions(
   access_uuid VARCHAR (255),
   refresh_uuid VARCHAR (255),
   user_id INTEGER
);
COMMIT;
