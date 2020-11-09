// Manage the third componenent:
//   - sending updates to the server (one message per check)
//   - printing the CLI

package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/math2001/gocmt/cmt"
)

func sendReports(fs *cmt.FrameworkSettings, checkResults <-chan *cmt.CheckResult) {
	var stdoutlock sync.Mutex

	writeReportHeaderToStdout(fs)

	var wg sync.WaitGroup
	for checkResult := range checkResults {
		wg.Add(1)
		go func(checkresult *cmt.CheckResult) {
			defer wg.Done()
			sendReport(fs, checkresult)

			// make sure that two tests don't write to stdout at the same time
			// (the output would be very messy)
			stdoutlock.Lock()
			writeReportToStdout(checkresult)
			stdoutlock.Unlock()
		}(checkResult)
	}
	wg.Wait()
}

func sendReport(fs *cmt.FrameworkSettings, checkresult *cmt.CheckResult) {
	log.Printf("TODO: send report")
}

func writeReportHeaderToStdout(fs *cmt.FrameworkSettings) {
	fmt.Println("======================================")
	fmt.Printf("%s:%s (ran %d check(s))\n", fs.CmtGroup, fs.CmtNode, len(fs.Checks))
	fmt.Println("======================================")
	fmt.Println()
}

func writeReportToStdout(checkresult *cmt.CheckResult) {
	// log.Printf("TODO: pretty print\n%#v\n", checkresult)

	fmt.Printf("::: %s\n", checkresult.Name())
	for _, checkitem := range checkresult.CheckItems() {
		fmt.Printf("%-10s %v %s -- %s\n", checkitem.Name, checkitem.Value, checkitem.Unit, checkitem.Description)
	}
	fmt.Println()
}
