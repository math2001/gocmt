package checks

import (
	"fmt"

	"github.com/math2001/gocmt/cmt"
	"github.com/shirou/gopsutil/disk"
)

func Disks(
	c *cmt.Check,
	args map[string]interface{},
) {
	path := args["path"].(string)
	alert_threshold := args["alert"].(int)

	disk, err := disk.Usage(path)
	if err != nil {
		c.AddError(err)
	}

	c.AddItem(&cmt.CheckItem{
		Name:        "disk",
		Value:       path,
		Description: "Path",
	})

	c.AddItem(&cmt.CheckItem{
		Name:        "disk_total",
		Value:       disk.Total,
		Unit:        "bytes",
		Description: "Total (bytes)",
	})

	c.AddItem(&cmt.CheckItem{
		Name:        "disk_free",
		Value:       disk.Free,
		Unit:        "bytes",
		Description: "Free (bytes)",
	})

	c.AddItem(&cmt.CheckItem{
		Name:        "disk_percent",
		Value:       disk.UsedPercent,
		Unit:        "%",
		Description: "Used (percent)",
	})

	if disk.UsedPercent > float64(alert_threshold) {
		c.SetAlert(fmt.Sprintf("check disk for %s - critical capacity alert (%.2f%%)", path, disk.UsedPercent))
	}
}
