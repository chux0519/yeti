CREATE TABLE maplegg (
    ign TEXT PRIMARY KEY,
    data JSON,
    profile_img BLOB,

    created_at TEXT NOT NULL, -- ISO8601 utc time
    updated_at TEXT NOT NULL
);