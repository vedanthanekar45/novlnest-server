CREATE TABLE shelves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userid UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE shelf_books (
    shelf_id UUID NOT NULL REFERENCES shelves(id) ON DELETE CASCADE,
    isbn13 TEXT NOT NULL,
    google_id TEXT NOT NULL,
    title TEXT NOT NULL,
    thumbnail_url TEXT,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (shelf_id, google_id)
);

CREATE TABLE book_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userid UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    isbn13 TEXT NOT NULL,
    google_id TEXT NOT NULL,
    title TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,

    status TEXT NOT NULL CHECK (status IN ('to_be_read', 'currently_reading', 'completed', 'dnf')),
    rating REAL CHECK (rating IS NULL OR (rating >= 0 AND rating <= 5 AND (rating * 2.0) = FLOOR(rating * 2.0))),
    review TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(userid, google_id)
);