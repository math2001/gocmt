// Manage the third componenent:
//   - sending updates to the server (one message per check)
//   - printing the CLI

package main

import (
	"log"
	"sync"

	"github.com/math2001/gocmt/cmt"
)

func sendReports(framework_conf map[string]interface{}, checkResults <-chan *cmt.CheckResult) {
	var stdoutlock sync.Mutex

	var wg sync.WaitGroup
	for checkResult := range checkResults {
		wg.Add(1)
		go func(checkresult *cmt.CheckResult) {
			defer wg.Done()
			sendReport(framework_conf, checkresult)

			// make sure that two tests don't write to stdout at the same time
			// (the output would be very messy)
			stdoutlock.Lock()
			writeReportToStdout(checkresult)
			stdoutlock.Unlock()
		}(checkResult)
	}
	wg.Wait()
}

func sendReport(framework_conf map[string]interface{}, checkresult *cmt.CheckResult) {
	log.Printf("TODO: send report")
}

func writeReportToStdout(checkresult *cmt.CheckResult) {
	log.Printf("TODO: pretty print\n%#v\n", checkresult)
}
