CREATE TABLE IF NOT EXISTS users (
    user_id BIGINT PRIMARY KEY NOT NULL,
    username VARCHAR(33) NULL,
    fullname VARCHAR(129) NULL,
    register_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS admins (
    user_id BIGINT PRIMARY KEY NOT NULL REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TYPE file_type_enum AS ENUM (
  'photo',
  'document',
  'video',
  'audio',
  'voice',
  'sticker',
  'animation'
);

CREATE TABLE IF NOT EXISTS files (
  file_id     TEXT PRIMARY KEY,
  file_key    VARCHAR(50) NOT NULL UNIQUE,
  file_type   file_type_enum NOT NULL,
  uploaded_by BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);