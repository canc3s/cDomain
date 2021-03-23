package main

import (
	"github.com/canc3s/cDomain/internal/runner"
)

func main() {
	options := runner.ParseOptions()

	runner.RunEnumeration(options)
}

