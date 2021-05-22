package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

var methods = []string{
	"/view/",
	"/recomTap/",
	"/addToBasket/",
	"/recomAddToBasket/",
	"/order/",
	"/visitorEvents/",
	"/categoryView/",
	"/search/",
	"/setContact/",
	"/recomBlockViewed/",
}

func main() {
	for _, method := range methods {
		http.HandleFunc(method, func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte("OK"))
		})
	}

	err := http.ListenAndServe(":80", logRequest(http.DefaultServeMux))
	if err != nil {
		panic(err)
	}
}

func logRequest(httpMux *http.ServeMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			log.Printf("%s %s %s\n%s\n\n", r.RemoteAddr, r.Method, r.URL, body)
			httpMux.ServeHTTP(w, r)
		}
	}
}
