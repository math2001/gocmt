// Manage the third componenent:
//   - sending updates to the server (one message per check)
//   - printing the CLI

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/math2001/gocmt/cmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const STDOUT_REPORT_WIDTH = 48

var httpclient = &http.Client{
	Timeout: 100 * time.Second, // default timeout is 0, meaning no timeout, which is bad
}

// consumes the reports from the check results channel (the check result are
// send on this channel as soon as they are finished, and the channel is closed
// once all the checks have finished)
func sendReport(fs *FrameworkSettings, checkResults <-chan *cmt.Check) {
	writeReportHeaderToStdout(fs)

	var g errgroup.Group
	for checkresult := range checkResults {
		checkresult := checkresult
		g.Go(func() error {
			return sendCheckResult(fs, checkresult)
		})

		writeReportToStdout(checkresult)
	}

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}

func sendCheckResult(fs *FrameworkSettings, c *cmt.Check) error {

	// don't send update if there are no check items
	if len(c.CheckItems()) == 0 {
		return nil
	}

	var g errgroup.Group
	for _, addr := range fs.GraylogHTTPGelfServers {
		addr := addr
		g.Go(func() error {
			return sendCheckResultHTTPGelf(c, addr, fs.CmtGroup, fs.CmtNode)
		})

	}

	for _, addr := range fs.GraylogUDPGelfServers {
		addr := addr
		g.Go(func() error {
			return sendCheckResultUDPGelf(c, addr, fs.CmtGroup, fs.CmtNode)
		})
	}

	for _, addr := range fs.TeamsChannel {
		addr := addr
		g.Go(func() error {
			return sendCheckResultTeamsChannel(c, addr, fs.CmtGroup, fs.CmtNode)
		})
	}

	return errors.Wrapf(g.Wait(), "reporting check result %q", c.Name())
}

func sendCheckResultHTTPGelf(c *cmt.Check, addr *HTTPGelfAddress, group string, node string) error {
	var buf bytes.Buffer
	payload := map[string]interface{}{
		"version":       "1.1",
		"host":          fmt.Sprintf("%s_%s", group, node),
		"short_message": fmt.Sprintf("cmt_check %s", c.Name()),
		"cmt_check":     c.Name(),
		"cmt_node":      node,
		"cmt_group":     group,
	}

	var hasAlert bool
	for _, ci := range c.CheckItems() {
		payload[fmt.Sprintf("cmt_%s", ci.Name)] = ci.Value
		if ci.IsAlert {
			hasAlert = true
		}
	}

	if hasAlert {
		payload["cmt_alert"] = "yes"
	} else {
		payload["cmt_alert"] = "no"
	}

	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return errors.Wrapf(err, "encoding payload in JSON")
	}

	req, err := http.NewRequest(http.MethodPost, addr.URL, &buf)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "text/plain")

	if err != nil {
		return errors.Wrapf(err, "preparing http request for %q", addr.Name)
	}

	res, err := httpclient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "performing http request for %q", addr.Name)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code 200 OK for %q, got %s", addr.Name, res.Status)
	}

	return nil
}

func sendCheckResultUDPGelf(c *cmt.Check, addr *UDPGelfAddress, group string, node string) error {
	return nil
}

func sendCheckResultTeamsChannel(c *cmt.Check, addr *TeamsAddress, group string, node string) error {
	return nil
}

func writeReportHeaderToStdout(fs *FrameworkSettings) {
	fmt.Println(strings.Repeat("=", STDOUT_REPORT_WIDTH))
	printCentered(fmt.Sprintf("%s:%s (ran %d checks)", fs.CmtGroup, fs.CmtNode, len(fs.Checks)), STDOUT_REPORT_WIDTH-1, ' ')
	fmt.Println(strings.Repeat("=", STDOUT_REPORT_WIDTH))
	fmt.Println()
}

func writeReportToStdout(checkresult *cmt.Check) {
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

	if msg, stack := checkresult.GetPanic(); msg != nil || stack != nil {
		fmt.Println("Panic:")
		fmt.Println(msg)
		fmt.Println(string(stack))
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
