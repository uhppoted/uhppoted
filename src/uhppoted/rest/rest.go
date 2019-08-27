package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type handlerfn func(*context.Context, http.ResponseWriter, *http.Request)

type handler struct {
	re *regexp.Regexp
	fn handlerfn
}

type dispatcher struct {
	ctx      *context.Context
	handlers []handler
}

func Run(ctx *context.Context) {
	d := dispatcher{
		ctx,
		make([]handler, 0),
	}

	d.Add("^/uhppote/device$", devices)

	log.Fatal(http.ListenAndServe(":8001", &d))
}

func Close() {
}

func (d *dispatcher) Add(path string, h handlerfn) {
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
			h.fn(d.ctx, w, r)
			return
		}
	}

	// Fall-through handler
	http.Error(w, "Unsupported API", http.StatusBadRequest)
}

func devices(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.Error(w, "NOT IMPLEMENTED", http.StatusNotImplemented)

	default:
		http.Error(w, fmt.Sprintf("Invalid method:%s - expected GET or POST", r.Method), http.StatusMethodNotAllowed)
	}
}
