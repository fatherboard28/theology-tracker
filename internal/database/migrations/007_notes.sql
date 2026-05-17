CREATE TABLE notes (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    title               TEXT    NOT NULL,
    body                TEXT    NOT NULL DEFAULT '',  -- markdown source
    -- Polymorphic parent: exactly one of Course / Unit / Topic / Method / Session.
    -- SQLite cannot enforce a cross-table FK on a polymorphic column, so
    -- referential integrity is handled at the application layer.
    primary_parent_type TEXT    NOT NULL
                            CHECK (primary_parent_type IN ('Course', 'Unit', 'Topic', 'Method', 'Session')),
    primary_parent_id   INTEGER NOT NULL,
    created_at          TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at          TEXT    NOT NULL DEFAULT (datetime('now'))
);

-- Auto-update updated_at whenever a note's body or title changes.
CREATE TRIGGER notes_updated_at
AFTER UPDATE ON notes
FOR EACH ROW
BEGIN
    UPDATE notes SET updated_at = datetime('now') WHERE id = OLD.id;
END;

CREATE INDEX idx_notes_parent ON notes(primary_parent_type, primary_parent_id);
CREATE INDEX idx_notes_updated ON notes(updated_at);
