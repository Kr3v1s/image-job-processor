package jobs

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func StartWorker(id int) {
	for job := range Queue {
		log.Println("Worker", id, "processing job", job.ID)

		job.Status = "running"

		resp, err := http.Get(job.URL)
		if err != nil {
			job.Status = "failed"
			job.Error = err.Error()
			log.Println("Worker", id, "failed job", job.ID, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			job.Status = "failed"
			job.Error = "failed to download image: " + resp.Status
			resp.Body.Close()
			log.Println("Worker", id, "failed job", job.ID, job.Error)
			continue
		}

		filename := filepath.Join("output", job.ID+".img")
		file, err := os.Create(filename)
		if err != nil {
			job.Status = "failed"
			job.Error = err.Error()
			resp.Body.Close()
			log.Println("Worker", id, "failed job", job.ID, err)
			continue
		}

		_, err = io.Copy(file, resp.Body)
		resp.Body.Close()
		file.Close()
		if err != nil {
			job.Status = "failed"
			job.Error = err.Error()
			log.Println("Worker", id, "failed job", job.ID, err)
			continue
		}

		job.Result = filename
		job.Status = "done"

		log.Println("Worker", id, "finished job", job.ID)
	}
}
