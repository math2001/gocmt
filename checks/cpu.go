package checks

import (
	"fmt"
	"time"

	"github.com/math2001/gocmt/cmt"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/cpu"
)

func CPU(
	c *cmt.Check,
	args map[string]interface{},
) {

	cpuPercent, err := cpu.Percent(2*time.Second, false)
	if err != nil {
		c.AddError(errors.Wrap(err, "cpu.Percent"))
		return
	}

	fmt.Fprintf(c.DebugBuffer(), "some debugging thing")

	c.AddItem(&cmt.CheckItem{
		Name:        "cmt_cpu",
		Value:       cpuPercent[0],
		Description: "CPU Percentage",
		Unit:        "%",
	})
}
