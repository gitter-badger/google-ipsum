package googleipsum

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/user"
)

func getUserContext(w http.ResponseWriter, r *http.Request) (appengine.Context, *user.User) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, nil
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return nil, nil
	}

	return c, u
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	getUserContext(w, r)

	page := template.Must(template.ParseFiles(
		"static/admin/_base.html",
		"static/admin/index.html",
		"static/admin/main.html",
	))

	if err := page.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addWord(w http.ResponseWriter, r *http.Request) {
	getUserContext(w, r)

	if r.Method == "GET" {
		page := template.Must(template.ParseFiles(
			"static/admin/_base.html",
			"static/admin/index.html",
			"static/admin/add-word.html",
		))

		if err := page.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
