CREATE TABLE courses (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    title                  TEXT    NOT NULL,
    description            TEXT    NOT NULL DEFAULT '',
    status                 TEXT    NOT NULL DEFAULT 'Active'
                               CHECK (status IN ('Active', 'Paused', 'Complete')),
    start_date             TEXT,
    target_completion_date TEXT,
    actual_completion_date TEXT,
    created_at             TEXT    NOT NULL DEFAULT (datetime('now'))
);
