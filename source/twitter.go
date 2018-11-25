package source

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/ChimeraCoder/anaconda"
	"github.com/athletifit/social-network-insights/models"
)

// Twitter what we need to sort our followers etc.
type Twitter struct {
	API  *anaconda.TwitterApi
	Conf Conf
}

// NewTwitter returns a struct with initialised twitter api.
func NewTwitter(env map[string]string, conf Conf) (*Twitter, error) {
	consumerKey := env["TWITTER_CONSUMER_KEY"]
	consumerSecret := env["TWITTER_CONSUMER_SECRET"]
	accessToken := env["TWITTER_ACCESS_TOKEN"]
	tokenSecret := env["TWITTER_TOKEN_SECRET"]

	api := anaconda.NewTwitterApiWithCredentials(accessToken, tokenSecret, consumerKey, consumerSecret)

	return &Twitter{
		API:  api,
		Conf: conf,
	}, nil
}

// getFriends sort twitter friends by count followers desc.
func (ts Twitter) getFriends(v url.Values) []models.User {
	users := make([]models.User, 0, 1)
	results := ts.API.GetFriendsListAll(v)

	// iterate through the result channel
	for result := range results {
		if result.Error != nil {
			fmt.Printf("%+v", result.Error)
			break
		}

		// iterate through each result page
		for _, f := range result.Friends {
			u := models.User{
				ScreenName:     f.ScreenName,
				FollowersCount: int64(f.FollowersCount),
				URL:            f.URL,
				Name:           f.Name,
				Email:          f.Email,
			}
			users = append(users, u)
		}
		fmt.Printf("T")
	}

	return users
}

// getFollowers sort twitter followers by count of followers desc.
func (ts Twitter) getFollowers(v url.Values) []models.User {
	users := make([]models.User, 0, 1)
	results := ts.API.GetFollowersListAll(v)

	// iterate through the result channel
	for result := range results {
		if result.Error != nil {
			fmt.Printf("%+v", result.Error)
			break
		}

		// iterate through each result page
		for _, f := range result.Followers {
			u := models.User{
				ScreenName:     f.ScreenName,
				FollowersCount: int64(f.FollowersCount),
				URL:            f.URL,
				Name:           f.Name,
				Email:          f.Email,
			}
			users = append(users, u)

		}
		fmt.Printf("T")
	}
	return users
}

// GetUsers returns the users for the source.
func (ts Twitter) GetUsers(ch models.UserSetChan, wg *sync.WaitGroup) {
	defer wg.Done()

	final := make([]models.User, 0, 1)
	// fixme we dont depuplicate here now but oh well..
	for _, screenName := range ts.Conf.Usernames {
		v := ts.getDefaultURLValues(screenName)

		tFol := ts.getFollowers(v)
		tFri := ts.getFriends(v)

		final = append(tFol, tFri...)
	}

	ch <- models.NewUserSet("twitter", final)
}

func (ts Twitter) getDefaultURLValues(screenName string) url.Values {
	v := url.Values{}
	v.Add("screen_name", screenName)
	v.Add("include_user_entities", "false")
	v.Add("count", "200")
	return v
}
