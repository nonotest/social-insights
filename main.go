package main

import (
	"flag"
	"fmt"

	"github.com/athletifit/social-network-insights/cmd"
	"github.com/athletifit/social-network-insights/models"
)

func main() {
	modePtr := flag.String("mode", "import", "Starts an import of your social data")
	var sourcesPtr models.ArrayFlags
	flag.Var(&sourcesPtr, "source", "Social Sources to use")
	flag.Parse()

	if *modePtr == "export" {
		cmd.ExportData(sourcesPtr)
		return
	}

	err := cmd.ImportData(sourcesPtr)
	if err != nil {
		fmt.Printf("%+v", err)
	}
}
