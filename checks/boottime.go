package checks

import (
	"time"

	"github.com/math2001/gocmt/cmt"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/host"
)

func Boottime(
	c *cmt.Check,
	args map[string]interface{},
) {
	boottime, err := host.BootTime()
	if err != nil {
		c.AddError(errors.Wrapf(err, "getting boottime"))
	}

	boottime_seconds := uint64(time.Now().Unix()) - boottime

	c.AddItem(&cmt.CheckItem{
		Name:        "cmt_boottime_seconds",
		Value:       boottime_seconds,
		Unit:        "seconds",
		Description: "Seconds since last reboot",
	})
	c.AddItem(&cmt.CheckItem{
		Name:        "cmt_boottime_days",
		Value:       boottime_seconds / (24 * 60 * 60),
		Unit:        "days",
		Description: "Days since last reboot",
	})
}
