package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"html"

	"gopkg.in/gorp.v1"
	"github.com/gorilla/mux"
)

func FileServerRouteG(m *mux.Router, path, dir string) {
	m.PathPrefix(path).Handler(
		http.StripPrefix(path, http.FileServer(http.Dir(dir))))
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/views/index.html")
}

//fetchLocations takes in a GORP DbMap and fetches all locations in the
//locations database table into a Location slice
func fetchLocations(dbMap *gorp.DbMap) []Location {
	var locations []Location

	_, err := dbMap.Select(&locations, "SELECT * FROM locations")
	if err != nil {
		log.Fatal(err)
	}

	return locations
}

//makeLocationsRoute takes in a GORP DbMap and makes a route that uses that
//DbMap to fetch all locations in the locations table and serve them as JSON.
func makeLocationsRoute(dbMap *gorp.DbMap) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		locations := fetchLocations(dbMap)

		locationsJSON, err := json.Marshal(locations)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, "%s", locationsJSON)
	}
}

func AngularReturnError(w http.ResponseWriter, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, `{"message":%q}`, err)
}

func HandleSignup(dbMap *gorp.DbMap) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := initUser()
  	err := json.NewDecoder(r.Body).Decode(user)
		if (err != nil) {
			log.Fatal(err)
		}

		// TODO: Some checks on the input

		err = user.save(dbMap)
		if (err != nil) {
			log.Fatal(err)
		}
	}
}

func HandleSignin(dbMap *gorp.DbMap) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := initUser()
  	err := json.NewDecoder(r.Body).Decode(user)
		if (err != nil) {
			log.Fatal(err)
		}

		dbUser := initUser()
		err = dbMap.SelectOne(dbUser, "select * from users where Username=$1", user.Username)
		if (err != nil) {
			AngularReturnError(w, "Wrong username or password")
			return
		}

		if (dbUser.Password != user.Password) {
				AngularReturnError(w, "Wrong username or password")
				return
		}

		dbUser.Password = "" // Not needed anymore
		user.Password = ""
		user.Salt = ""

		fmt.Fprintf(w, "%s", dbUser.JSON())
	}
}

//AddStaticRoutes takes in a Gorilla mux Router and an alternating set of URL
//paths and directory paths and for each pair of strings, the router is given a
//FileServer Handler where the first string is the URL path and the second
//string is the directory path to serve files from.
func AddStaticRoutes(m *mux.Router, pathsAndDirs ...string) {
	for i := 0; i < len(pathsAndDirs)-1; i += 2 {
		FileServerRouteG(m, pathsAndDirs[i], pathsAndDirs[i+1])
	}
}

//initRouter takes in a GORP DbMap and initializes the router's routes while
//using the DbMap to handle database functionality.
func initRouter(dbMap *gorp.DbMap) *mux.Router {
	r := mux.NewRouter()

	//Add static routes for the public directory
	AddStaticRoutes(r, "/partials/", "public/partials",
		"/scripts/", "public/scripts", "/styles/", "public/styles",
		"/images/", "public/images")

	//Add the locations route API with makeLocationsRoute
	r.HandleFunc("/locations", makeLocationsRoute(dbMap))

	r.HandleFunc("/auth/signup", HandleSignup(dbMap)).Methods("POST")
	r.HandleFunc("/auth/signin", HandleSignin(dbMap)).Methods("POST")

	//Serve all other requests with index.html, and ultimately the front-end
	//Angular.js app.
	r.PathPrefix("/").HandlerFunc(indexRoute)

	return r
}
