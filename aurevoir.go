package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

// Command line args def
var port = flag.Int("p", 9096, "http port to run")
var devMode = flag.Bool("dev", false, "run in development mode")
var root = flag.String("root", ".", "The root of the repository")

var templates *template.Template

type commit struct {
	Id   string
	Msg  string
	Data string
}

type CommitReader interface {
	Commits() map[string]commit
}

func init() {
	// Parse all templates
	templates = template.Must(template.New("app").ParseGlob("web/tmpl/*.html"))
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	dirs, err := ioutil.ReadDir(*root)
	if err != nil {
		log.Fatal(err)
	}

	var projects []string
	for _, dir := range dirs {
		projects = append(projects, dir.Name())
	}

	templates.ExecuteTemplate(w, "index.html", projects)
}

func commitsHandler(w http.ResponseWriter, req *http.Request) {
	project := mux.Vars(req)["project"]

	cr := newCommitReader(*root + "/" + project)

	params := map[string]interface{}{
		"project": project,
		"commits": cr.Commits(),
	}

	templates.ExecuteTemplate(w, "commits.html", params)
}

func commitHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	project := params["project"]
	commitId := params["id"]

	commit := getCommit(project, commitId)
	commit.Msg = "This is a message"

	templates.ExecuteTemplate(w, "commit.html", commit)
}

func getCommit(project, commitId string) commit {
	cmd := exec.Command("git", "show", commitId)
	cmd.Dir = *root + "/" + project
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("Couldn't get commit " + commitId)
	}

	return commit{Id: commitId, Data: "\n" + string(out)}
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

// This has to go in a git.go file of package aurevoir
type gitCommitReader struct {
	baseDir string
}

func newCommitReader(dir string) CommitReader {
	return &gitCommitReader{dir}
}

func (cr *gitCommitReader) Commits() map[string]commit {
	cmd := exec.Command("git", "log", "--oneline", "-20", "--no-merges")
	cmd.Dir = cr.baseDir
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	commits := make(map[string]commit)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.SplitN(scanner.Text(), " ", 2)
		c := commit{Id: line[0], Msg: line[1]}
		commits[c.Id] = c
	}

	return commits
}
