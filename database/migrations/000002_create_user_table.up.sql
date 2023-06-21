BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS image
(
    id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    bucket_name          text,
    path                 text,
    created_at           timestamp with time zone default now(),
    updated_at           timestamp with time zone,
    archived_at          timestamp with time zone
                                       );

COMMIT ;

CREATE TABLE IF NOT EXISTS users
(
    id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name                  text,
    address              text,
    email                 text,
    phone                text,
    dob                  text,
    image_id             uuid references image(id),
    created_at           timestamp with time zone default now(),
    updated_at           timestamp with time zone,
    archived_at          timestamp with time zone
                                       );

COMMIT ;