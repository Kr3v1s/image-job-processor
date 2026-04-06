package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"image-job-processor/internal/jobs"
)

type CreateJobRequest struct {
	ImageURL string `json:"image_url"`
}

func CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateJobRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ImageURL == "" {
		http.Error(w, "image_url is required", http.StatusBadRequest)
		return
	}

	jobID := strconv.FormatInt(time.Now().UnixNano(), 10)

	job := &jobs.Job{
		ID:     jobID,
		Status: "pending",
		URL:    req.ImageURL,
	}

	jobs.Store[jobID] = job
	jobs.Queue <- job
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

func GetJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	job, exists := jobs.Store[id]
	if !exists {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}
