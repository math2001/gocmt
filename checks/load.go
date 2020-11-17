package checks

import (
	"github.com/math2001/gocmt/cmt"
	"github.com/shirou/gopsutil/load"
)

func Load(
	c *cmt.CheckResult,
	args map[string]interface{},
) {
	loadavg, err := load.Avg()
	if err != nil {
		c.AddError(err)
		return
	}

	c.AddItem(&cmt.CheckItem{
		Name:        "load1",
		Value:       loadavg.Load1,
		Description: "CPU Load Average, one minute",
	})

	c.AddItem(&cmt.CheckItem{
		Name:        "load5",
		Value:       loadavg.Load5,
		Description: "CPU Load Average, five minutes",
	})

	c.AddItem(&cmt.CheckItem{
		Name:        "load15",
		Value:       loadavg.Load15,
		Description: "CPU Load Average, fifteen minutes",
	})

}
