package source

import (
	"sync"

	"github.com/athletifit/social-network-insights/models"
)

// MockSource is a mock source that we use to dev.
type MockSource struct {
	Name string
}

// NewMockSource returns a new MockSource struct.
func NewMockSource(name string) (*MockSource, error) {
	return &MockSource{Name: name}, nil
}

func (m MockSource) getTestUsers() []models.User {
	friends := []models.User{
		models.User{
			ScreenName:     "nonotest",
			FollowersCount: 100,
		}, models.User{
			ScreenName:     "paul",
			FollowersCount: 200,
		}, models.User{
			ScreenName:     "jerome",
			FollowersCount: 50,
		}, models.User{
			ScreenName:     "lol",
			FollowersCount: 1,
		}, models.User{
			ScreenName:     "lolllll",
			FollowersCount: 2,
		}}
	return friends
}

// GetUsers returns the dataset for the source.
func (m MockSource) GetUsers(ch models.UserSetChan, wg *sync.WaitGroup) {
	defer wg.Done()

	tFol := m.getTestUsers()

	ch <- models.NewUserSet("mock", tFol)
}
