CREATE TABLE IF NOT EXISTS recipes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    url TEXT,
    ingredients TEXT,
    steps TEXT,
    summary TEXT,
    summary_embedding BLOB
);