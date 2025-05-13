package storage

import (
	"database/sql"
	"fmt"
	"log"
	"music-library/internal/model"
	"strings"
)

type Storage interface {
	CreateSong(song *model.Song) error
	UpdateSong(id int, song *model.Song) error
	DeleteSong(id int) error
	GetSong(id int) (*model.Song, error)
	GetSongs(filter map[string]string, page, limit int) ([]model.Song, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) CreateSong(song *model.Song) error {
	query := `
		INSERT INTO songs ("group", song, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	return s.db.QueryRow(
		query,
		song.Group,
		song.Song,
		song.ReleaseDate,
		song.Text,
		song.Link,
	).Scan(&song.ID)
}

func (s *PostgresStorage) UpdateSong(id int, song *model.Song) error {
	query := `
		UPDATE songs
		SET "group" = $1, song = $2
		WHERE id = $3`
	_, err := s.db.Exec(query, song.Group, song.Song, id)
	return err
}

func (s *PostgresStorage) DeleteSong(id int) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *PostgresStorage) GetSong(id int) (*model.Song, error) {
	query := `
		SELECT id, "group", song, release_date, text, link
		FROM songs WHERE id = $1`
	var song model.Song
	err := s.db.QueryRow(query, id).Scan(
		&song.ID,
		&song.Group,
		&song.Song,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
	)
	return &song, err
}

func (s *PostgresStorage) GetSongs(filter map[string]string, page, limit int) ([]model.Song, error) {
	query := `SELECT id, "group", song, release_date, text, link FROM songs`
	var args []interface{}
	var conditions []string
	counter := 1

	for key, val := range filter {
		if val != "" {
			conditions = append(conditions, fmt.Sprintf("%s = $%d", key, counter))
			args = append(args, val)
			counter++
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, (page-1)*limit)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []model.Song
	for rows.Next() {
		var song model.Song
		if err = rows.Scan(
			&song.ID,
			&song.Group,
			&song.Song,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
		); err != nil {
			log.Printf("Error scanning song: %v", err)
			continue
		}
		songs = append(songs, song)
	}

	return songs, nil
}
