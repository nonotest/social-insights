package source

import (
	"fmt"
	"sync"

	"github.com/athletifit/social-network-insights/models"
	"gopkg.in/ahmdrz/goinsta.v2"
)

// Instagram what we need to sort our followers etc.
type Instagram struct {
	API  *goinsta.Instagram
	Conf Conf
}

// NewInstagram returns a struct with initialised instagram api.
func NewInstagram(env map[string]string, conf Conf) (*Instagram, error) {
	insta, err := goinsta.Import("ig-cookie")
	if err != nil {
		return nil, err
	}

	return &Instagram{
		API:  insta,
		Conf: conf,
	}, nil
}

// getFriends sort twitter friends by count followers desc.
// TODO: make recursive ?
// Ideally we want to get friends of our friends etc (depth)
func (ig Instagram) getFriends() []models.User {
	user, err := ig.API.Profiles.ByName(ig.Conf.Usernames[0])
	if err != nil {
		fmt.Printf("%+v", err)
		return nil
	}

	friends := user.Following()
	users := make([]models.User, 0, 1)

	for friends.Next() {
		for _, user := range friends.Users {
			// we get more details from a profile search
			// will cost us more queries but its useful.
			friend, err := ig.API.Profiles.ByName(user.Username)
			if err != nil {
				fmt.Printf("%+v", err)
				continue
			}
			users = append(users, models.User{
				ScreenName:     friend.Username,
				FollowersCount: int64(friend.FollowerCount),
				Name:           friend.FullName,
				Email:          friend.Email,
				URL:            friend.ExternalURL,
			})
		}
		fmt.Printf("IG")
	}

	return users
}

// getFollowers sort ig followers by count followers desc.
func (ig Instagram) getFollowers() []models.User {

	user, err := ig.API.Profiles.ByName(ig.Conf.Usernames[0])
	if err != nil {
		fmt.Printf("%+v", err)
		return nil
	}

	followers := user.Followers()
	users := make([]models.User, 0, 1)

	for followers.Next() {
		for _, user := range followers.Users {
			follower, err := ig.API.Profiles.ByName(user.Username)
			if err != nil {
				fmt.Printf("%+v", err)
				continue
			}

			users = append(users, models.User{
				ScreenName:     follower.Username,
				FollowersCount: int64(follower.FollowerCount),
				Name:           follower.FullName,
				Email:          follower.Email,
				URL:            follower.ExternalURL,
			})
		}
		fmt.Printf("IG")
	}

	return users
}

// GetUsers returns the dataset for the source.
func (ig Instagram) GetUsers(ch models.UserSetChan, wg *sync.WaitGroup) {
	defer wg.Done()

	tFol := ig.getFollowers()
	tFri := ig.getFriends()

	final := append(tFol, tFri...)

	ch <- models.NewUserSet("instagram", final)

}
