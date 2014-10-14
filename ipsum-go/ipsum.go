package googleipsum

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"
)

func getIpsum(r *http.Request) (string, error) {
	var s string

	n := paragraphLength()

	for i := 0; i < n; i++ {
		word := getWord(r)
		if word != "" {
			s += word + " "
		} else {
			s += "google "
		}
	}
	s = strings.TrimSpace(s)

	return s, nil
}

func paragraphLength() int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)

	return n
}

func getWord(r *http.Request) string {
	type Content struct {
		Date    time.Time
		Content string
	}

	c := appengine.NewContext(r)

	q := datastore.NewQuery("Ipsum")
	var words []Content
	_, err := q.GetAll(c, &words)
	if err != nil {
		c.Errorf("fetching ipsum text: %v", err)
		return ""
	}

	rand.Seed(time.Now().UnixNano())

	return words[rand.Intn(len(words))].Content
}
