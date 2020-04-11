/*
AUTHOR: Patrick Archer (@pjarcher913)
DATE CREATED: 10 April 2020
Copyright (c) 2020 Patrick Archer
*/

/*
Package main is the primary entry file for program execution.
It runs a preliminary setup routine that, after successful completion, initializes the web server and page router.
*/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

// Only enable if debugging and want extra log info recorded
const DEBUG_MODE = true  // Limits functionality for dev. purposes

// Where log files are to be generated and stored
const LOG_PATH = "./logs/"

// Host port to serve on
const SERVER_PORT = ":3000"

const (
	PATH_TO_HOME_HTML = "./src/web/pages/home/home.html" // Location of home.html (used to render page)
)

// Unique id tag included into newly-generated log file names
var LOG_STAMP = "main_" + time.RFC3339


func main() {
	// Perform all preliminary setups before server goes live and starts listening for requests.
	// If prelimSetup() completes, it returns true and the program continues.
	// If it fails, it returns false and fatal error.
	if prelimSetup() {
		// Initialize route handler
		routeHandler := initRouter()
		// Begin serving the site and listening on assigned web_old port
		initWebServer(routeHandler)
	} else {
		log.Fatalln("prelimSetup() failed.")
	}
}

// initLogger() initializes Logrus and configures it for future utilization.
func initLogger() {
	file, err := os.OpenFile(LOG_PATH + LOG_STAMP + ".log", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//defer file.Close()  // Because file will be closed in GC, we can leave it open so the logger is uninterrupted
	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})

	// Minimum severity level to log, according to if the program is running in debug mode or not
	if DEBUG_MODE { log.SetLevel(log.DebugLevel) } else { log.SetLevel(log.ErrorLevel) }

	log.Infoln("Logger initialized successfully.")
}

// prelimSetup() performs all preliminary setups before server goes live and starts listening for requests.
func prelimSetup() bool {
	// Initialize logger
	initLogger()

	// Services init'd, thus prelimSetup() is complete
	log.Println("Initializing services via prelimSetup().")
	return true
}

// initRouter() initializes Mux's routing services and configures them according to the website's defined route hierarchy.
func initRouter() *mux.Router {
	log.Infoln("Executing initRouter().")

	// Init mux router object
	r := mux.NewRouter()

	// Init route handlers
	/* 404 */
	r.NotFoundHandler = http.HandlerFunc(prh_404)
	/* GETs */
	r.HandleFunc("/", prh_GET_Home).Methods("GET")

	/* POSTs */
	r.HandleFunc("/{rootParam}", prh_POST_Home).Methods("POST")

	return r
}

// initWebServer() initializes the web server and begins serving clients connecting to the pre-configured SERVER_PORT
func initWebServer(routeHandler *mux.Router) {
	log.Infoln("Executing initWebServer().")

	// Serve website and listen on configured SERVER_PORT
	// http.ListenAndServe() always returns a non-nil error, and the error is its only returned value.
	// However, http.ListenAndServe() should never return (unless error or intentional termination).
	fmt.Println("Now serving on 127.0.0.1" + SERVER_PORT)
	log.Infoln("Now serving on 127.0.0.1" + SERVER_PORT)
	err := http.ListenAndServe(SERVER_PORT, routeHandler)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
