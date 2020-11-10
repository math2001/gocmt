package main

import (
	"fmt"

	"github.com/math2001/gocmt/cmt"
)

type checkerfunction func(
	cr *cmt.CheckResult,
	args map[string]interface{},
)

func main() {
	conf := loadConf()

	checkResults := runChecks(conf)
	sendReport(conf.FrameworkSettings, checkResults)

	fmt.Println("CMT done")
}
