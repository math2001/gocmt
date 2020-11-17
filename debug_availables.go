package main

import (
	"fmt"
	"strings"

	"github.com/math2001/gocmt/checks"
)

var debugAvailable = map[string]debugAvailableFunction{
	"mount":   checks.AvailMounts,
	"process": checks.AvailProcess,
}

func debugAvailables() {
	for name, fn := range debugAvailable {
		fmt.Println(name)
		fmt.Println(strings.Repeat("=", stdoutReportWidth))
		fn()
	}
}
