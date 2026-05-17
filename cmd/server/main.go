package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"theology-tracker/internal/config"
	"theology-tracker/internal/database"
	"theology-tracker/internal/handlers"
	appmw "theology-tracker/internal/middleware"
	"theology-tracker/internal/store"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	s := store.New(db)
	h := handlers.New(s)

	r := chi.NewRouter()
	r.Use(chimw.Recoverer)
	r.Use(chimw.Logger)
	r.Use(appmw.MethodOverride) // must come before routing so _method is applied first

	// Static assets
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Dashboard
	r.Get("/", h.Dashboard)

	// Courses
	r.Get("/courses", h.ListCourses)
	r.Get("/courses/new", h.NewCourseForm)
	r.Post("/courses", h.CreateCourse)
	r.Get("/courses/{courseID}", h.GetCourse)
	r.Get("/courses/{courseID}/edit", h.EditCourseForm)
	r.Put("/courses/{courseID}", h.UpdateCourse)
	r.Delete("/courses/{courseID}", h.DeleteCourse)
	r.Post("/courses/{courseID}/status", h.UpdateCourseStatus)

	// Units (nested under courses)
	r.Get("/courses/{courseID}/units/new", h.NewUnitForm)
	r.Post("/courses/{courseID}/units", h.CreateUnit)
	r.Get("/courses/{courseID}/units/{unitID}/edit", h.EditUnitForm)
	r.Put("/courses/{courseID}/units/{unitID}", h.UpdateUnit)
	r.Delete("/courses/{courseID}/units/{unitID}", h.DeleteUnit)
	r.Post("/courses/{courseID}/units/reorder", h.ReorderUnits)

	log.Printf("server listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
