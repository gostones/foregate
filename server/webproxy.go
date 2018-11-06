package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func toUrl(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func StripPrefix(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// /p/port/
		if p := strings.SplitN(r.URL.Path, "/", 4); len(p) > 3 {
			fmt.Printf("StripPrefix: %v %v\n", p, r.URL)

			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p[3]
			h.ServeHTTP(w, r2)
		} else {
			http.NotFound(w, r)
		}
	}
}
