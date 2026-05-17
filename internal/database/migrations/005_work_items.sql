CREATE TABLE work_items (
    id                         INTEGER PRIMARY KEY AUTOINCREMENT,
    type                       TEXT    NOT NULL
                                   CHECK (type IN ('Reading', 'Assignment', 'Paper', 'Practice Session')),
    title                      TEXT    NOT NULL,
    status                     TEXT    NOT NULL DEFAULT 'Not Started'
                                   CHECK (status IN ('Not Started', 'In Progress', 'Complete')),
    estimated_duration_minutes INTEGER,
    due_date                   TEXT,
    -- set automatically by the application when status is changed to Complete
    completion_date            TEXT,
    general_notes              TEXT    NOT NULL DEFAULT '',
    -- nullable: a work item owned by a topic has no unit parent
    unit_id                    INTEGER REFERENCES units(id) ON DELETE CASCADE,
    -- nullable: a work item inside a unit may have no owning topic
    -- RESTRICT prevents deleting a topic that is the sole parent of work items
    owning_topic_id            INTEGER REFERENCES topics(id) ON DELETE RESTRICT,
    created_at                 TEXT    NOT NULL DEFAULT (datetime('now')),
    -- spec rule: every work item must have at least one of unit_id or owning_topic_id
    CHECK (unit_id IS NOT NULL OR owning_topic_id IS NOT NULL)
);

CREATE INDEX idx_work_items_unit     ON work_items(unit_id);
CREATE INDEX idx_work_items_topic    ON work_items(owning_topic_id);
CREATE INDEX idx_work_items_status   ON work_items(status);
CREATE INDEX idx_work_items_due_date ON work_items(due_date);
