package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var Template = template.Must(template.ParseGlob("*.html"))

func job_handler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/job/")
	job := get_job(id)
	if job == nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Job      *Job
		Dispatch *JobDispatch
		Results  *Results
	}{Job: job}

	if job.Status == "pending" || job.Status == "complete" || job.Status == "error" {
		dispatch := get_dispatch(job.Id)
		data.Dispatch = &dispatch
	}

	if job.Status == "complete" {
		results := get_full_results(job.Id)
		data.Results = &results
	}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	err := Template.ExecuteTemplate(w, "results.html", data)
	if err != nil {
		fmt.Println(err)
	}
}

func jobs_handler(w http.ResponseWriter, r *http.Request) {
	jobs := read_jobs()
	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	err := Template.ExecuteTemplate(w, "jobs.html", jobs)
	if err != nil {
		fmt.Println(err)
	}
}

func artifacts_handler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/artifacts/")
	job := get_job(id)
	if job == nil {
		http.NotFound(w, r)
		return
	}

	artifacts := get_artifacts(job.Id)

	data := struct {
		Job       *Job
		Artifacts []Artifact
	}{Job: job, Artifacts: artifacts}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	err := Template.ExecuteTemplate(w, "artifacts.html", data)
	if err != nil {
		fmt.Println(err)
	}
}
func main() {
	http.HandleFunc("/jobs/", jobs_handler)
	http.HandleFunc("/job/", job_handler)
	http.HandleFunc("/artifacts/", artifacts_handler)
	http.Handle("/", http.RedirectHandler("/jobs/", http.StatusFound))
	http.ListenAndServe(":8080", nil)
}
