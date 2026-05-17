package handlers

import (
	"net/http"

	"theology-tracker/internal/store"
	"theology-tracker/internal/templates/layout"
)

// Handler holds shared dependencies for all HTTP handlers.
// Each domain gets its own *_handler.go file with methods on this struct.
type Handler struct {
	store *store.Store
}

func New(s *store.Store) *Handler {
	return &Handler{store: s}
}

// isHTMX returns true when the request was initiated by HTMX.
// Handlers use this to decide whether to render a full page or a fragment.
func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	component := layout.Base("Dashboard")
	component.Render(r.Context(), w)
}
