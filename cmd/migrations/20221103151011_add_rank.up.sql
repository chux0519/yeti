CREATE TABLE maplegg (
    ign TEXT PRIMARY KEY,
    data JSON,
    avatar BLOB,

    created_at TEXT NOT NULL, -- ISO8601 utc time
    updated_at TEXT NOT NULL
);