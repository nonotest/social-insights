package source

import (
	"errors"
	"sync"

	"github.com/athletifit/social-network-insights/models"
)

// TODO: Remove hardcoded strings.

// Source an interface that defines various methods for our source.
type Source interface {
	GetUsers(ch models.UserSetChan, wg *sync.WaitGroup)
}

// NewSource returns a new source.
func NewSource(srcName string, env map[string]string) (Source, error) {
	switch srcName {
	case "instagram":
		return NewInstagram(env, Conf{[]string{"stocksinplayau"}, 1})
	case "facebook":
		return NewFacebook(env, Conf{[]string{env["FACEBOOK_PAGE_NAME"]}, 1})
	case "twitter":
		conf := Conf{[]string{env["TWITTER_USERNAME"], "GregOxfordPG"}, 1}
		return NewTwitter(env, conf)
	case "mock":
		return NewMockSource("mock")
	}

	return nil, errors.New(srcName + "Source not found")
}

// Conf is a struct that we will use to pass various configuration elements
// for a source.
type Conf struct {
	Usernames []string
	Depth     int64
}
