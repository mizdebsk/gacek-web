package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var funcMap = template.FuncMap{
	"emoji":  emoji,
	"is_bad": is_bad,
}

var Template = template.Must(template.New("").Funcs(funcMap).ParseGlob("*.html"))

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

	if job.Status == "pending" || job.Status == "complete" {
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

func main() {
	http.HandleFunc("/jobs/", jobs_handler)
	http.HandleFunc("/job/", job_handler)
	http.Handle("/", http.RedirectHandler("/jobs/", http.StatusFound))
	http.ListenAndServe(":8080", nil)
}
