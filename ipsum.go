package googleipsum

import (
	"math/rand"
	"strings"
	"time"
)

func getIpsum() (string, error) {
	var s string

	n := paragraphLength()

	for i := 0; i < n; i++ {
		s += "google "
	}
	s = strings.TrimSpace(s)

	return s, nil
}

func paragraphLength() int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)

	return n
}
