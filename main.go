package main

import (
	"fmt"

	"github.com/math2001/gocmt/cmt"
)

type checkerfunction func(
	cr *cmt.CheckResult,
	global_settings map[string]interface{},
	check_settings map[string]interface{},
)

func main() {
	conf := loadConf()

	checkResults := runChecks(conf)
	sendReports(conf.FrameworkSettings, checkResults)

	fmt.Println("CMT done")
}
