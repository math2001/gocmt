package main

import (
	"flag"
	"fmt"

	"github.com/math2001/gocmt/cmt"
)

type checkerfunction func(
	c *cmt.CheckResult,
	args map[string]interface{},
)

func main() {

	listAvailable := flag.Bool(
		"available", false,
		"display available entries found for each checks (manual run on target)")

	flag.Parse()

	conf := loadConf()

	if *listAvailable {
		conf.FrameworkSettings.Available = true
	}

	checkResults := runChecks(conf)
	sendReport(conf.FrameworkSettings, checkResults)

	fmt.Println("CMT done")
}
