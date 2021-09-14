BEGIN;
create table if not exists users (
   id BIGSERIAL PRIMARY KEY,
   user_name VARCHAR (255) NOT NULL,
   password VARCHAR(255) NOT NULL,
   email VARCHAR (255) UNIQUE NOT NULL,
   age INTEGER NOT NULL,
   categories INTEGER[],
   updated_at BIGINT DEFAULT 0 NOT NULL,
   created_at BIGINT DEFAULT 0 NOT NULL,
   deleted_at BIGINT DEFAULT 0 NOT NULL
);
COMMIT;
