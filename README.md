# Disclaimer

This is a random script that served some internal usage.
Nothing fancy to see!

Not much err checking/logging so dont panic if you see random panics. (PUN!)

# Setup

Run `docker-compose up` if you need a local redis instance.
At the moment the only datastore is redis to simplify things.

Copy `.env.example` to `.env` and fill whatever you need.

Later we should be able to split the package to offer:

- a command line tool to import data
- a rest api + a web ui to visualise/export the data

Examples:

`go run main.go --mode=import --source=instagram --source=twitter`

Will import list of users that follow you or you are following on IG/Twitter
with whatever details we can get.

`go run main.go --mode=export --source=instagram --source=twitter`

Will create a google sheet with the exported details.
There is a filter view available that allow you to sort ASC/DESC the data!

Deps with `dep`

# Sources

We can import followers from the list below;

## INSTAGRAM

If you want to use IG run `go run export-insta.go` first to create a cookie.
The cookie will be used so you dont have to login/logout all the time.

Unknown rate limit

## TWITTER

Twitter developer account (for twitter).

## FB

Turns out you cant really get anything interesting from this.
No more friends, no more subscribes/fans to a page...

Get a page or page app token.

via `https://developers.facebook.com/tools/explorer/`

## LINKEDIN

TODO

# Export

We can export to

## GOOGLE SHEET

Export takes 2 params: depth (how many levels do you want to go through (may be expensive in api calls))
baseScreenName, what to use to query first level.
