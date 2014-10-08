package googleipsum

import (
	"html/template"
	"net/http"

	"github.com/bmizerany/pat"
)

// define the routes during package initization.
func init() {
	router := pat.New()
	// handle application paths
	router.Get("/", http.HandlerFunc(rootHandler))
	http.Handle("/", router)
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
