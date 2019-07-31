package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"uhppote-simulator/simulator"
	"uhppote/types"
)

type SwipeRequest struct {
	Door       uint8  `json:"door"`
	CardNumber uint32 `json:"card-number"`
}

type SwipeResponse struct {
	Granted bool   `json:"access-granted"`
	Message string `json:"message"`
}

type context struct {
	simulators []*simulator.Simulator
}

type handlerfn func(*context, http.ResponseWriter, *http.Request)

type handler struct {
	re *regexp.Regexp
	fn handlerfn
}

type dispatcher struct {
	ctx      *context
	handlers []handler
}

func Run(simulators []*simulator.Simulator) {
	d := new(dispatcher)

	d.ctx = new(context)
	d.ctx.simulators = simulators

	d.Add("^/uhppote/simulator/[0-9]+/swipe$", swipe)

	log.Fatal(http.ListenAndServe(":8080", d))
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

func swipe(ctx *context, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("Invalid method:%s - expected POST", r.Method), http.StatusBadRequest)
		return
	}

	url := r.URL.Path
	matches := regexp.MustCompile("^/uhppote/simulator/([0-9]+)/swipe$").FindStringSubmatch(url)
	deviceId, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	request := SwipeRequest{}
	err = json.Unmarshal(blob, &request)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	granted := false
	message := "Access denied"

	for _, s := range ctx.simulators {
		if s.SerialNumber == types.SerialNumber(deviceId) {
			for _, c := range s.Cards {
				if c.CardNumber == request.CardNumber {
					granted = c.Doors[request.Door]
					message = "Access granted"
				}
			}
		}
	}

	response := SwipeResponse{
		Granted: granted,
		Message: message,
	}

	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
