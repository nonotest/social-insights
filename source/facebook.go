package source

import (
	"sync"

	"github.com/athletifit/social-network-insights/models"
	fb "github.com/huandu/facebook"
)

// Facebook what we need to sort our followers etc.
type Facebook struct {
	API  *fb.Session
	Conf Conf
}

// NewFacebook returns a struct with initialised instagram api.
func NewFacebook(env map[string]string, cnf Conf) (*Facebook, error) {
	// Get Facebook access token.

	/// Create a global App var to hold app id and secret.
	var globalApp = fb.New(env["FACEBOOK_APP_ID"], env["FACEBOOK_APP_SECRET"])

	// Facebook asks for a valid redirect uri when parsing signed request.
	// It's a new enforced policy starting as of late 2013.
	globalApp.RedirectUri = "http://stocksinplay.com/"

	token := env["SHORT_LIVED_FACEBOOK_TOKEN"]
	// If there is another way to get decoded access token,
	// this will return a session created directly from the token.
	session := globalApp.Session(token)

	// This validates the access token by ensuring that the current user ID is properly returned. err is nil if token is valid.
	err := session.Validate()

	if err != nil {
		return nil, err
	}

	return &Facebook{
		API:  session,
		Conf: cnf,
	}, nil
}

// GetUsers gets facebook users.
func (fbk Facebook) GetUsers(ch models.UserSetChan, wg *sync.WaitGroup) {
	defer wg.Done()

	users := make([]models.User, 0, 1)
	res, _ := fbk.API.Get("/me/insights", nil)

	// create a paging structure.
	paging, _ := res.Paging(fbk.API)
	var allResults []fb.Result

	// append first page of results to slice of Result
	allResults = append(allResults, paging.Data()...)

	for {
		// get next page.
		noMore, err := paging.Next()
		if err != nil {
			panic(err)
		}
		if noMore {
			// No more results available
			break
		}
		// append current page of results to slice of Result
		allResults = append(allResults, paging.Data()...)
	}

	ch <- models.NewUserSet("facebook", users)
}
