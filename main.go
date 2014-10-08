package googleipsum

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/bmizerany/pat"
	"github.com/yosssi/gcss"
)

// define the routes during package initization.
func init() {
	router := pat.New()
	// handle asset paths
	router.Get("/css/:file", http.HandlerFunc(compileCSS))

	// handle application paths
	router.Get("/", http.HandlerFunc(rootHandler))
	http.Handle("/", router)
}

func compileCSS(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get(":file")
	if file == "" {
		log.Println("did not get a filename")
		http.Error(w, "no CSS file requested", http.StatusInternalServerError)
	}

	f := gcss.Path(file)
	f = fmt.Sprintf("static/css/%s", f)
	log.Printf("GCSS file: %s\n", f)

	css, err := os.Open(f)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	defer func() {
		if err := css.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	w.Header().Set("Content-Type", "text/css")
	_, err = gcss.Compile(w, css)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	page := template.Must(template.ParseFiles(
		"static/_base.html",
		"static/index.html",
	))

	if err := page.Execute(w, nil); err != nil {
		http.Error(w, "failed to load page", http.StatusInternalServerError)
		return
	}
}
