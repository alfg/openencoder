package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alfg/enc/types"
	"github.com/alfg/enc/worker"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

type request struct {
	Task        string `json:"task"`
	Source      string `json:"source"`
	Destination string `json:"dest"`
	Delay       string `json:"delay"`
}

type response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// NewServer creates a new server
func NewServer(port string) {
	http.HandleFunc("/api/workflow", workflowHandler)
	http.HandleFunc("/api/encode", encodeHandler)

	log.Info("started server on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func workflowHandler(w http.ResponseWriter, r *http.Request) {

	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Decode json.
	decoder := json.NewDecoder(r.Body)
	var t request
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	// Parse the durations.
	delay, err := time.ParseDuration(t.Delay)
	if err != nil {
		http.Error(w, "Bad delay value: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate delay is in range 1 to 10 seconds.
	if delay.Seconds() < 1 || delay.Seconds() > 10 {
		http.Error(w, "The delay must be between 1 and 10 seconds, inclusively.", http.StatusBadRequest)
		return
	}

	// Set name and validate value.
	task := t.Task
	if task == "" {
		http.Error(w, "You must specify a task.", http.StatusBadRequest)
		return
	}

	// Create Job and push the work onto the jobCh.
	// job := Job{name, duration, t.Job.Task, t.Job.Source, t.Job.Dest}
	job := types.Job{
		ID:          xid.New().String(),
		Task:        t.Task,
		Source:      t.Source,
		Destination: t.Destination,
		Delay:       delay,
	}
	go func() {
		log.Info("added: %s %s\n", job.Task, job.Delay)
		worker.Jobs <- job
	}()

	// Create response.
	resp := response{
		Message: "Job created",
		Status:  200,
	}
	j, _ := json.Marshal(resp)

	// Render success.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func encodeHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Decode json.
	decoder := json.NewDecoder(r.Body)
	var t request
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	// Parse the durations.
	delay, err := time.ParseDuration(t.Delay)
	if err != nil {
		http.Error(w, "Bad delay value: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate delay is in range 1 to 10 seconds.
	if delay.Seconds() < 1 || delay.Seconds() > 10 {
		http.Error(w, "The delay must be between 1 and 10 seconds, inclusively.", http.StatusBadRequest)
		return
	}

	// Set name and validate value.
	task := t.Task
	if task == "" {
		http.Error(w, "You must specify a task.", http.StatusBadRequest)
		return
	}

	// Create Job and push the work onto the jobs channel.
	job := types.Job{
		ID:          xid.New().String(),
		Task:        t.Task,
		Source:      t.Source,
		Destination: t.Destination,
		Delay:       delay,
	}
	go func() {
		log.Infof("added: %s %s\n", job.Task, job.Delay)
		worker.Jobs <- job
	}()

	// Create response.
	resp := response{
		Message: "Job created",
		Status:  200,
	}
	j, _ := json.Marshal(resp)

	// Render success.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}
