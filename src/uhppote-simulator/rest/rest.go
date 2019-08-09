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
)

type Device struct {
	DeviceId   uint32 `json:"device-id"`
	DeviceType string `json:"device-type"`
	URI        string `json:"uri"`
}

type DeviceList struct {
	Devices []Device `json:"devices"`
}

type NewDeviceRequest struct {
	DeviceId   uint32 `json:"device-id"`
	DeviceType string `json:"device-type"`
	Compressed bool   `json:"compressed"`
}

type SwipeRequest struct {
	Door       uint8  `json:"door"`
	CardNumber uint32 `json:"card-number"`
}

type SwipeResponse struct {
	Granted bool   `json:"access-granted"`
	Opened  bool   `json:"door-opened"`
	Message string `json:"message"`
}

type handlerfn func(*simulator.Context, http.ResponseWriter, *http.Request)

type handler struct {
	re *regexp.Regexp
	fn handlerfn
}

type dispatcher struct {
	ctx      *simulator.Context
	handlers []handler
}

func Run(ctx *simulator.Context) {
	d := dispatcher{
		ctx,
		make([]handler, 0),
	}

	d.Add("^/uhppote/simulator$", devices)
	d.Add("^/uhppote/simulator/[0-9]+$", device)
	d.Add("^/uhppote/simulator/[0-9]+/swipe$", swipe)

	log.Fatal(http.ListenAndServe(":8080", &d))
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

func devices(ctx *simulator.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list(ctx, w, r)

	case http.MethodPost:
		create(ctx, w, r)

	default:
		http.Error(w, fmt.Sprintf("Invalid method:%s - expected GET or POST", r.Method), http.StatusMethodNotAllowed)
	}
}

func device(ctx *simulator.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		delete(ctx, w, r)

	default:
		http.Error(w, fmt.Sprintf("Invalid method:%s - expected DELETE", r.Method), http.StatusMethodNotAllowed)
	}
}

func list(ctx *simulator.Context, w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	devices := make([]Device, 0)

	ctx.DeviceList.Apply(func(s simulator.Simulator) {
		devices = append(devices, Device{
			DeviceId:   s.DeviceID(),
			DeviceType: "UTC3011-L04",
			URI:        fmt.Sprintf("/uhppote/simulator/%d", s.DeviceID()),
		})
	})

	response := DeviceList{devices}
	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func create(ctx *simulator.Context, w http.ResponseWriter, r *http.Request) {
	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	request := NewDeviceRequest{}
	err = json.Unmarshal(blob, &request)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if request.DeviceId < 1 {
		http.Error(w, "Missing device ID", http.StatusBadRequest)
		return
	}

	if request.DeviceType != "UTC0311-L04" {
		http.Error(w, "Invalid  device type - expected UTC0311-L04", http.StatusBadRequest)
		return
	}

	err = ctx.DeviceList.Add(request.DeviceId, request.Compressed, ctx.Directory)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating device %d: %v", request.DeviceId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/uhppote/simulator/%d", request.DeviceId))
	w.Header().Set("Content-Type", "application/json")
}

func delete(ctx *simulator.Context, w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	matches := regexp.MustCompile("^/uhppote/simulator/([0-9]+)$").FindStringSubmatch(url)
	deviceId, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	_, err = ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	if s := ctx.DeviceList.Find(uint32(deviceId)); s == nil {
		http.Error(w, fmt.Sprintf("No device with ID %d", deviceId), http.StatusNotFound)
		return
	}

	ctx.DeviceList.Delete(uint32(deviceId))

	w.Header().Set("Content-Type", "application/json")
}

func swipe(ctx *simulator.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("Invalid method:%s - expected POST", r.Method), http.StatusMethodNotAllowed)
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

	s := ctx.DeviceList.Find(uint32(deviceId))
	if s == nil {
		http.Error(w, fmt.Sprintf("No device with ID %d", deviceId), http.StatusNotFound)
		return
	}

	granted, eventId := s.Swipe(uint32(deviceId), request.CardNumber, request.Door)
	opened := false
	message := "Access denied"

	if granted {
		opened = true
		message = "Access granted"
	}

	response := SwipeResponse{
		Granted: granted,
		Opened:  opened,
		Message: message,
	}

	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/uhppote/simulator/%d/events/%d", s.DeviceID(), eventId))
	w.Write(b)
}
