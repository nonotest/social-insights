package datastore

import "github.com/athletifit/social-network-insights/models"

type DataStore interface {
	LoadUsers(source string) (*models.UserMap, error)
	SaveLastTwitterCursor(cursor int64)
	SaveUsers(us *models.UserSet) error
}

func NewDataStore(store string) DataStore {
	switch store {
	case "redis":
		return NewRedisDataStore()
	}

	return nil
}
