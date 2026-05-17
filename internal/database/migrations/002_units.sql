CREATE TABLE units (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    course_id              INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title                  TEXT    NOT NULL,
    description            TEXT    NOT NULL DEFAULT '',
    -- position drives display order within a course; updated by drag-and-drop
    position               INTEGER NOT NULL DEFAULT 0,
    target_completion_date TEXT,
    -- set automatically by the application when all work items in the unit complete
    actual_completion_date TEXT,
    created_at             TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_units_course_id ON units(course_id);
