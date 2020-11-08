package checks

import (
	"time"

	"github.com/math2001/gocmt/cmt"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/cpu"
)

func CheckCPU(globals map[string]interface{}, settings map[string]interface{}) (check *cmt.CheckResult) {
	check = &cmt.CheckResult{}
	cpuPercent, err := cpu.Percent(2*time.Second, false)
	if err != nil {
		check.AddError(errors.Wrap(err, "cpu.Percent"))
		return
	}

	check.AddItem(&cmt.CheckItem{
		Name:        "CPU",
		Value:       cpuPercent,
		Description: "CPU Percentage",
		Unit:        "%",
	})

	return
}
