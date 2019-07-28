package rest

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"uhppote-simulator/simulator"
)

type handler struct {
	re *regexp.Regexp
	fn http.HandlerFunc
}

type dispatcher struct {
	handlers []handler
}

func Run(simulators []*simulator.Simulator) {
	d := new(dispatcher)

	d.Add("^/uhppote/simulator/[0-9]+/door/[1-4]/swipe$", swipe)

	log.Fatal(http.ListenAndServe(":8080", d))
}

func (d *dispatcher) Add(path string, h http.HandlerFunc) {
	re := regexp.MustCompile(path)
	d.handlers = append(d.handlers, handler{re, h})
}

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// CORS pre-flight request ?
	if r.Method == http.MethodOptions {
		return
	}

	// Dispatch to handler
	url := r.URL.Path
	for _, h := range d.handlers {
		if h.re.MatchString(url) {
			h.fn(w, r)
			return
		}
	}

	// Fall-through handler
	http.Error(w, "Unsupported API", http.StatusBadRequest)
}

func swipe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("Invalid method:%s - expected POST", r.Method), http.StatusBadRequest)
		return
	}

	url := r.URL.Path
	matches := regexp.MustCompile("^/uhppote/simulator/([0-9]+)/door/([1-4])/swipe$").FindStringSubmatch(url)
	deviceId := matches[1]
	door := matches[2]

	json, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(deviceId)
	fmt.Println(door)
	fmt.Println(json)
	fmt.Fprintf(w, "Greeting earthling, %q", html.EscapeString(r.URL.Path))
}
