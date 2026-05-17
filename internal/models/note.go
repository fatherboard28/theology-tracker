package models

type NoteParentType string

const (
	NoteParentCourse  NoteParentType = "Course"
	NoteParentUnit    NoteParentType = "Unit"
	NoteParentTopic   NoteParentType = "Topic"
	NoteParentMethod  NoteParentType = "Method"
	NoteParentSession NoteParentType = "Session"
)

type Note struct {
	ID               int64          `db:"id"`
	Title            string         `db:"title"`
	Body             string         `db:"body"` // markdown source
	PrimaryParentType NoteParentType `db:"primary_parent_type"`
	PrimaryParentID  int64          `db:"primary_parent_id"`
	CreatedAt        string         `db:"created_at"`
	UpdatedAt        string         `db:"updated_at"`
}

// NoteSummary is used in listing views — avoids loading the full markdown body.
type NoteSummary struct {
	ID                int64          `db:"id"`
	Title             string         `db:"title"`
	PrimaryParentType NoteParentType `db:"primary_parent_type"`
	PrimaryParentID   int64          `db:"primary_parent_id"`
	UpdatedAt         string         `db:"updated_at"`
	// Populated by JOIN queries, not DB columns:
	TopicTags   []TopicSummary
	BodyPreview string // first ~120 chars of body, stripped of markdown
}
