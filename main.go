package main

import (
	"log"
	"net/http"

	handler "image-job-processor/internal/http"
	jobs "image-job-processor/internal/jobs"
)

func main() {
	for i := 1; i <= 3; i++ {
		go jobs.StartWorker(i)
	}

	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("./internal/http/static"))))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateJobHandler(w, r)
			return
		}

		if r.Method == http.MethodGet {
			handler.GetJobHandler(w, r)
			return
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
