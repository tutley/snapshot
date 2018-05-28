package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/integrii/headlessChrome"
)

var chromePath string
var listenURI string
var sleepTime string

func init() {
	// Grab the environment variables to setup this service
	// listenURL will be localhost:LISTEN_PORT with default 9999
	listenPort := getEnv("LISTEN_PORT", "9999")
	listenURI = fmt.Sprintf(":%s", listenPort)

	// This is the path to the chrome binary
	chromePath = getEnv("CHROME_PATH", `/usr/bin/chromium-browser`)

	// sleep time is the amount of time we should wait for chrome to render the page
	sleepTime = getEnv("SLEEP_TIME", "2")
}

func main() {
	// this is the webserver part, just one route '/'
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse the requested URL from the request
		u := r.URL.Query().Get("url")
		if len(u) == 0 {
			log.Println("Error: No URL was provided")
			http.Error(w, "There was no URL parameter", http.StatusBadRequest)
			return
		}
		// This function should check the URL to make sure it is sensible
		// maybe use the URL package here to do some more error checking?

		// doTheThing and return the string
		log.Println("Fetching ", u)
		resp, err := doTheThing(u)
		if err != nil {
			log.Println("Error, page load didn't work. ", err.Error())
			http.Error(w, "Something happened and the parse didn't work", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", resp)
	})
	http.ListenAndServe(listenURI, r)
}

// doTheThing runs chromium browser to fetch the page and generate a string
func doTheThing(u string) (string, error) {
	// make a new session
	headlessChrome.ChromePath = chromePath
	browser, err := headlessChrome.NewBrowser(u)
	if err != nil {
		return "", err
	}
	// Close the browser process when this func returns
	defer browser.Exit()

	// sleep while content is rendered.
	st := fmt.Sprintf("%vs", sleepTime)
	sd, err := time.ParseDuration(st)
	if err != nil {
		return "", err
	}
	time.Sleep(sd)

	// Query all the HTML from the web site
	browser.Write(`document.documentElement.outerHTML`)
	time.Sleep(time.Second)

	// loop over all the output that came from the output channel
	// and print it to the console
	var result string
	for len(browser.Output) > 0 {
		result = <-browser.Output
		//fmt.Println(<-browser.Output)
	}

	return result, err
}

// this is a helper func to fetch environment variables
func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
