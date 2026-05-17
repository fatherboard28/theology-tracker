package store

import (
	"fmt"
	"strings"
	"time"

	"theology-tracker/internal/models"
)

// ── Param structs ─────────────────────────────────────────────────────────────

type CreateCourseParams struct {
	Title                string
	Description          string
	Status               models.CourseStatus
	StartDate            *string
	TargetCompletionDate *string
}

type UpdateCourseParams struct {
	ID                   int64
	Title                string
	Description          string
	Status               models.CourseStatus
	StartDate            *string
	TargetCompletionDate *string
}

type CreateUnitParams struct {
	CourseID             int64
	Title                string
	Description          string
	TargetCompletionDate *string
}

type UpdateUnitParams struct {
	ID                   int64
	Title                string
	Description          string
	TargetCompletionDate *string
}

// ── Courses ───────────────────────────────────────────────────────────────────

func (s *Store) ListCourses() ([]models.CourseWithProgress, error) {
	type row struct {
		models.Course
		TotalWorkItems     int `db:"total_work_items"`
		CompletedWorkItems int `db:"completed_work_items"`
	}

	const q = `
		SELECT
			c.*,
			COUNT(wi.id)                                                   AS total_work_items,
			COALESCE(SUM(CASE WHEN wi.status = 'Complete' THEN 1 ELSE 0 END), 0) AS completed_work_items
		FROM courses c
		LEFT JOIN units u      ON u.course_id = c.id
		LEFT JOIN work_items wi ON wi.unit_id  = u.id
		GROUP BY c.id
		ORDER BY c.created_at DESC`

	var rows []row
	if err := s.db.Select(&rows, q); err != nil {
		return nil, fmt.Errorf("listing courses: %w", err)
	}

	out := make([]models.CourseWithProgress, len(rows))
	for i, r := range rows {
		out[i] = models.CourseWithProgress{
			Course:             r.Course,
			TotalWorkItems:     r.TotalWorkItems,
			CompletedWorkItems: r.CompletedWorkItems,
		}
	}
	return out, nil
}

func (s *Store) GetCourse(id int64) (*models.CourseWithProgress, error) {
	var c models.CourseWithProgress
	if err := s.db.Get(&c, `
		SELECT
			c.*,
			COALESCE(COUNT(wi.id), 0)                                              AS total_work_items,
			COALESCE(SUM(CASE WHEN wi.status = 'Complete' THEN 1 ELSE 0 END), 0)   AS completed_work_items
		FROM courses c
		LEFT JOIN units      u  ON u.course_id = c.id
		LEFT JOIN work_items wi ON wi.unit_id  = u.id
		WHERE c.id = ?
		GROUP BY c.id
	`, id); err != nil {
		return nil, fmt.Errorf("getting course %d: %w", id, err)
	}
	return &c, nil
}

func (s *Store) CreateCourse(p CreateCourseParams) (int64, error) {
	res, err := s.db.Exec(`
		INSERT INTO courses (title, description, status, start_date, target_completion_date)
		VALUES (?, ?, ?, ?, ?)`,
		p.Title, p.Description, p.Status, p.StartDate, p.TargetCompletionDate,
	)
	if err != nil {
		return 0, fmt.Errorf("creating course: %w", err)
	}
	return res.LastInsertId()
}

func (s *Store) UpdateCourse(p UpdateCourseParams) error {
	_, err := s.db.Exec(`
		UPDATE courses
		SET title = ?, description = ?, status = ?, start_date = ?, target_completion_date = ?
		WHERE id = ?`,
		p.Title, p.Description, p.Status, p.StartDate, p.TargetCompletionDate, p.ID,
	)
	if err != nil {
		return fmt.Errorf("updating course %d: %w", p.ID, err)
	}
	return nil
}

func (s *Store) DeleteCourse(id int64) error {
	_, err := s.db.Exec(`DELETE FROM courses WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("deleting course %d: %w", id, err)
	}
	return nil
}

// SetCourseStatus updates the status and manages actual_completion_date:
// setting Complete records today's date; any other status clears it.
func (s *Store) SetCourseStatus(id int64, status models.CourseStatus) error {
	var completionDate *string
	if status == models.CourseStatusComplete {
		d := time.Now().Format("2006-01-02")
		completionDate = &d
	}
	_, err := s.db.Exec(
		`UPDATE courses SET status = ?, actual_completion_date = ? WHERE id = ?`,
		status, completionDate, id,
	)
	if err != nil {
		return fmt.Errorf("setting course status: %w", err)
	}
	return nil
}

// ── Units ─────────────────────────────────────────────────────────────────────

func (s *Store) ListUnitsForCourse(courseID int64) ([]models.UnitWithProgress, error) {
	type row struct {
		models.Unit
		TotalWorkItems     int `db:"total_work_items"`
		CompletedWorkItems int `db:"completed_work_items"`
	}

	const q = `
		SELECT
			u.*,
			COUNT(wi.id)                                                    AS total_work_items,
			COALESCE(SUM(CASE WHEN wi.status = 'Complete' THEN 1 ELSE 0 END), 0) AS completed_work_items
		FROM units u
		LEFT JOIN work_items wi ON wi.unit_id = u.id
		WHERE u.course_id = ?
		GROUP BY u.id
		ORDER BY u.position ASC, u.created_at ASC`

	var rows []row
	if err := s.db.Select(&rows, q, courseID); err != nil {
		return nil, fmt.Errorf("listing units for course %d: %w", courseID, err)
	}

	out := make([]models.UnitWithProgress, len(rows))
	for i, r := range rows {
		out[i] = models.UnitWithProgress{
			Unit:               r.Unit,
			TotalWorkItems:     r.TotalWorkItems,
			CompletedWorkItems: r.CompletedWorkItems,
		}
	}
	return out, nil
}

func (s *Store) GetUnit(id int64) (*models.Unit, error) {
	var u models.Unit
	if err := s.db.Get(&u, `SELECT * FROM units WHERE id = ?`, id); err != nil {
		return nil, fmt.Errorf("getting unit %d: %w", id, err)
	}
	return &u, nil
}

func (s *Store) CreateUnit(p CreateUnitParams) (int64, error) {
	// Place the new unit at the end of the course's unit list.
	var maxPos int
	_ = s.db.QueryRow(
		`SELECT COALESCE(MAX(position), -1) FROM units WHERE course_id = ?`, p.CourseID,
	).Scan(&maxPos)

	res, err := s.db.Exec(`
		INSERT INTO units (course_id, title, description, position, target_completion_date)
		VALUES (?, ?, ?, ?, ?)`,
		p.CourseID, p.Title, p.Description, maxPos+1, p.TargetCompletionDate,
	)
	if err != nil {
		return 0, fmt.Errorf("creating unit: %w", err)
	}
	return res.LastInsertId()
}

func (s *Store) UpdateUnit(p UpdateUnitParams) error {
	_, err := s.db.Exec(`
		UPDATE units
		SET title = ?, description = ?, target_completion_date = ?
		WHERE id = ?`,
		p.Title, p.Description, p.TargetCompletionDate, p.ID,
	)
	if err != nil {
		return fmt.Errorf("updating unit %d: %w", p.ID, err)
	}
	return nil
}

func (s *Store) DeleteUnit(id int64) error {
	_, err := s.db.Exec(`DELETE FROM units WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("deleting unit %d: %w", id, err)
	}
	return nil
}

// ReorderUnits sets each unit's position to its index in orderedIDs.
// It validates that all IDs belong to courseID before committing.
func (s *Store) ReorderUnits(courseID int64, orderedIDs []int64) error {
	if len(orderedIDs) == 0 {
		return nil
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	for pos, id := range orderedIDs {
		res, err := tx.Exec(
			`UPDATE units SET position = ? WHERE id = ? AND course_id = ?`,
			pos, id, courseID,
		)
		if err != nil {
			return fmt.Errorf("reordering unit %d: %w", id, err)
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return fmt.Errorf("unit %d does not belong to course %d", id, courseID)
		}
	}

	return tx.Commit()
}

// CheckAndCompleteUnit evaluates whether all work items in a unit are done.
// If so, it stamps actual_completion_date; if not, it clears it.
// Call this whenever a work item's status changes.
func (s *Store) CheckAndCompleteUnit(unitID int64) error {
	var total, completed int
	err := s.db.QueryRow(`
		SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'Complete' THEN 1 ELSE 0 END), 0)
		FROM work_items
		WHERE unit_id = ?`, unitID,
	).Scan(&total, &completed)
	if err != nil {
		return fmt.Errorf("checking unit completion: %w", err)
	}

	if total > 0 && total == completed {
		_, err = s.db.Exec(`
			UPDATE units
			SET actual_completion_date = date('now')
			WHERE id = ? AND actual_completion_date IS NULL`, unitID,
		)
	} else {
		_, err = s.db.Exec(`
			UPDATE units SET actual_completion_date = NULL WHERE id = ?`, unitID,
		)
	}
	if err != nil {
		return fmt.Errorf("updating unit completion date: %w", err)
	}
	return nil
}

// NilIfEmpty converts an empty string to a nil *string — used when
// optional date fields from HTML forms arrive as empty strings.
func NilIfEmpty(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}
