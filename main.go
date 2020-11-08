package main

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	conf := loadConf()
	checkResults := runChecks(conf)
	var g errgroup.Group
	for checkResult := range checkResults {
		g.Go(func() {
			report(checkResult)
		})
	}
	err := g.Wait()
	if err != nil {
		fmt.Printf("error(s) reporting: %s\n", err)
	}
	fmt.Println("CMT done")
}
