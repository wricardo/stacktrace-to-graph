package subdir

import stacktracetograph "github.com/wricardo/stacktrace-to-graph"

func DoSomething() {
	stacktracetograph.ReportStacktrace()
}
