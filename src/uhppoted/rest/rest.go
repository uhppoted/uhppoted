package rest

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"uhppote"
)

type RestD struct {
	HttpEnabled        bool
	HttpPort           uint16
	HttpsEnabled       bool
	HttpsPort          uint16
	TLSKeyFile         string
	TLSCertificateFile string
	CACertificateFile  string
}

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

func (r *RestD) Run(u *uhppote.UHPPOTE, l *log.Logger) {
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
	d.Add("^/uhppote/device/[0-9]+/card/[0-9]+$", http.MethodGet, getCard)
	d.Add("^/uhppote/device/[0-9]+/card$", http.MethodDelete, deleteCards)
	d.Add("^/uhppote/device/[0-9]+/card/[0-9]+$", http.MethodDelete, deleteCard)
	d.Add("^/uhppote/device/[0-9]+/event$", http.MethodGet, getEvents)
	d.Add("^/uhppote/device/[0-9]+/event/[0-9]+$", http.MethodGet, getEvent)

	var wg sync.WaitGroup

	if r.HttpEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("... listening on port %d\n", r.HttpPort)
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", r.HttpPort), &d))
		}()
	}

	if r.HttpsEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("... listening on port %d\n", r.HttpsPort)

			ca, err := ioutil.ReadFile(r.CACertificateFile)
			if err != nil {
				log.Fatal(err)
			}

			certificates := x509.NewCertPool()
			if !certificates.AppendCertsFromPEM(ca) {
				log.Fatal("Unable failed to parse CA certificate")
			}

			tlsConfig := tls.Config{
				ClientAuth: tls.RequireAndVerifyClientCert,
				ClientCAs:  certificates,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				},
				PreferServerCipherSuites: true,
				MinVersion:               tls.VersionTLS12,
			}

			tlsConfig.BuildNameToCertificate()

			httpsd := &http.Server{
				Addr:      fmt.Sprintf(":%d", r.HttpsPort),
				Handler:   &d,
				TLSConfig: &tlsConfig,
			}

			log.Fatal(httpsd.ListenAndServeTLS(r.TLSCertificateFile, r.TLSKeyFile))
		}()
	}

	wg.Wait()
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

	matches = regexp.MustCompile("^/uhppote/device/[0-9]+/door/([1-4])(?:$|/.*$)").FindStringSubmatch(url)
	if matches != nil {
		door, err := strconv.ParseUint(matches[1], 10, 8)
		if err == nil {
			ctx = context.WithValue(ctx, "door", uint8(door))
		}
	}

	matches = regexp.MustCompile("^/uhppote/device/[0-9]+/card/([0-9]+)$").FindStringSubmatch(url)
	if matches != nil {
		cardNumber, err := strconv.ParseUint(matches[1], 10, 32)
		if err == nil {
			ctx = context.WithValue(ctx, "card-number", uint32(cardNumber))
		}
	}

	matches = regexp.MustCompile("^/uhppote/device/[0-9]+/event/([0-9]+)$").FindStringSubmatch(url)
	if matches != nil {
		eventId, err := strconv.ParseUint(matches[1], 10, 32)
		if err == nil {
			ctx = context.WithValue(ctx, "event-id", uint32(eventId))
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

func debug(ctx context.Context, serialNumber uint32, operation string, r *http.Request) {
	ctx.Value("log").(*log.Logger).Printf("DEBUG %-12d %-20s %v\n", serialNumber, operation, *r)
}

func warn(ctx context.Context, serialNumber uint32, operation string, err error) {
	ctx.Value("log").(*log.Logger).Printf("WARN  %-12d %-20s %v\n", serialNumber, operation, err)
}
