package checks

import (
	"fmt"

	"github.com/math2001/gocmt/cmt"
	"github.com/shirou/gopsutil/process"
)

func Process(c *cmt.CheckResult, args map[string]interface{}) {
	name := args["name"].(string)
	psname := args["psname"].(string)

	c.AddItem(&cmt.CheckItem{
		Name:  "process_name",
		Value: name,
	})

	pids, err := process.Pids()
	if err != nil {
		panic(err)
	}

	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			panic(err)
		}

		actualpsname, err := p.Name()
		if err != nil {
			panic(err)
		}

		if actualpsname != psname {
			continue
		}

		infos, err := p.MemoryInfo()
		if err != nil {
			panic(err)
		}
		c.AddItem(&cmt.CheckItem{
			Name:        "process_memory",
			Value:       infos.RSS,
			Description: "rss",
			Unit:        "byte",
		})
		times, err := p.Times()
		if err != nil {
			panic(err)
		}
		c.AddItem(&cmt.CheckItem{
			Name:        "process_cpu",
			Value:       times.User,
			Description: "cpu time, user",
			Unit:        "seconds",
		})
		return
	}

	c.AddItem(&cmt.CheckItem{
		Name:        "process_status",
		Value:       "nok",
		Description: "ok/nok",

		IsAlert:      true,
		AlertMessage: fmt.Sprintf("check_process - %s missing (%s)", name, psname),
	})

}

func AvailProcess() {
	pids, err := process.Pids()
	if err != nil {
		panic(err)
	}
	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			panic(err)
		}
		name, err := p.Name()
		if err != nil {
			panic(err)
		}
		mempc, err := p.MemoryPercent()
		if err != nil {
			panic(err)
		}
		cpupc, err := p.CPUPercent()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%-30s(%6d): mem %.1f%% cpu %.1f%%\n", name, p.Pid, mempc, cpupc)
	}
}
