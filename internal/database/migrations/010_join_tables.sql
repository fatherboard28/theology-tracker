-- ── Topic tagging join tables ─────────────────────────────────────────────────
-- Any entity can be tagged with any number of topics.
-- The composite PK prevents duplicate tags and doubles as an index.

CREATE TABLE course_topics (
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    topic_id  INTEGER NOT NULL REFERENCES topics(id)  ON DELETE CASCADE,
    PRIMARY KEY (course_id, topic_id)
);
CREATE INDEX idx_course_topics_topic ON course_topics(topic_id);

CREATE TABLE unit_topics (
    unit_id  INTEGER NOT NULL REFERENCES units(id)  ON DELETE CASCADE,
    topic_id INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    PRIMARY KEY (unit_id, topic_id)
);
CREATE INDEX idx_unit_topics_topic ON unit_topics(topic_id);

CREATE TABLE work_item_topics (
    work_item_id INTEGER NOT NULL REFERENCES work_items(id) ON DELETE CASCADE,
    topic_id     INTEGER NOT NULL REFERENCES topics(id)     ON DELETE CASCADE,
    PRIMARY KEY (work_item_id, topic_id)
);
CREATE INDEX idx_work_item_topics_topic ON work_item_topics(topic_id);

CREATE TABLE note_topics (
    note_id  INTEGER NOT NULL REFERENCES notes(id)  ON DELETE CASCADE,
    topic_id INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    PRIMARY KEY (note_id, topic_id)
);
CREATE INDEX idx_note_topics_topic ON note_topics(topic_id);

CREATE TABLE session_topics (
    session_id INTEGER NOT NULL REFERENCES study_sessions(id) ON DELETE CASCADE,
    topic_id   INTEGER NOT NULL REFERENCES topics(id)         ON DELETE CASCADE,
    PRIMARY KEY (session_id, topic_id)
);
CREATE INDEX idx_session_topics_topic ON session_topics(topic_id);

-- ── Note ↔ Work Item references (many-to-many) ────────────────────────────────
-- Assignment, Paper, and Practice Session work items can reference any number
-- of notes. A note can be referenced by any number of work items.
-- This is distinct from topic tagging — it's a "this note is source material
-- or output for this work item" relationship.

CREATE TABLE note_work_items (
    note_id      INTEGER NOT NULL REFERENCES notes(id)      ON DELETE CASCADE,
    work_item_id INTEGER NOT NULL REFERENCES work_items(id) ON DELETE CASCADE,
    PRIMARY KEY (note_id, work_item_id)
);
CREATE INDEX idx_note_work_items_work_item ON note_work_items(work_item_id);
