package models

// User is a struct that holds a user's details.
type User struct {
	Email          string `redis:"email"`
	Name           string `redis:"name"`
	ScreenName     string `redis:"screenName"`
	FollowersCount int64  `redis:"followersCount"`
	URL            string `redis:"url"`
}

type UserMap map[string]User
