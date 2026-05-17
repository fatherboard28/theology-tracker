package models

type StudySession struct {
	ID              int64  `db:"id"`
	Date            string `db:"date"`             // YYYY-MM-DD
	DurationMinutes int    `db:"duration_minutes"`
	WorkItemID      *int64 `db:"work_item_id"` // nil for free-form sessions
	MethodID        *int64 `db:"method_id"`    // nil if no method used
	Reflection      string `db:"reflection"`
	CreatedAt       string `db:"created_at"`
}

// StudySessionWithContext enriches a session with joined display fields
// used in listings and the calendar view.
type StudySessionWithContext struct {
	StudySession
	WorkItemTitle *string `db:"work_item_title"` // nil if no work item
	MethodName    *string `db:"method_name"`     // nil if no method
	TopicTags     []TopicSummary
	ScriptureTags []ScriptureTag
}

// StreakResult holds the current and all-time streak counts,
// computed by the store from session dates.
type StreakResult struct {
	Current int
	AllTime int
}
