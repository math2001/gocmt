package checks

import (
	"fmt"

	"github.com/math2001/gocmt/cmt"
	"github.com/math2001/gocmt/columnprint"
	"github.com/shirou/gopsutil/disk"
)

func Mounts(c *cmt.CheckResult, args map[string]interface{}) {
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

func AvailMounts() {
	partitions, err := disk.Partitions(false)
	if err != nil {
		panic(err)
	}

	var u columnprint.U
	u.SetColumns("%s", "%s", "%s", "%s")
	u.WouldPrintLiteral("process", "device", "fstype", "opts")
	for _, p := range partitions {
		u.WouldPrint(p.Mountpoint, p.Device, p.Fstype, p.Opts)
	}
	u.PrintLiteral("process", "device", "fstype", "opts")
	for _, p := range partitions {
		u.Print(p.Mountpoint, p.Device, p.Fstype, p.Opts)
	}
}
