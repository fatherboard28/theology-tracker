CREATE TABLE study_sessions (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    date             TEXT    NOT NULL,              -- ISO 8601: YYYY-MM-DD
    duration_minutes INTEGER NOT NULL CHECK (duration_minutes > 0),
    -- Both associations are optional: a session can be logged free-form
    -- with no work item and no method.
    work_item_id     INTEGER REFERENCES work_items(id) ON DELETE SET NULL,
    method_id        INTEGER REFERENCES methods(id)    ON DELETE SET NULL,
    reflection       TEXT    NOT NULL DEFAULT '',
    created_at       TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_sessions_date        ON study_sessions(date);
CREATE INDEX idx_sessions_work_item   ON study_sessions(work_item_id);
CREATE INDEX idx_sessions_method      ON study_sessions(method_id);
