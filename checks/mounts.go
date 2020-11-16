package checks

import (
	"fmt"

	"github.com/math2001/gocmt/cmt"
	"github.com/shirou/gopsutil/disk"
)

func Mounts(c *cmt.Check, args map[string]interface{}) {
	mountpoint := args["mountpoint"].(string)
	c.AddItem(&cmt.CheckItem{
		Name:  "mount",
		Value: mountpoint,
	})

	partitions, err := disk.Partitions(false)
	if err != nil {
		panic(err)
	}
	ci := &cmt.CheckItem{
		Name:        "mount_status",
		Description: "ok/nok",
	}
	for _, partition := range partitions {
		if partition.Mountpoint == mountpoint {
			ci.Value = "ok"
			c.AddItem(ci)
			return
		}
	}

	ci.Value = "nok"
	ci.AlertMessage = fmt.Sprintf("check_mount - %s not found", mountpoint)
	ci.IsAlert = true

	c.AddItem(ci)
}
