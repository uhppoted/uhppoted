package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"uhppote"
)

type handlerfn func(context.Context, http.ResponseWriter, *http.Request)

type handler struct {
	re     *regexp.Regexp
	method string
	fn     handlerfn
}

type dispatcher struct {
	uhppote  *uhppote.UHPPOTE
	log      *log.Logger
	handlers []handler
}

func Run(u *uhppote.UHPPOTE, l *log.Logger) {
	d := dispatcher{
		u,
		l,
		make([]handler, 0),
	}

	d.Add("^/uhppote/device$", http.MethodGet, getDevices)
	d.Add("^/uhppote/device/[0-9]+$", http.MethodGet, getDevice)
	d.Add("^/uhppote/device/[0-9]+/status$", http.MethodGet, getStatus)
	d.Add("^/uhppote/device/[0-9]+/time$", http.MethodGet, getTime)
	d.Add("^/uhppote/device/[0-9]+/time$", http.MethodPut, setTime)
	d.Add("^/uhppote/device/[0-9]+/door/[1-4]/delay$", http.MethodGet, getDoorDelay)
	d.Add("^/uhppote/device/[0-9]+/door/[1-4]/delay$", http.MethodPut, setDoorDelay)
	d.Add("^/uhppote/device/[0-9]+/card$", http.MethodGet, getCards)

	log.Fatal(http.ListenAndServe(":8001", &d))
}

func Close() {
}

func (d *dispatcher) Add(path string, method string, h handlerfn) {
	re := regexp.MustCompile(path)
	d.handlers = append(d.handlers, handler{re, method, h})
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
		if h.re.MatchString(url) && r.Method == h.method {
			ctx := context.WithValue(context.Background(), "uhppote", d.uhppote)
			ctx = context.WithValue(ctx, "log", d.log)
			ctx = parse(ctx, r)
			h.fn(ctx, w, r)
			return
		}
	}

	// Fall-through handler
	http.Error(w, "Unsupported API", http.StatusBadRequest)
}

func parse(ctx context.Context, r *http.Request) context.Context {
	url := r.URL.Path

	matches := regexp.MustCompile("^/uhppote/device/([0-9]+)(?:$|/.*$)").FindStringSubmatch(url)
	if matches != nil {
		deviceId, err := strconv.ParseUint(matches[1], 10, 32)
		if err == nil {
			ctx = context.WithValue(ctx, "device-id", uint32(deviceId))
		}
	}

	matches = regexp.MustCompile("^/uhppote/device/[0-9]+/door/([1-4])(?:$|/.*$$)").FindStringSubmatch(url)
	if matches != nil {
		door, err := strconv.ParseUint(matches[1], 10, 8)
		if err == nil {
			ctx = context.WithValue(ctx, "door", uint8(door))
		}
	}

	return ctx
}

func reply(ctx context.Context, w http.ResponseWriter, response interface{}) {
	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
