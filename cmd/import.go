package cmd

import (
	"fmt"
	"sync"

	"github.com/athletifit/social-network-insights/datastore"
	"github.com/athletifit/social-network-insights/models"
	"github.com/athletifit/social-network-insights/source"
	"github.com/joho/godotenv"
)

// ImportData import sources passed in cmd flags.
func ImportData(sourcesPtr models.ArrayFlags) error {
	ch := make(models.UserSetChan)
	wg := &sync.WaitGroup{}

	var env map[string]string
	env, err := godotenv.Read()
	if err != nil {
		return err
	}

	for _, s := range sourcesPtr {
		src, err := source.NewSource(s, env)
		if err != nil {
			fmt.Printf("%+v", err)
			continue
		}

		wg.Add(1)
		go src.GetUsers(ch, wg)
		fmt.Printf("=== %s started === \n", s)
	}

	// wait that all goroutines have sent the data to close the channel
	go func() {
		wg.Wait()
		close(ch)
	}()

	dataStore := datastore.NewRedisDataStore()

	for userSet := range ch {

		fmt.Printf("=== Done with %s ===\n", userSet.Title)
		err := dataStore.SaveUsers(userSet)
		if err != nil {
			fmt.Printf("%+v", err)
		}
	}

	fmt.Println("=== Import finished ===")
	return nil
}
