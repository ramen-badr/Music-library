CREATE TABLE IF NOT EXISTS songs (
     id SERIAL PRIMARY KEY,
     "group" TEXT NOT NULL,
     song TEXT NOT NULL,
     release_date TEXT,
     text TEXT,
     link TEXT
);