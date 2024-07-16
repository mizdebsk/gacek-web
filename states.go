package main

// tmt test results: pass info warn fail skip error
// Ref: tmt.git spec/plans/results.fmf
// Ref: https://tmt.readthedocs.io/en/stable/spec/plans.html#results-format
const (
	TMT_PASS TestStatus = iota
	TMT_INFO
	TMT_WARN
	TMT_FAIL
	TMT_SKIP
	TMT_ERROR
)

func (s TestStatus) String() string {
	switch s {
	case TMT_PASS:
		return "pass"
	case TMT_INFO:
		return "info"
	case TMT_WARN:
		return "warn"
	case TMT_FAIL:
		return "fail"
	case TMT_SKIP:
		return "skip"
	default:
		return "error"
	}
}

func (status TestStatus) IsBad() bool {
	switch status {
	case TMT_PASS, TMT_INFO, TMT_SKIP:
		return false
	default:
		return true
	}
}

func (status TestStatus) Emoji() string {
	switch status {
	case TMT_PASS:
		return "✅"
	case TMT_INFO:
		return "ℹ️"
	case TMT_WARN:
		return "❓"
	case TMT_FAIL:
		return "⛔"
	case TMT_SKIP:
		return "❎"
	default:
		return "❌"
	}
}

func (job Job) Emoji() string {
	if job.Results != nil {
		return job.Results.Overall.Emoji()
	}
	if job.Status == "new" || job.Status == "pending" {
		return "🔧"
	}
	return "❌"
}

// TF test results: passed info needs_inspection failed not_applicable error
// Ref: gluetool-modules.git gluetool_modules_framework/testing/test_schedule_tmt.py
// Ref: https://pagure.io/fedora-ci/messages/blob/master/f/schemas/test-complete.yaml
func parse_tf_result(tf_state string) TestStatus {
	switch tf_state {
	case "passed":
		return TMT_PASS
	case "info":
		return TMT_INFO
	case "needs_inspection":
		return TMT_WARN
	case "failed":
		return TMT_FAIL
	case "not_applicable":
		return TMT_SKIP
	default:
		return TMT_ERROR
	}
}
