package checks

import (
	"github.com/math2001/gocmt/cmt"
	"github.com/shirou/gopsutil/mem"
)

func Swap(c *cmt.CheckResult, args map[string]interface{}) {
	swapMem, err := mem.SwapMemory()
	if err != nil {
		panic(err)
	}
	c.AddItem(&cmt.CheckItem{
		Name:        "swap_percent",
		Value:       swapMem.UsedPercent,
		Description: "Swap used (percent)",
		Unit:        "%",
	})
	c.AddItem(&cmt.CheckItem{
		Name:        "swap_used",
		Value:       swapMem.Used,
		Description: "Swap used (bytes)",
		Unit:        "bytes",
	})
}
