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
	Commits() (map[string]commit, error)
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
	commits, err := cr.Commits()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	params := map[string]interface{}{
		"project": project,
		"commits": commits,
	}

	templates.ExecuteTemplate(w, "commits.html", params)
}

func commitHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	project := params["project"]
	commitId := params["id"]

	commit, err := getCommit(project, commitId)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	templates.ExecuteTemplate(w, "commit.html", commit)
}

func getCommit(project, commitId string) (commit, error) {
	cmd := exec.Command("git", "show", commitId, "--oneline")
	cmd.Dir = *root + "/" + project
	out, err := cmd.Output()
	if err != nil {
		return commit{}, fmt.Errorf("'%s' doesn't appear to be a valid git repo", project)
	}

	data := strings.SplitN(string(out), "\n", 2)
	commitLine := strings.SplitN(data[0], " ", 2)

	return commit{Id: commitLine[0], Msg: commitLine[1], Data: "\n" + data[1]}, nil
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

func (cr *gitCommitReader) Commits() (map[string]commit, error) {
	cmd := exec.Command("git", "log", "--oneline", "-20", "--no-merges")
	cmd.Dir = cr.baseDir
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("'%s' doesn't appear to be a valid git repo", cr.baseDir)
	}

	commits := make(map[string]commit)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.SplitN(scanner.Text(), " ", 2)
		c := commit{Id: line[0], Msg: line[1]}
		commits[c.Id] = c
	}

	return commits, nil
}
