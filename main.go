package main

import (
	"flag"
	"fmt"
	"os"
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

	flag.Parse()

	if *available {
		debugAvailables()
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
