package jobs

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func StartWorker(id int) {
	for jobID := range Queue {
		processJob(id, jobID)
	}
}

func processJob(id int, jobID string) {
	job, exists := JobStore.Get(jobID)
	if !exists {
		log.Println("Worker", id, "job not found", jobID)
		return
	}

	log.Println("Worker", id, "processing job", job.ID)

	updateJob(job.ID, func(job *Job) {
		job.Status = "running"
		job.Error = ""
	})

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := downloadWithRetries(client, job, 3)
	if err != nil {
		failJob(job.ID, err.Error())
		log.Println("Worker", id, "failed job", job.ID, err)
		return
	}
	defer resp.Body.Close()

	savedFilePath := filepath.Join("output", job.ID+detectImageExtension(resp.Header.Get("Content-Type")))
	file, err := os.Create(savedFilePath)
	if err != nil {
		failJob(job.ID, err.Error())
		log.Println("Worker", id, "failed job", job.ID, err)
		return
	}

	_, err = io.Copy(file, resp.Body)
	closeErr := file.Close()
	if err != nil {
		failJob(job.ID, err.Error())
		log.Println("Worker", id, "failed job", job.ID, err)
		return
	}
	if closeErr != nil {
		failJob(job.ID, closeErr.Error())
		log.Println("Worker", id, "failed job", job.ID, closeErr)
		return
	}

	thumbPath := filepath.Join("output", "thumb_"+job.ID+".jpg")
	err = CreateThumbnail(savedFilePath, thumbPath)
	if err != nil {
		failJob(job.ID, err.Error())
		log.Println("Worker", id, "failed thumbnail job", job.ID, err)
		return
	}

	updateJob(job.ID, func(job *Job) {
		job.Result = savedFilePath
		job.Thumb = thumbPath
		job.Status = "done"
	})

	log.Println("Worker", id, "finished job", job.ID)
}

func downloadWithRetries(client *http.Client, job *Job, maxRetries int) (*http.Response, error) {
	var resp *http.Response
	var err error
	var lastStatus string

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = client.Get(job.URL)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		if resp != nil {
			lastStatus = resp.Status
			resp.Body.Close()
		}

		log.Println("retry", attempt, "for job", job.ID)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed after retries: %w", err)
	}

	return nil, fmt.Errorf("failed after retries: status %s", lastStatus)
}

func failJob(jobID string, message string) {
	updateJob(jobID, func(job *Job) {
		job.Status = "failed"
		job.Error = message
	})
}

func updateJob(jobID string, updateFn func(job *Job)) {
	_, _ = JobStore.Update(jobID, updateFn)
}

func detectImageExtension(contentType string) string {
	contentType = strings.TrimSpace(strings.Split(contentType, ";")[0])
	if contentType == "" {
		return ".jpg"
	}

	extensions, err := mime.ExtensionsByType(contentType)
	if err != nil || len(extensions) == 0 {
		return ".jpg"
	}

	switch extensions[0] {
	case ".jpe":
		return ".jpg"
	default:
		return extensions[0]
	}
}
