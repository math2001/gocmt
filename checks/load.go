package checks

import (
	"github.com/math2001/gocmt/cmt"
	"github.com/shirou/gopsutil/load"
)

func Load(
	cr *cmt.CheckResult,
	globals map[string]interface{},
	settings map[string]interface{},
) {
	loadavg, err := load.Avg()
	if err != nil {
		cr.AddError(err)
		return
	}

	cr.AddItem(&cmt.CheckItem{
		Name:        "load1",
		Value:       loadavg.Load1,
		Description: "CPU Load Average, one minute",
	})

	cr.AddItem(&cmt.CheckItem{
		Name:        "load5",
		Value:       loadavg.Load5,
		Description: "CPU Load Average, five minutes",
	})

	cr.AddItem(&cmt.CheckItem{
		Name:        "load15",
		Value:       loadavg.Load15,
		Description: "CPU Load Average, fifteen minutes",
	})

}
