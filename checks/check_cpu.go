package checks

import (
	"time"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/cpu"
)

func CheckCPU() *cmt.Check {
	check := &cmt.Check{}
	cpuPercent, err := cpu.Percent(2*time.Second, false)
	if err != nil {
		check.AddError(errors.Wrap(err, "cpu.Percent"))
		return nil
	}

	check.AddItem(&cmt.CheckItem{
		Name:        "CPU",
		Value:       cpuPercent,
		Description: "CPU Percentage",
		Unit:        "%",
	})

	return check
}
