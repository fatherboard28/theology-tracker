package models

type CourseStatus string

const (
	CourseStatusActive   CourseStatus = "Active"
	CourseStatusPaused   CourseStatus = "Paused"
	CourseStatusComplete CourseStatus = "Complete"
)

type Course struct {
	ID                   int64        `db:"id"`
	Title                string       `db:"title"`
	Description          string       `db:"description"`
	Status               CourseStatus `db:"status"`
	StartDate            *string      `db:"start_date"`             // YYYY-MM-DD or nil
	TargetCompletionDate *string      `db:"target_completion_date"` // YYYY-MM-DD or nil
	ActualCompletionDate *string      `db:"actual_completion_date"` // set on completion
	CreatedAt            string       `db:"created_at"`
}

// IsComplete returns true when the course has been marked complete.
func (c *Course) IsComplete() bool {
	return c.Status == CourseStatusComplete
}
