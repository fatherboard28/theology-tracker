package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"theology-tracker/internal/models"
	"theology-tracker/internal/store"
	coursepages "theology-tracker/internal/templates/pages/courses"
)

// ── Helpers ───────────────────────────────────────────────────────────────────

func courseIDParam(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "courseID"), 10, 64)
}

func unitIDParam(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "unitID"), 10, 64)
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ── Course handlers ───────────────────────────────────────────────────────────

func (h *Handler) ListCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.store.ListCourses()
	if err != nil {
		http.Error(w, "could not load courses", http.StatusInternalServerError)
		return
	}
	coursepages.List(courses).Render(r.Context(), w)
}

func (h *Handler) NewCourseForm(w http.ResponseWriter, r *http.Request) {
	coursepages.Form(coursepages.CourseFormData{IsNew: true}).Render(r.Context(), w)
}

func (h *Handler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	data := coursepages.CourseFormData{
		IsNew:       true,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Status:      r.FormValue("status"),
		StartDate:   r.FormValue("start_date"),
		TargetDate:  r.FormValue("target_completion_date"),
	}

	if data.Title == "" {
		data.Error = "Title is required."
		coursepages.Form(data).Render(r.Context(), w)
		return
	}

	status := models.CourseStatus(data.Status)
	if status == "" {
		status = models.CourseStatusActive
	}

	id, err := h.store.CreateCourse(store.CreateCourseParams{
		Title:                data.Title,
		Description:          data.Description,
		Status:               status,
		StartDate:            nilIfEmpty(data.StartDate),
		TargetCompletionDate: nilIfEmpty(data.TargetDate),
	})
	if err != nil {
		data.Error = "Could not create course. Please try again."
		coursepages.Form(data).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/courses/%d", id), http.StatusSeeOther)
}

func (h *Handler) GetCourse(w http.ResponseWriter, r *http.Request) {
	courseID, err := courseIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	course, err := h.store.GetCourse(courseID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	units, err := h.store.ListUnitsForCourse(courseID)
	if err != nil {
		http.Error(w, "could not load units", http.StatusInternalServerError)
		return
	}

	coursepages.Detail(course, units).Render(r.Context(), w)
}

func (h *Handler) EditCourseForm(w http.ResponseWriter, r *http.Request) {
	courseID, err := courseIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	course, err := h.store.GetCourse(courseID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	coursepages.Form(coursepages.CourseFormData{
		IsNew:       false,
		CourseID:    course.ID,
		Title:       course.Title,
		Description: course.Description,
		Status:      string(course.Status),
		StartDate:   derefStr(course.StartDate),
		TargetDate:  derefStr(course.TargetCompletionDate),
	}).Render(r.Context(), w)
}

func (h *Handler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	courseID, err := courseIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	data := coursepages.CourseFormData{
		IsNew:       false,
		CourseID:    courseID,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Status:      r.FormValue("status"),
		StartDate:   r.FormValue("start_date"),
		TargetDate:  r.FormValue("target_completion_date"),
	}

	if data.Title == "" {
		data.Error = "Title is required."
		coursepages.Form(data).Render(r.Context(), w)
		return
	}

	if err := h.store.UpdateCourse(store.UpdateCourseParams{
		ID:                   courseID,
		Title:                data.Title,
		Description:          data.Description,
		Status:               models.CourseStatus(data.Status),
		StartDate:            nilIfEmpty(data.StartDate),
		TargetCompletionDate: nilIfEmpty(data.TargetDate),
	}); err != nil {
		data.Error = "Could not update course. Please try again."
		coursepages.Form(data).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/courses/%d", courseID), http.StatusSeeOther)
}

func (h *Handler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	courseID, err := courseIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := h.store.DeleteCourse(courseID); err != nil {
		http.Error(w, "could not delete course", http.StatusInternalServerError)
		return
	}

	// HTMX delete from the list page: trigger a redirect to /courses.
	w.Header().Set("HX-Redirect", "/courses")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateCourseStatus(w http.ResponseWriter, r *http.Request) {
	courseID, err := courseIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	status := models.CourseStatus(r.FormValue("status"))
	if err := h.store.SetCourseStatus(courseID, status); err != nil {
		http.Error(w, "could not update status", http.StatusInternalServerError)
		return
	}

	// Reload the full detail page so all derived fields (completion date, etc.) refresh.
	http.Redirect(w, r, fmt.Sprintf("/courses/%d", courseID), http.StatusSeeOther)
}

// ── Unit handlers ─────────────────────────────────────────────────────────────

func (h *Handler) NewUnitForm(w http.ResponseWriter, r *http.Request) {
	courseID, err := courseIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	coursepages.UnitForm(coursepages.UnitFormData{
		IsNew:    true,
		CourseID: courseID,
	}).Render(r.Context(), w)
}

func (h *Handler) CreateUnit(w http.ResponseWriter, r *http.Request) {
	courseID, err := courseIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	data := coursepages.UnitFormData{
		IsNew:      true,
		CourseID:   courseID,
		Title:      r.FormValue("title"),
		Desc:       r.FormValue("description"),
		TargetDate: r.FormValue("target_completion_date"),
	}

	if data.Title == "" {
		data.Error = "Title is required."
		coursepages.UnitForm(data).Render(r.Context(), w)
		return
	}

	if _, err := h.store.CreateUnit(store.CreateUnitParams{
		CourseID:             courseID,
		Title:                data.Title,
		Description:          data.Desc,
		TargetCompletionDate: nilIfEmpty(data.TargetDate),
	}); err != nil {
		data.Error = "Could not create unit. Please try again."
		coursepages.UnitForm(data).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/courses/%d", courseID), http.StatusSeeOther)
}

func (h *Handler) EditUnitForm(w http.ResponseWriter, r *http.Request) {
	courseID, err := courseIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	unitID, err := unitIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	unit, err := h.store.GetUnit(unitID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	coursepages.UnitForm(coursepages.UnitFormData{
		IsNew:      false,
		CourseID:   courseID,
		UnitID:     unit.ID,
		Title:      unit.Title,
		Desc:       unit.Description,
		TargetDate: derefStr(unit.TargetCompletionDate),
	}).Render(r.Context(), w)
}

func (h *Handler) UpdateUnit(w http.ResponseWriter, r *http.Request) {
	courseID, err := courseIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	unitID, err := unitIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	data := coursepages.UnitFormData{
		IsNew:      false,
		CourseID:   courseID,
		UnitID:     unitID,
		Title:      r.FormValue("title"),
		Desc:       r.FormValue("description"),
		TargetDate: r.FormValue("target_completion_date"),
	}

	if data.Title == "" {
		data.Error = "Title is required."
		coursepages.UnitForm(data).Render(r.Context(), w)
		return
	}

	if err := h.store.UpdateUnit(store.UpdateUnitParams{
		ID:                   unitID,
		Title:                data.Title,
		Description:          data.Desc,
		TargetCompletionDate: nilIfEmpty(data.TargetDate),
	}); err != nil {
		data.Error = "Could not update unit. Please try again."
		coursepages.UnitForm(data).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/courses/%d", courseID), http.StatusSeeOther)
}

func (h *Handler) DeleteUnit(w http.ResponseWriter, r *http.Request) {
	courseID, err := courseIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	unitID, err := unitIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := h.store.DeleteUnit(unitID); err != nil {
		http.Error(w, "could not delete unit", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/courses/%d", courseID))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ReorderUnits(w http.ResponseWriter, r *http.Request) {
	courseID, err := courseIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	rawIDs := r.Form["ids"]
	orderedIDs := make([]int64, 0, len(rawIDs))
	for _, raw := range rawIDs {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		orderedIDs = append(orderedIDs, id)
	}

	if err := h.store.ReorderUnits(courseID, orderedIDs); err != nil {
		http.Error(w, "could not reorder units", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// nilIfEmpty wraps the store helper so handlers can call it directly.
func nilIfEmpty(s string) *string {
	return store.NilIfEmpty(s)
}
