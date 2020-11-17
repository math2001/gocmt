package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/math2001/gocmt/cmt"
)

type checkerFunction func(
	c *cmt.CheckResult,
	args map[string]interface{},
)

type debugAvailableFunction func()

func main() {

	available := flag.Bool(
		"available", false,
		"display available entries found for each checks (manual run on target)")
	listchecks := flag.Bool(
		"listchecks", false,
		"display available checks")

	flag.Parse()

	if *available {
		debugAvailables()
		return
	}

	if *listchecks {
		// print all the names of the checks, sorted
		names := make([]string, len(allchecks))
		i := 0
		for name := range allchecks {
			names[i] = name
			i++
		}

		sort.Slice(names, func(i, j int) bool {
			return names[i] < names[j]
		})
		lineLength := 0
		for _, name := range names {
			if lineLength > 80 {
				lineLength = 0
				fmt.Println()
			}
			lineLength += len(name)
			fmt.Print(name, " ")
		}
		if lineLength != 0 {
			fmt.Println()
		}

		return
	}

	conf := loadConf()
	go func() {
		time.Sleep(50 * time.Second)
		fmt.Print("execution took too long, exiting")
		os.Exit(1)
	}()

	checkResults := runChecks(conf)
	sendReport(conf.FrameworkSettings, checkResults)

	fmt.Println("CMT done")
}
