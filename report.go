// Manage the third componenent:
//   - sending updates to the server (one message per check)
//   - printing the CLI

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/math2001/gocmt/cmt"
)

const STDOUT_REPORT_WIDTH = 48

func sendReports(fs *FrameworkSettings, checkResults <-chan *cmt.CheckResult) {
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

func sendReport(fs *FrameworkSettings, checkresult *cmt.CheckResult) {
	// log.Printf("TODO: send report")
}

func writeReportHeaderToStdout(fs *FrameworkSettings) {
	fmt.Println(strings.Repeat("=", STDOUT_REPORT_WIDTH))
	printCentered(fmt.Sprintf("%s:%s (ran %d check(s))", fs.CmtGroup, fs.CmtNode, len(fs.Checks)), STDOUT_REPORT_WIDTH-1, ' ')
	fmt.Println(strings.Repeat("=", STDOUT_REPORT_WIDTH))
	fmt.Println()
}

func writeReportToStdout(checkresult *cmt.CheckResult) {
	printCentered(checkresult.Name(), STDOUT_REPORT_WIDTH, '-')
	if checkresult.ArgumentSet() != nil {
		fmt.Printf("argument set: %v\n", checkresult.ArgumentSet())
	}

	for _, checkitem := range checkresult.CheckItems() {
		fmt.Printf("%-20s %v %s -> %s\n", checkitem.Name, checkitem.Value, checkitem.Unit, checkitem.Description)
	}
	fmt.Println()

	if len(checkresult.Errors()) > 0 {
		fmt.Println("Errors:")
		for _, err := range checkresult.Errors() {
			fmt.Println(err)
		}
		fmt.Println()
	}

	if checkresult.DebugBuffer().Len() > 0 {
		fmt.Println("Debug output:")

		// write all the characters, except the last one, to check if it is a newline.
		// if it isn't, then one is added automatically
		io.CopyN(os.Stdout, checkresult.DebugBuffer(), int64(checkresult.DebugBuffer().Len()-1))

		lastchar, err := checkresult.DebugBuffer().ReadByte()
		if err == nil {
			fmt.Printf("%c", lastchar)
		} else if err != io.EOF {
			fmt.Println()
			fmt.Println("Internal error when reading debug buffer: ", err)
			fmt.Println()
		}

		if lastchar != '\n' {
			fmt.Println()
		}

		fmt.Println()
	}

}

func printCentered(text string, width int, paddingChar rune) {
	// -2 for the spaces
	for i := 0; i < (width-len(text)-2)/2; i++ {
		fmt.Printf("%c", paddingChar)
	}
	fmt.Printf(" %s ", text)
	for i := 0; i < (width-len(text)-2)/2; i++ {
		fmt.Printf("%c", paddingChar)
	}
	fmt.Println()
}
