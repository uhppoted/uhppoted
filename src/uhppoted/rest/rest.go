package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"uhppote"
)

type handlerfn func(context.Context, http.ResponseWriter, *http.Request)

type handler struct {
	re *regexp.Regexp
	fn handlerfn
}

type dispatcher struct {
	u        *uhppote.UHPPOTE
	handlers []handler
}

func Run(u *uhppote.UHPPOTE) {
	d := dispatcher{
		u,
		make([]handler, 0),
	}

	d.Add("^/uhppote/device$", devices)
	d.Add("^/uhppote/device/[0-9]+$", device)

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
			ctx := context.WithValue(context.Background(), "uhppote", d.u)
			h.fn(ctx, w, r)
			return
		}
	}

	// Fall-through handler
	http.Error(w, "Unsupported API", http.StatusBadRequest)
}

func devices(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
		GetDevices(u, w, r)

	default:
		http.Error(w, fmt.Sprintf("Invalid method:%s - expected GET", r.Method), http.StatusMethodNotAllowed)
	}
}

func device(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
		url := r.URL.Path
		matches := regexp.MustCompile("^/uhppote/device/([0-9]+)$").FindStringSubmatch(url)
		deviceId, err := strconv.ParseUint(matches[1], 10, 32)
		if err != nil {
			http.Error(w, "Error reading request", http.StatusInternalServerError)
			return
		}

		GetDevice(uint32(deviceId), u, w, r)

	default:
		http.Error(w, fmt.Sprintf("Invalid method:%s - expected GET", r.Method), http.StatusMethodNotAllowed)
	}
}
