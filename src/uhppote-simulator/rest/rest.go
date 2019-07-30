package rest

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"uhppote-simulator/simulator"
	"uhppote/types"
)

type AccessCard struct {
	CardNumber uint32 `json:"card-number"`
}

type Swipe struct {
	Card AccessCard `json:"card"`
}

type context struct {
	simulators []*simulator.Simulator
}

type handler struct {
	re *regexp.Regexp
	fn func(*context, http.ResponseWriter, *http.Request)
}

type dispatcher struct {
	ctx      *context
	handlers []handler
}

func Run(simulators []*simulator.Simulator) {
	d := new(dispatcher)

	d.ctx = new(context)
	d.ctx.simulators = simulators

	d.Add("^/uhppote/simulator/[0-9]+/door/[1-4]/swipe$", swipe)

	log.Fatal(http.ListenAndServe(":8080", d))
}

func (d *dispatcher) Add(path string, h func(*context, http.ResponseWriter, *http.Request)) {
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
	matches := regexp.MustCompile("^/uhppote/simulator/([0-9]+)/door/([1-4])/swipe$").FindStringSubmatch(url)
	deviceId, _ := strconv.ParseUint(matches[1], 10, 32)
	door, _ := strconv.ParseUint(matches[2], 10, 8)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var swipe Swipe
	err = json.Unmarshal(blob, &swipe)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
	} else {
		fmt.Printf("SWIPE: %v\n", swipe)
		fmt.Println(deviceId)
		fmt.Println(door)
		fmt.Println(string(blob))

		for _, s := range ctx.simulators {
			if s.SerialNumber == types.SerialNumber(deviceId) {
				for _, c := range s.Cards {
					if c.CardNumber == swipe.Card.CardNumber {
						switch door {
						case 1:
							if c.Door1 {
								fmt.Fprintf(w, "Greeting earthling, %q", html.EscapeString(r.URL.Path))
								return
							}
						case 2:
							if c.Door2 {
								fmt.Fprintf(w, "Greeting earthling, %q", html.EscapeString(r.URL.Path))
								return
							}
						case 3:
							if c.Door3 {
								fmt.Fprintf(w, "Greeting earthling, %q", html.EscapeString(r.URL.Path))
								return
							}
						case 4:
							if c.Door4 {
								fmt.Fprintf(w, "Greeting earthling, %q", html.EscapeString(r.URL.Path))
								return
							}

						}
					}
				}
			}
		}

		fmt.Fprintf(w, "Access denied")
	}
}
