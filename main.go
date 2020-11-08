package main

import (
	"fmt"
	"sync"

	"github.com/math2001/gocmt/cmt"
)

func main() {
	conf := loadConf()
	checkResults := runChecks(conf)
	var wg sync.WaitGroup
	for checkResult := range checkResults {
		wg.Add(1)
		go func(checkResult *cmt.CheckResult) {
			defer wg.Done()
			report(checkResult)
		}(checkResult)
	}
	wg.Wait()
	fmt.Println("CMT done")
}
