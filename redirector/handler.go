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
	uri = config.Redirector.GetURI()
	log.Println("Connected to Redis on " + config.Redis.GetAddr())

}

type CreateRedirectURLBody struct {
	URL       string `json:"url"`
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
	new_url := fmt.Sprintf("%09s", string(base62.FormatUint(counter)))

	// Store the URL in Redis
	rdb.Set(ctx, "redirector:url:"+new_url+":l", body.URL, 0).Err()

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"url": uri + "/" + new_url,
	}
	json.NewEncoder(w).Encode(response)
	log.Printf(">> %s %s %d", body.URL, new_url, counter)
}

func GetAllRedirectURLs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	keys, err := rdb.Keys(ctx, "redirector:url:*:l").Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a list of all redirect URLs
	w.WriteHeader(http.StatusOK)
	response := make(map[string]string)
	for _, key := range keys {
		url := rdb.Get(ctx, key).Val()
		response[strings.Split(key, ":")[2]] = url
	}
	json.NewEncoder(w).Encode(response)
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
