package jobs

import (
	"log"
	"time"
)

func StartWorker(id int) {
	for job := range Queue {
		log.Println("Worker", id, "processing job", job.ID)

		job.Status = "running"

		// simula processamento pesado
		time.Sleep(2 * time.Second)

		job.Result = "processed: " + job.URL
		job.Status = "done"

		log.Println("Worker", id, "finished job", job.ID)
	}
}
