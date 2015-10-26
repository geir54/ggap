package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"html"
	"text/template"

	"gopkg.in/gorp.v1"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func FileServerRouteG(m *mux.Router, path, dir string) {
	m.PathPrefix(path).Handler(
		http.StripPrefix(path, http.FileServer(http.Dir(dir))))
}

func indexRoute(dbMap *gorp.DbMap) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("public/views/index.html")
		if err != nil {
			log.Fatal(err)
		}

		session, err := store.Get(r, "sessions")
		if err != nil {
			log.Println(err)
			return
		}

		type Response struct {
		    User   string
		}
		var user Response

		dbUser := initUser()
		err = dbMap.SelectOne(dbUser, "select * from users where Id=$1", session.Values["id"])
		if (err != nil) {
			user = Response{User : "\"\""}
		} else {
			dbUser.Password = "" // Not needed anymore

			user = Response{User : dbUser.JSON()}
		}

	  err =	t.Execute(w, user)
		if err != nil {
			log.Fatal(err)
		}
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
		err = dbMap.SelectOne(dbUser, "select * from users where username=$1", user.Username)
		if (err != nil) {
			AngularReturnError(w, "Wrong username or password")
			return
		}

		if (!dbUser.checkPassword(user.Password)) {
				AngularReturnError(w, "Wrong username or password")
				return
		}

		session, err := store.Get(r, "sessions")
		if err != nil {
			log.Println(err)
			return
		}
		session.Values["id"] = dbUser.Id
		session.Save(r, w)

		dbUser.Password = "" // Not needed anymore
		user.Password = ""

		fmt.Fprintf(w, "%s", dbUser.JSON())
	}
}

func HandleSignout() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "sessions")
			if err != nil {
				log.Println(err)
				return
			}

			delete(session.Values, "id")
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusFound)
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

func serveSingle(filename string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filename)
	}
}

//initRouter takes in a GORP DbMap and initializes the router's routes while
//using the DbMap to handle database functionality.
func initRouter(dbMap *gorp.DbMap) *mux.Router {
	r := mux.NewRouter()

	//Add static routes for the public directory
	AddStaticRoutes(r, "/lib/", "public/lib",
		"/modules/", "public/modules")

	r.HandleFunc("/robots.txt", serveSingle("public/robots.txt"))
	r.HandleFunc("/application.js", serveSingle("public/application.js"))
	r.HandleFunc("/config.js", serveSingle("public/config.js"))
	r.HandleFunc("/config.js", serveSingle("public/humans.txt"))

	r.HandleFunc("/auth/signup", HandleSignup(dbMap)).Methods("POST")
	r.HandleFunc("/auth/signin", HandleSignin(dbMap)).Methods("POST")
	r.HandleFunc("/auth/signout", HandleSignout())

	//Serve all other requests with index.html, and ultimately the front-end
	//Angular.js app.
	r.PathPrefix("/").HandlerFunc(indexRoute(dbMap))

	return r
}
