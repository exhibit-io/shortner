package main

import (
	"log"
	"net/http"

	"github.com/exhibit-io/redirector"
	"github.com/exhibit-io/redirector/config"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	config := config.LoadConfig()

	redirector.Init(config)

	// Use a Router to handle requests
	router := httprouter.New()

	// POST new redirect link
	router.POST("/", redirector.CreateRedirectURL)

	// GET all redirect links
	router.GET("/", redirector.GetAllRedirectURLs)

	// GET redirect link by path
	router.GET("/:url", redirector.HandleURLRedirection)

	// Setup middlewares.  For this we're basically adding:
	//	- Support for CORS to make JSONP work.
	handler := cors.Default().Handler(router)

	log.Println("Starting HTTP server on:", config.Redirector.GetAddr())
	log.Fatal(http.ListenAndServe(config.Redirector.GetAddr(), handler))
}
