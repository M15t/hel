package tools

import (
	"fmt"
	"hel/fshare"
	"hel/redis"
	"time"
)

// FSCredentials struct
type FSCredentials struct {
	Token     string `redis:"token"`
	SessionID string `redis:"session_id"`
	UpdatedAt int64  `redis:"updated_at"`
}

// GetCreds retrieves Fshare credentials from Redis cache or Fshare API if not found in cache.
// It takes an instance of fshare.Fshare as input and returns the FSCredentials struct and an error.
// If the credentials are not found in the cache, it fetches them from the Fshare API, stores them in the cache,
// and returns the credentials. If the credentials are found in the cache, it retrieves and returns them.
// The function handles Redis operations and error checking internally.
func GetCreds(fs *fshare.Fshare) (*FSCredentials, error) {
	addr := "127.0.0.1:6379"
	key := "fs:creds"

	rd := redis.New(&redis.Config{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	if !rd.Exists(key) {
		t, err := fs.GetToken()
		if err != nil {
			return nil, err
		}

		now := time.Now().Unix()
		if err := rd.HSet(key, map[string]interface{}{
			"token":      t.Token,
			"session_id": t.SessionID,
			"updated_at": now,
		}, time.Hour*6); err != nil {
			return nil, err
		}

		return &FSCredentials{
			Token:     t.Token,
			SessionID: t.SessionID,
			UpdatedAt: now,
		}, nil
	}

	var creds FSCredentials
	if err := rd.HGetAll(key, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

// StartQueue creates and returns a channel of strings to process the results of fshare.FileResponse objects.
// It populates the channel with formatted URLs based on the Linkcode field of each fshare.FileResponse object in the results slice.
// The channel has a buffer size of 100 to handle concurrent processing efficiently.
// The function launches a goroutine to populate the channel and closes it when all results are processed.
// The returned channel can be used to asynchronously consume the generated URLs.
func StartQueue(results []*fshare.FileResponse) <-chan string {
	queue := make(chan string, 100)

	go func() {
		defer close(queue)
		for i := 0; i < len(results); i++ {
			queue <- fmt.Sprintf("https://www.fshare.vn/file/%s", results[i].Linkcode)
		}
	}()

	return queue
}
