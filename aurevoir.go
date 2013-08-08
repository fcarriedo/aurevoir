package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

// Command line args def
var port = flag.Int("p", 9096, "http port to run")
var devMode = flag.Bool("dev", false, "run in development mode")

var templates *template.Template

func init() {
	// Parse all templates
	templates = template.Must(template.New("app").ParseGlob("web/tmpl/*.html"))
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func commitsHandler(w http.ResponseWriter, req *http.Request) {
	project := mux.Vars(req)["project"]

	params := map[string]interface{}{
		"project": project,
		"commits": map[string]string {
      "eee": "This is the 5th commit",
      "ddd": "This is the 4th commit",
      "ccc": "This is the 3er commit",
      "bbb": "This is the 2nd commit",
      "aaa": "This is the 1st commit",
    },
	}

  templates.ExecuteTemplate(w, "listing.html", params)
}

func commitHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	project := params["project"]
	commit := params["id"]

	fmt.Fprintf(w, "Showing the diff for %s/%s", project, commit)
}

func main() {
	flag.Parse()

	// Create the mux router
	router := mux.NewRouter()

	// Static resources - resource ending in common know web file formats
	// (css, html, jpg, etc.) get handled directly by the fileServer
	router.Handle("/{static-res:(.+\\.)(js|css|jpg|png|ico|gif)$}", http.FileServer(http.Dir("web/")))

	// Home handler
	router.HandleFunc("/", homeHandler)

	// Mapping handling
	router.HandleFunc("/{project}/commits", commitsHandler)
	router.HandleFunc("/{project}/commits/{id}", commitHandler)

	// Hook it with http pkg
	http.Handle("/", router)

	host := fmt.Sprintf(":%d", *port)
	fmt.Printf("Server up and listening on %s\n", host)
	log.Fatal(http.ListenAndServe(host, nil))
}
