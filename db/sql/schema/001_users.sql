CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL DEFAULT UNIQUE,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL  DEFAULT NOW(),
    password_hash TEXT NOT NULL,
    google_id TEXT UNIQUE,
    name TEXT,
    avatar_url TEXT
)