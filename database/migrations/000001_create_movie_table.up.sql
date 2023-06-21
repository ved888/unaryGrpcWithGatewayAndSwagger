BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS movie
(
    id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title                text,
    genre               text,
    created_at           timestamp with time zone default now(),
    updated_at           timestamp with time zone,
    archived_at          timestamp with time zone
);

COMMIT ;