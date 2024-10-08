package redirector

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
	"github.com/jxskiss/base62"

	// Import the config package
	"github.com/exhibit-io/redirector/config"
)

var (
	ctx = context.Background()
	rdb *redis.Client
	uri string
)

func Init(config *config.Config) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.Redis.GetAddr(), // Redis server address
		Password: config.Redis.Password,
	})
	if rdb.Ping(ctx).Err() != nil {
		log.Fatal("Failed to connect to Redis")
		panic("Failed to connect to Redis")
	}
	uri = config.Redirector.PublicURL
	log.Println("Connected to Redis on " + config.Redis.GetAddr())

}

type CreateRedirectURLBody struct {
	URL       string `json:"url"`
	ExpiresIn int64  `json:"expiresIn"`
}

type RedirectObject struct {
	Fragment  string `json:"fragment"`
	URL       string `json:"url"`
	Original  string `json:"original"`
	ExpiresIn int64  `json:"expiresIn"`
}

func CreateRedirectURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body CreateRedirectURLBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	counter := uint64(rdb.IncrBy(ctx, "redirector:counter", 1).Val())
	fragment := fmt.Sprintf("%09s", string(base62.FormatUint(counter)))

	// Store the URL in Redis
	rdb.Set(ctx, "redirector:url:"+fragment+":l", body.URL, 0).Err()

	// Respond with a success message
	response := RedirectObject{
		Fragment:  fragment,
		URL:       uri + "/" + fragment,
		Original:  body.URL,
		ExpiresIn: body.ExpiresIn,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Could not encode response: %v", err), http.StatusInternalServerError)
	}
	log.Printf(">> %s %s %d", body.URL, response.Original, counter)
}

func GetAllRedirectURLs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	keys, err := rdb.Keys(ctx, "redirector:url:*:l").Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a list of all redirect URLs
	response := []RedirectObject{}
	for _, key := range keys {
		url := rdb.Get(ctx, key).Val()
		fragment := strings.Split(key, ":")[2]
		response = append(response, RedirectObject{
			Fragment:  fragment,
			URL:       uri + "/" + fragment,
			Original:  url,
			ExpiresIn: -1,
		})
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Could not encode response: %v", err), http.StatusInternalServerError)
	}
}

func HandleURLRedirection(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := rdb.Get(ctx, "redirector:url:"+ps.ByName("url")+":l").Val()
	if url == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Count visits to this url
	visits := rdb.Incr(ctx, "redirector:url:"+ps.ByName("url")+":v").Val()

	http.Redirect(w, r, url, http.StatusFound)
	log.Printf("<< %s: %s %d", ps.ByName("url"), url, visits)
}
