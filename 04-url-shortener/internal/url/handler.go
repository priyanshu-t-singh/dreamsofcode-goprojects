package url

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/04-url-shortener/internal/core/models"
)

type Handler struct {
	Repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{Repository: repository}
}

type CreateShortURLSQL struct {
	ShortCode   string
	OriginalURL string
}

func sendErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
}

func (h *Handler) SaveShortURL(w http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("url")

	if originalURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		sendErrorResponse(w, fmt.Errorf("URL is required"))
		return
	}

	shortCode, err := generateCode(6)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendErrorResponse(w, fmt.Errorf("Internal Server Error"))
		return
	}

	err2 := h.Repository.SaveShortURL(&CreateShortURLSQL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
	})

	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendErrorResponse(w, err2)
		return
	}

	http.Redirect(w, r, "/r/"+shortCode, http.StatusFound)
}

func (h *Handler) RedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	pageId := r.PathValue("pageId")
	originalURL, err := h.Repository.GetOriginalURL(pageId)

	if err != nil {
		if originalURL == "-1" {
			w.WriteHeader(http.StatusNotFound)
			sendErrorResponse(w, err)
			return
		}
		sendErrorResponse(w, err)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
