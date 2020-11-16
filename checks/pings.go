package checks

import (
	"fmt"
	"os/exec"

	"github.com/math2001/gocmt/cmt"
)

func Pings(c *cmt.Check, args map[string]interface{}) {
	hostname := args["hostname"].(string)
	c.AddItem(&cmt.CheckItem{
		Name:  "ping",
		Value: hostname,
	})
	ci := &cmt.CheckItem{
		Name:        "ping_status",
		Description: "ok/nok",
	}

	cmd := exec.Command("ping", "-c", "1", "-W", "2", hostname)
	err := cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		// exited with non-zero code
		ci.Value = "nok"
		ci.IsAlert = true
		ci.AlertMessage = fmt.Sprintf("check_ping - %s not responding", hostname)
	} else if err != nil {
		panic(err) // error like failed to write to stdout
	} else {
		ci.Value = "ok"
	}
	c.AddItem(ci)

}
