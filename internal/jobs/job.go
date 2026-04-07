package jobs

import "sync"

type Job struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	URL    string `json:"image_url"`
	Result string `json:"result"`
	Thumb  string `json:"thumbnail"`
	Error  string `json:"error"`
}

type Store struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

func NewStore() *Store {
	return &Store{
		jobs: make(map[string]*Job),
	}
}

func (s *Store) Save(job *Job) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs[job.ID] = cloneJob(job)
}

func (s *Store) Get(id string) (*Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, ok := s.jobs[id]
	if !ok {
		return nil, false
	}

	return cloneJob(job), true
}

func (s *Store) Update(id string, updateFn func(job *Job)) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		return nil, false
	}

	updateFn(job)
	return cloneJob(job), true
}

func cloneJob(job *Job) *Job {
	if job == nil {
		return nil
	}

	jobCopy := *job
	return &jobCopy
}

var JobStore = NewStore()
var Queue = make(chan string, 100)
