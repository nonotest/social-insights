package cmd

import (
	"fmt"

	"github.com/athletifit/social-network-insights/datastore"
	"github.com/athletifit/social-network-insights/export"
	"github.com/athletifit/social-network-insights/models"
)

func ExportData(sourcesPtr models.ArrayFlags) {
	store := datastore.NewRedisDataStore()

	// fixme change to users
	sets := make([]models.DataSet, 0, 1)

	for _, s := range sourcesPtr {
		users, err := store.LoadUsers(s)
		if err != nil {
			fmt.Printf("%+v", err)
			continue
		}
		set := models.NewDataSet(s, *users)
		sets = append(sets, set)
	}

	// later, what type of exporter?
	ex := export.NewSheetExporter()

	document := export.NewDocument("Influencers", sets)
	ex.Export(document)
}
