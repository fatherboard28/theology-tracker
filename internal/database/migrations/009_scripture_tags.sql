CREATE TABLE scripture_tags (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Validated by the application against the known book abbreviation list.
    -- Format: "[3-char book] [chapter]:[verse]" e.g. "Rom 8:28", "Gen 1:1-5"
    reference   TEXT    NOT NULL,
    -- Polymorphic: the entity this tag is attached to.
    -- Allowed types: WorkItem, Topic, Session
    entity_type TEXT    NOT NULL
                    CHECK (entity_type IN ('WorkItem', 'Topic', 'Session')),
    entity_id   INTEGER NOT NULL
);

CREATE INDEX idx_scripture_entity ON scripture_tags(entity_type, entity_id);
-- Supports the reference view: "show all items tagged with Rom 8"
CREATE INDEX idx_scripture_ref    ON scripture_tags(reference);
