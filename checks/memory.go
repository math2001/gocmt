package checks

import (
	"github.com/math2001/gocmt/cmt"
	"github.com/shirou/gopsutil/mem"
)

func Memory(c *cmt.CheckResult, args map[string]interface{}) {
	virtMem, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}
	c.AddItem(&cmt.CheckItem{
		Name:        "memory_percent",
		Value:       virtMem.UsedPercent,
		Description: "Memory used (percent)",
		Unit:        "%",
	})
	c.AddItem(&cmt.CheckItem{
		Name:        "memory_used",
		Value:       virtMem.Used,
		Description: "Memory used (bytes)",
		Unit:        "bytes",
	})
	c.AddItem(&cmt.CheckItem{
		Name:        "memory_available",
		Value:       virtMem.Available,
		Description: "Memory available (bytes)",
		Unit:        "bytes",
	})
}
