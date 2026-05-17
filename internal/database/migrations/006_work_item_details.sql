-- Each table has a 1-to-1 relationship with work_items via the PK.
-- A row exists here only when work_items.type matches the table's type.
-- The application enforces consistency at creation time.

CREATE TABLE readings (
    work_item_id INTEGER PRIMARY KEY REFERENCES work_items(id) ON DELETE CASCADE,
    source       TEXT    NOT NULL DEFAULT '',   -- book title, article name, or 'Scripture'
    author       TEXT    NOT NULL DEFAULT '',
    location     TEXT    NOT NULL DEFAULT '',   -- page range or scripture reference
    format       TEXT    NOT NULL DEFAULT 'Physical Book'
                     CHECK (format IN ('Physical Book', 'PDF', 'Online Article', 'Scripture'))
);

CREATE TABLE assignments (
    work_item_id INTEGER PRIMARY KEY REFERENCES work_items(id) ON DELETE CASCADE,
    description  TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE papers (
    work_item_id      INTEGER PRIMARY KEY REFERENCES work_items(id) ON DELETE CASCADE,
    prompt_or_topic   TEXT    NOT NULL DEFAULT '',
    word_count_target INTEGER,
    score_or_grade    TEXT    NOT NULL DEFAULT ''   -- freeform: letter, %, or narrative
);

CREATE TABLE practice_sessions (
    work_item_id      INTEGER PRIMARY KEY REFERENCES work_items(id) ON DELETE CASCADE,
    method_id         INTEGER REFERENCES methods(id) ON DELETE SET NULL,
    scripture_passage TEXT    NOT NULL DEFAULT '',
    duration_minutes  INTEGER
);
