package main

import (
	"encoding/xml"
)

type Job struct {
	Id      string
	Status  string
	Results *Results
}

type TmtTest struct {
	Name        string   `yaml:"name"`
	Summary     string   `yaml:"summary"`
	Description string   `yaml:"description"`
	Tier        string   `yaml:"tier"`
	Duration    string   `yaml:"duration"`
	Contacts    []string `yaml:"contact"`
	Tags        []string `yaml:"tag"`
}

type Results struct {
	XMLName    xml.Name   `xml:"testsuites"`
	OverallStr string     `xml:"overall-result,attr"`
	Overall    TestStatus `xml:"-"`
	Plans      []Plan     `xml:"testsuite"`
}

type Plan struct {
	XMLName   xml.Name   `xml:"testsuite"`
	Name      string     `xml:"name,attr"`
	ResultStr string     `xml:"result,attr"`
	Result    TestStatus `xml:"-"`
	Tests     []Test     `xml:"testcase"`
	Logs      []Log      `xml:"logs>log"`
}

type Test struct {
	XMLName   xml.Name   `xml:"testcase"`
	IntId     int        `xml:"-"`
	Name      string     `xml:"name,attr"`
	ResultStr string     `xml:"result,attr"`
	Result    TestStatus `xml:"-"`
	Time      string     `xml:"time,attr"`
	Logs      []Log      `xml:"logs>log"`
	Journal   *Log       `xml:"-"`
	Testout   *Log       `xml:"-"`
	Info      *TmtTest   `xml:"-"`
	Link      string     `xml:"-"`
}

type Log struct {
	XMLName xml.Name `xml:"log"`
	Name    string   `xml:"name,attr"`
	Url     string   `xml:"href,attr"`
}

type JobDispatch struct {
	XMLName xml.Name `xml:"dispatch"`
	TfId    string   `xml:"tfId"`
}

type TestStatus int

type Artifact struct {
	Name    string `xml:"name"`
	NVR     string `xml:"nvr"`
	BuildId string `xml:"buildId"`
	TaskId  string `xml:"taskId"`
}
