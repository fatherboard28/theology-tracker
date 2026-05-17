package models

type WorkItemType string

const (
	WorkItemTypeReading         WorkItemType = "Reading"
	WorkItemTypeAssignment      WorkItemType = "Assignment"
	WorkItemTypePaper           WorkItemType = "Paper"
	WorkItemTypePracticeSession WorkItemType = "Practice Session"
)

type WorkItemStatus string

const (
	WorkItemStatusNotStarted WorkItemStatus = "Not Started"
	WorkItemStatusInProgress WorkItemStatus = "In Progress"
	WorkItemStatusComplete   WorkItemStatus = "Complete"
)

type ReadingFormat string

const (
	ReadingFormatPhysicalBook  ReadingFormat = "Physical Book"
	ReadingFormatPDF           ReadingFormat = "PDF"
	ReadingFormatOnlineArticle ReadingFormat = "Online Article"
	ReadingFormatScripture     ReadingFormat = "Scripture"
)

// WorkItem holds the fields common to all work item types.
type WorkItem struct {
	ID                       int64          `db:"id"`
	Type                     WorkItemType   `db:"type"`
	Title                    string         `db:"title"`
	Status                   WorkItemStatus `db:"status"`
	EstimatedDurationMinutes *int           `db:"estimated_duration_minutes"`
	DueDate                  *string        `db:"due_date"`        // YYYY-MM-DD or nil
	CompletionDate           *string        `db:"completion_date"` // set automatically
	GeneralNotes             string         `db:"general_notes"`
	UnitID                   *int64         `db:"unit_id"`        // nil if topic-owned
	OwningTopicID            *int64         `db:"owning_topic_id"` // nil if unit-owned
	CreatedAt                string         `db:"created_at"`
}

func (w *WorkItem) IsComplete() bool {
	return w.Status == WorkItemStatusComplete
}

// Reading holds the type-specific fields for WorkItemTypeReading.
type Reading struct {
	WorkItemID int64         `db:"work_item_id"`
	Source     string        `db:"source"` // book title, article name, or "Scripture"
	Author     string        `db:"author"`
	Location   string        `db:"location"` // page range or scripture reference
	Format     ReadingFormat `db:"format"`
}

// Assignment holds the type-specific fields for WorkItemTypeAssignment.
type Assignment struct {
	WorkItemID  int64  `db:"work_item_id"`
	Description string `db:"description"`
}

// Paper holds the type-specific fields for WorkItemTypePaper.
type Paper struct {
	WorkItemID     int64  `db:"work_item_id"`
	PromptOrTopic  string `db:"prompt_or_topic"`
	WordCountTarget *int  `db:"word_count_target"`
	ScoreOrGrade   string `db:"score_or_grade"` // freeform: letter, %, or narrative
}

// PracticeSession holds the type-specific fields for WorkItemTypePracticeSession.
type PracticeSession struct {
	WorkItemID       int64  `db:"work_item_id"`
	MethodID         *int64 `db:"method_id"` // nil if no method used
	ScripturePassage string `db:"scripture_passage"`
	DurationMinutes  *int   `db:"duration_minutes"`
}

// WorkItemWithDetails bundles a WorkItem with its type-specific detail
// and derived fields used in list views. The Detail field holds one of
// *Reading, *Assignment, *Paper, or *PracticeSession depending on Type.
type WorkItemWithDetails struct {
	WorkItem
	Detail          any    // type-specific struct
	TotalTimeLogged int    // derived from associated study sessions (minutes)
	TopicTags       []TopicSummary
	ScriptureTags   []ScriptureTag
}
