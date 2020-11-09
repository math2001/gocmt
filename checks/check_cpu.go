package checks

import (
	"time"

	"github.com/math2001/gocmt/cmt"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/cpu"
)

func CheckCPU(cr *cmt.CheckResult, globals map[string]interface{}, settings map[string]interface{}) {
	cpuPercent, err := cpu.Percent(2*time.Second, false)
	if err != nil {
		cr.AddError(errors.Wrap(err, "cpu.Percent"))
		return
	}

	cr.AddItem(&cmt.CheckItem{
		Name:        "CPU",
		Value:       cpuPercent[0],
		Description: "CPU Percentage",
		Unit:        "%",
	})
}
