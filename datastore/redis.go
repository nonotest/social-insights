package datastore

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/athletifit/social-network-insights/models"
	"github.com/garyburd/redigo/redis"
)

// TODO: look into TLS using stunnel.
// Spiped

const (
	// redisSourcesIndex the index of the sources db.
	redisSourcesIndex = 0

	// RedisSourcesPoolName is the name of the sources pool.
	RedisSourcesPoolName = "sources"
)

var pools map[string]*redis.Pool
var once = sync.Once{}

// redisDBs our maps with the mapping redis db <-> index.
var redisDBs = map[string]int{
	RedisSourcesPoolName: redisSourcesIndex,
}

// Pool wraps a new Pool function in a once.Do to only do it once...
func Pool() map[string]*redis.Pool {
	once.Do(newPools)
	return pools
}

// Init our redis pools with the different dbs.
func newPools() {
	addr := "localhost:6379" // fix later use env.
	pools = make(map[string]*redis.Pool, len(redisDBs))

	for k, v := range redisDBs {
		pools[k] = &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: (func(index int) func() (redis.Conn, error) {
				connDialer := func() (redis.Conn, error) {
					// pwd := redis.DialPassword(redisPwd)
					c, err := redis.Dial("tcp", addr)
					if err != nil {
						return nil, err
					}

					if _, err := c.Do("SELECT", strconv.Itoa(index)); err != nil {
						c.Close()
						return nil, err
					}

					return c, nil
				}
				return connDialer
			})(v),
		}
	}
}

type RedisDataStore struct {
	pools map[string]*redis.Pool
}

// NewRedisDataStore returns a new redis data store.
func NewRedisDataStore() *RedisDataStore {
	pools := Pool()
	return &RedisDataStore{
		pools: pools,
	}
}

// SaveLastTwitterCursor todo when we do deeper searches on twitter.
func (rds *RedisDataStore) SaveLastTwitterCursor(cursor int64) {
	fmt.Printf("Save cursor: %d", cursor)
}

// LoadUsers loads user from a redis hash into a user map.
func (rds *RedisDataStore) LoadUsers(source string) (*models.UserMap, error) {
	conn := rds.pools[RedisSourcesPoolName].Get()
	defer conn.Close()

	values, err := redis.StringMap(conn.Do("HGETALL", source))
	if err != nil {
		return nil, err
	}
	if err != nil && err != redis.ErrNil {
		return nil, err
	}

	users := make(models.UserMap, 0)
	for k, v := range values {
		var u models.User
		json.Unmarshal([]byte(v), &u)
		users[k] = u
	}

	return &users, nil
}

// SaveUsers persits the users of a set in a redis hash.
func (rds *RedisDataStore) SaveUsers(userSet *models.UserSet) error {
	conn := rds.pools[RedisSourcesPoolName].Get()
	defer conn.Close()

	// convert to something good for redis...
	usersMap := make(map[string][]byte, 0)
	for _, u := range userSet.Users {
		if _, ok := usersMap[u.ScreenName]; !ok {
			b, _ := json.Marshal(&u)
			usersMap[u.ScreenName] = b
		}
	}

	// only if we have users to save.
	if len(usersMap) > 0 {
		_, err := conn.Do("HMSET", redis.Args{userSet.Title}.AddFlat(usersMap)...)
		if err != nil {
			return err
		}
	}

	return nil
}
