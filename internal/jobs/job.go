package jobs

type Job struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	URL    string `json:"image_url"`
	Result string `json:"result"`
	Error  string `json:"error"`
}

var Store = make(map[string]*Job)
var Queue = make(chan *Job, 100)
