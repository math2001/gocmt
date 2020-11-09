package main

import (
	"fmt"

	"github.com/math2001/gocmt/cmt"
)

type checkerfunction func(map[string]interface{}, map[string]interface{}) *cmt.CheckResult

func main() {
	conf := loadConf()

	checkResults := runChecks(conf)
	sendReports(conf.FrameworkSettings, checkResults)

	fmt.Println("CMT done")
}
