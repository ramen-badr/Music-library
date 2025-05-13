package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	_ "music-library/docs"
	"music-library/internal/client"
	"music-library/internal/config"
	"music-library/internal/handler"
	"music-library/internal/storage"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	if err = runMigrations(dsn); err != nil {
		log.Fatal("Migration failed:", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	apiClient := client.NewAPIClient(cfg.APIBaseURL)
	str := storage.NewPostgresStorage(db)
	songHandler := handler.NewSongHandler(str, apiClient)

	// Роутер
	router := mux.NewRouter()
	router.HandleFunc("/songs", songHandler.GetSongs).Methods("GET")
	router.HandleFunc("/songs/{id}/text", songHandler.GetSongText).Methods("GET")
	router.HandleFunc("/songs/{id}", songHandler.DeleteSong).Methods("DELETE")
	router.HandleFunc("/songs/{id}", songHandler.UpdateSong).Methods("PUT")
	router.HandleFunc("/songs", songHandler.CreateSong).Methods("POST")

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}

func runMigrations(dsn string) error {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
