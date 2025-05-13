package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"music-library/internal/client"
	"music-library/internal/model"
	"music-library/internal/storage"

	"github.com/gorilla/mux"
)

type SongHandler struct {
	storage   storage.Storage
	apiClient *client.APIClient
}

func NewSongHandler(storage storage.Storage, apiClient *client.APIClient) *SongHandler {
	return &SongHandler{
		storage:   storage,
		apiClient: apiClient,
	}
}

// GetSongs возвращает список песен с фильтрацией и пагинацией
// @Summary Get songs
// @Description Get songs with filtering and pagination
// @Tags songs
// @Param group query string false "Filter by group"
// @Param song query string false "Filter by song name"
// @Param releaseDate query string false "Filter by release date"
// @Param link query string false "Filter by link"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10)"
// @Success 200 {array} model.Song
// @Router /songs [get]
func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	filter := map[string]string{
		"group":        r.URL.Query().Get("group"),
		"song":         r.URL.Query().Get("song"),
		"release_date": r.URL.Query().Get("releaseDate"),
		"link":         r.URL.Query().Get("link"),
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}

	songs, err := h.storage.GetSongs(filter, page, limit)
	if err != nil {
		log.Printf("Error fetching songs: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// GetSongText возвращает текст песни с пагинацией по куплетам
// @Summary Get song text
// @Description Get song text paginated by verses
// @Tags songs
// @Param id path int true "Song ID"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Verses per page (default 5)"
// @Success 200 {array} string
// @Router /songs/{id}/text [get]
func (h *SongHandler) GetSongText(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	song, err := h.storage.GetSong(id)
	if err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 5
	}

	verses := strings.Split(song.Text, "\n\n")
	start := (page - 1) * limit
	end := start + limit

	if start > len(verses) {
		start = len(verses)
	}

	if end > len(verses) {
		end = len(verses)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(verses[start:end])
}

// DeleteSong удаляет песню по ID
// @Summary Delete a song
// @Description Delete a song by ID
// @Tags songs
// @Param id path int true "Song ID"
// @Success 204
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := h.storage.DeleteSong(id); err != nil {
		log.Printf("Error deleting song: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateSong обновляет данные песни
// @Summary Update a song
// @Description Update song details
// @Tags songs
// @Param id path int true "Song ID"
// @Param song body model.Song true "Updated song data"
// @Success 200 {object} model.Song
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var song model.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.storage.UpdateSong(id, &song); err != nil {
		log.Printf("Error updating song: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// CreateSong добавляет новую песню
// @Summary Create a song
// @Description Add a new song with details from external API
// @Tags songs
// @Param song body model.Song true "Song data"
// @Success 201 {object} model.Song
// @Router /songs [post]
func (h *SongHandler) CreateSong(w http.ResponseWriter, r *http.Request) {
	var song model.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	detail, err := h.apiClient.FetchSongDetails(song.Group, song.Song)
	if err != nil {
		log.Printf("Error fetching song details: %v", err)
		http.Error(w, "Failed to fetch song details", http.StatusInternalServerError)
		return
	}

	song.ReleaseDate = detail.ReleaseDate
	song.Text = detail.Text
	song.Link = detail.Link

	if err := h.storage.CreateSong(&song); err != nil {
		log.Printf("Error creating song: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}
