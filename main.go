package googleipsum

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bmizerany/pat"
	"github.com/yosssi/gcss"
)

// define the routes during package initization.
func init() {
	router := pat.New()
	// handle asset paths
	router.Get("/css/:file", http.HandlerFunc(compileCSS))

	// handle application paths
	router.Get("/admin", http.HandlerFunc(adminHandler))
	router.Get("/admin/word/add", http.HandlerFunc(addWord))
	router.Post("/generate", http.HandlerFunc(generateIpsum))
	router.Get("/", http.HandlerFunc(rootHandler))
	http.Handle("/", router)
}

// compileCSS gets the CSS name from the URL, loads the appropate GCSS file,
// then serves the client the compiled version.
func compileCSS(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get(":file")
	if file == "" {
		log.Println("did not get a filename")
		http.Error(w, "no CSS file requested", http.StatusInternalServerError)
	}

	// convert the .css extension to .gcss, and build out path to the file
	f := gcss.Path(file)
	f = fmt.Sprintf("static/css/%s", f)
	log.Printf("GCSS file: %s\n", f)

	// read in the GCSS file
	css, err := os.Open(f)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// close out the file resource once done
	defer func() {
		if err := css.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	// set the content type header so browsers will know how to handle it
	w.Header().Set("Content-Type", "text/css")

	// build out the CSS and serve it to the browser
	_, err = gcss.Compile(w, css)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// rootHandler handles the root path, and catch all for unmatched routes
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// if the path is not the root path then it is a 404
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	}

	// load the template files used for this page
	page := template.Must(template.ParseFiles(
		"static/_base.html",
		"static/index.html",
	))

	// render the template files and serve the page
	if err := page.Execute(w, nil); err != nil {
		http.Error(w, "failed to load page", http.StatusInternalServerError)
	}
}

func generateIpsum(w http.ResponseWriter, r *http.Request) {
	type content struct {
		// Ipsum contains the generated ipsum text
		Ipsum []string
	}
	t := content{}

	s := r.FormValue("paragraph")
	n, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	s = r.FormValue("lead")
	if s == "" {
		s = "false"
	}

	l, err := strconv.ParseBool(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	for i := 0; i < n; i++ {
		var paragraph string
		if i == 0 {
			if l {
				paragraph += "Google ipsum dolor sit amet "
			}
		}
		// TODO generate ipsum text
		s, _ := getIpsum()
		paragraph = paragraph + s

		t.Ipsum = append(t.Ipsum, paragraph)
	}

	page := template.Must(template.ParseFiles(
		"static/_base.html",
		"static/index.html",
	))

	if err = page.Execute(w, t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
