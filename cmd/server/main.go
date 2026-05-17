package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"theology-tracker/internal/config"
	"theology-tracker/internal/database"
	"theology-tracker/internal/handlers"
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

	// Static assets
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Pages — routes will grow as each phase is implemented
	r.Get("/", h.Dashboard)

	log.Printf("server listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
