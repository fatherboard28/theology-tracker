package models

// CourseWithProgress extends Course with work item completion counts
// derived from a JOIN query. Used in list and detail views.
type CourseWithProgress struct {
	Course
	TotalWorkItems     int `db:"total_work_items"`
	CompletedWorkItems int `db:"completed_work_items"`
}

func (c *CourseWithProgress) ProgressPercent() int {
	if c.TotalWorkItems == 0 {
		return 0
	}
	return (c.CompletedWorkItems * 100) / c.TotalWorkItems
}

// UnitWithProgress extends Unit with work item completion counts.
type UnitWithProgress struct {
	Unit
	TotalWorkItems     int `db:"total_work_items"`
	CompletedWorkItems int `db:"completed_work_items"`
}

func (u *UnitWithProgress) ProgressPercent() int {
	if u.TotalWorkItems == 0 {
		return 0
	}
	return (u.CompletedWorkItems * 100) / u.TotalWorkItems
}
