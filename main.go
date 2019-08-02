package main

import (
	"fmt"

	cmd "github.com/elastic/ece-support-diagnostics/cmd/ece-support-diagnostics"
)

var Version string
var Build string

func main() {
	fullVersion := fmt.Sprintf("%s (%s)", Version, Build)
	cmd.Execute(fullVersion)
}
