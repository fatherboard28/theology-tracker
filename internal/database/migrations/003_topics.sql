CREATE TABLE topics (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    title           TEXT    NOT NULL,
    description     TEXT    NOT NULL DEFAULT '',
    type            TEXT    NOT NULL DEFAULT 'Other'
                        CHECK (type IN ('Book of Bible', 'Theological Theme', 'Doctrine', 'Other')),
    -- one level of nesting only; enforced at the application layer
    parent_topic_id INTEGER REFERENCES topics(id) ON DELETE SET NULL,
    created_at      TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_topics_parent ON topics(parent_topic_id);
