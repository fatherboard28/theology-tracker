package models

type Unit struct {
	ID                   int64   `db:"id"`
	CourseID             int64   `db:"course_id"`
	Title                string  `db:"title"`
	Description          string  `db:"description"`
	Position             int     `db:"position"`               // display order within the course
	TargetCompletionDate *string `db:"target_completion_date"` // YYYY-MM-DD or nil
	ActualCompletionDate *string `db:"actual_completion_date"` // set when all work items complete
	CreatedAt            string  `db:"created_at"`
}

// IsComplete returns true when all work items in the unit have been completed.
// The actual_completion_date is set by the application at that moment.
func (u *Unit) IsComplete() bool {
	return u.ActualCompletionDate != nil
}
