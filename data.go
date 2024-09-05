package main

import (
	"encoding/xml"
	"fmt"
	//"io/ioutil"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"slices"
	"strings"
)

var gacek_home = os.Getenv("GACEK_HOME")
var jobs_dir = gacek_home + "/jobs"
var queues_dir = gacek_home + "/queues"

func read_jobs() []Job {
	queues := []string{"new", "pending", "complete", "error"}
	jobs := []Job{}
	for _, queue := range queues {
		dir, err := os.Open(queues_dir + "/" + queue)
		if err != nil {
			log.Fatal(err)
		}
		defer dir.Close()
		entries, err := dir.Readdirnames(0)
		if err != nil {
			log.Fatal(err)
		}
		for _, entry := range entries {
			job := Job{Id: entry, Status: queue}
			if job.Status == "complete" {
				results := get_results(job.Id)
				job.Results = &results
			}
			jobs = append(jobs, job)
		}
	}

	slices.SortStableFunc(jobs, func(ja, jb Job) int {
		aa := strings.Split(ja.Id, "-")
		bb := strings.Split(jb.Id, "-")
		at := strings.Join(aa[1:], "-")
		bt := strings.Join(bb[1:], "-")
		r := -strings.Compare(at, bt)
		if r != 0 {
			return r
		}
		return strings.Compare(aa[0], bb[0])
	})
	return jobs
}

func get_job(id string) *Job {
	for _, job := range read_jobs() {
		if job.Id == id {
			return &job
		}
	}
	return nil
}

func parse_test_name(test *Test) {
	chunks := strings.Split(strings.TrimPrefix(test.Name, "/"), "/")
	if len(chunks) < 2 {
		test.Component = ""
		test.Path = strings.Join(chunks, "/")
		test.Link = ""
		log.Printf("Unable to parse test name: %s\n", test.Name)
		return
	}
	if chunks[0] == "mbici" {
		test.Component = chunks[1]
		test.Path = strings.Join(chunks[2:], "/")
		test.Link = ""
	} else if chunks[0] == "tests" {
		test.Component = chunks[1]
		test.Path = strings.Join(chunks[2:], "/")
		test.Link = fmt.Sprintf("https://src.fedoraproject.org/tests/%s/blob/main/f/%s",
			test.Component, test.Path)
	} else if chunks[1] == "tests" {
		test.Component = chunks[0]
		test.Path = strings.Join(chunks[2:], "/")
		test.Link = fmt.Sprintf("https://src.fedoraproject.org/rpms/%s/blob/main/f/tests/%s",
			test.Component, test.Path)
	} else if chunks[0] == "javapackages-validator" {
		test.Component = chunks[0]
		test.Path = strings.Join(chunks[1:], "/")
		if chunks[1] == "jp" {
			test.Link = "https://src.fedoraproject.org/tests/javapackages/tree/main"
		} else {
			test.Link = "https://github.com/fedora-java/javapackages-validator"
		}
	} else if chunks[0] == "runit" {
		test.Component = chunks[0]
		test.Path = strings.Join(chunks[1:], "/")
		test.Link = "https://src.fedoraproject.org/tests/javapackages/blob/main/f/runit"
	} else {
		test.Component = ""
		test.Path = strings.Join(chunks, "/")
		test.Link = ""
		log.Printf("Unable to parse test name: %s\n", test.Name)
	}
}

func get_results(job_id string) Results {
	bytes, err := os.ReadFile(jobs_dir + "/" + job_id + "/results.xml")
	if err != nil {
		log.Fatal(err)
	}
	var results Results
	if err := xml.Unmarshal(bytes, &results); err != nil {
		log.Fatal(err)
	}

	results.Overall = parse_tf_result(results.OverallStr)
	iid := 0
	for pi, tfPlan := range results.Plans {
		results.Plans[pi].Result = parse_tf_result(tfPlan.ResultStr)
		for ti, tfTest := range tfPlan.Tests {
			results.Plans[pi].Tests[ti].Result = parse_tf_result(tfTest.ResultStr)
			results.Plans[pi].Tests[ti].IntId = iid
			iid++
			for li, log := range tfTest.Logs {
				if log.Name == "journal.txt" {
					results.Plans[pi].Tests[ti].Journal = &results.Plans[pi].Tests[ti].Logs[li]
				}
				if log.Name == "testout.log" {
					results.Plans[pi].Tests[ti].Testout = &results.Plans[pi].Tests[ti].Logs[li]
				}
			}

			parse_test_name(&results.Plans[pi].Tests[ti])
		}
		slices.SortStableFunc(results.Plans[pi].Tests, func(ta, tb Test) int {
			if ta.Component != tb.Component {
				return strings.Compare(ta.Component, tb.Component)
			}
			return strings.Compare(ta.Path, tb.Path)
		})
	}

	slices.SortStableFunc(results.Plans, func(pa, pb Plan) int {
		return strings.Compare(pa.Name, pb.Name)
	})
	return results
}

func get_full_results(job_id string) Results {
	results := get_results(job_id)

	for pi, tfPlan := range results.Plans {
		tmtTests := get_tests(job_id, strings.TrimPrefix(tfPlan.Name, "/"))
		for ti, tfTest := range tfPlan.Tests {
			for tti, tmtTest := range tmtTests {
				if tmtTest.Name == tfTest.Name {
					results.Plans[pi].Tests[ti].Info = &tmtTests[tti]
				}
			}
			if results.Plans[pi].Tests[ti].Info == nil {
				for tti, tmtTest := range tmtTests {
					if strings.HasPrefix(tfTest.Name+"/", tmtTest.Name) {
						results.Plans[pi].Tests[ti].Info = &tmtTests[tti]
					}
				}
			}
		}
	}

	return results
}

func get_tests(job_id string, plan_name string) []TmtTest {
	bytes, err := os.ReadFile(jobs_dir + "/" + job_id + "/" + plan_name + "-tests.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var tests []TmtTest
	if err := yaml.Unmarshal(bytes, &tests); err != nil {
		log.Fatal(err)
	}

	return tests
}

func get_dispatch(job_id string) JobDispatch {
	bytes, err := os.ReadFile(jobs_dir + "/" + job_id + "/" + "tf-dispatch.xml")
	if err != nil {
		log.Fatal(err)
	}
	var dispatch JobDispatch
	if err := xml.Unmarshal(bytes, &dispatch); err != nil {
		log.Fatal(err)
	}

	return dispatch
}

func get_artifacts(job_id string) []Artifact {
	bytes, err := os.ReadFile(jobs_dir + "/" + job_id + "/" + "artifacts.xml")
	if err != nil {
		log.Fatal(err)
	}
	artifacts := struct {
		Artifacts []Artifact `xml:"artifact"`
	}{}
	if err := xml.Unmarshal(bytes, &artifacts); err != nil {
		log.Fatal(err)
	}

	return artifacts.Artifacts
}
