package rest

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

type wrapper struct {
	mux *http.ServeMux
}

func Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/uhppote/simulator/", simulator)
	mux.HandleFunc("/", unsupported)

	log.Fatal(http.ListenAndServe(":8080", wrapper{mux}))
}

func (h wrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	h.mux.ServeHTTP(w, r)
}

func unsupported(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unsupported API", http.StatusBadRequest)
}

func simulator(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Greeting earthling, %q", html.EscapeString(r.URL.Path))
}
