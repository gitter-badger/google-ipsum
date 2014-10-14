package googleipsum

import (
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type Word struct {
	Date    time.Time
	Content string
}

func ipsumKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Ipsum", "default_ipsum", 0, nil)
}

func getUserContext(r *http.Request) (appengine.Context, *user.User) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	return c, u
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	c, u := getUserContext(r)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	q := datastore.NewQuery("Ipsum").
		Order("-Date").
		Limit(5)

	var ipsum []Word
	_, err := q.GetAll(c, &ipsum)
	if err != nil {
		c.Errorf("fetching ipsum text: %v", err)
		return
	}

	page := template.Must(template.ParseFiles(
		"static/admin/_base.html",
		"static/admin/index.html",
		"static/admin/main.html",
	))

	if err := page.Execute(w, ipsum); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addWord(w http.ResponseWriter, r *http.Request) {
	c, u := getUserContext(r)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

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
	} else {
		word := Word{
			Content: r.FormValue("newWord"),
			Date:    time.Now(),
		}

		key := datastore.NewIncompleteKey(c, "Ipsum", ipsumKey(c))
		_, err := datastore.Put(c, key, &word)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusFound)
	}
}
