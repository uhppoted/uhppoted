package rest

import (
	"compress/gzip"
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
	"strings"
	"sync"
	"uhppote"
)

type OpenApi struct {
	Enabled   bool
	Directory string
}

type RestD struct {
	HttpEnabled        bool
	HttpPort           uint16
	HttpsEnabled       bool
	HttpsPort          uint16
	TLSKeyFile         string
	TLSCertificateFile string
	CACertificateFile  string
	CORSEnabled        bool
	OpenApi
}

type handlerfn func(context.Context, http.ResponseWriter, *http.Request)

type handler struct {
	re     *regexp.Regexp
	method string
	fn     handlerfn
}

type dispatcher struct {
	corsEnabled bool
	uhppote     *uhppote.UHPPOTE
	log         *log.Logger
	handlers    []handler
	openapi     http.Handler
}

func (r *RestD) Run(u *uhppote.UHPPOTE, l *log.Logger) {
	d := dispatcher{
		uhppote: u,
		handlers: []handler{
			handler{regexp.MustCompile("^/uhppote/device$"), http.MethodGet, getDevices},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+$"), http.MethodGet, getDevice},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/status$"), http.MethodGet, getStatus},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/time$"), http.MethodGet, getTime},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/time$"), http.MethodPut, setTime},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/door/[1-4]/delay$"), http.MethodGet, getDoorDelay},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/door/[1-4]/delay$"), http.MethodPut, setDoorDelay},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/door/[1-4]/state$"), http.MethodGet, getDoorControlState},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/door/[1-4]/state$"), http.MethodPut, setDoorControlState},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/card$"), http.MethodGet, getCards},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/card/[0-9]+$"), http.MethodGet, getCard},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/card$"), http.MethodDelete, deleteCards},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/card/[0-9]+$"), http.MethodDelete, deleteCard},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/event$"), http.MethodGet, getEvents},
			handler{regexp.MustCompile("^/uhppote/device/[0-9]+/event/[0-9]+$"), http.MethodGet, getEvent},
		},
		log:         l,
		corsEnabled: r.CORSEnabled,
		openapi:     http.NotFoundHandler(),
	}

	if r.OpenApi.Enabled {
		d.openapi = http.StripPrefix("/openapi", http.FileServer(http.Dir(r.OpenApi.Directory)))
	}

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

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	compression := "none"

	for key, headers := range r.Header {
		if http.CanonicalHeaderKey(key) == "Accept-Encoding" {
			for _, header := range headers {
				encodings := strings.Split(header, ",")
				for _, encoding := range encodings {
					if strings.TrimSpace(encoding) == "gzip" {
						compression = "gzip"
					}
				}
			}
		}
	}

	if d.corsEnabled {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// CORS pre-flight request ?
		if r.Method == http.MethodOptions {
			return
		}
	}

	// OpenAPI ?

	if strings.HasPrefix(r.URL.Path, "/openapi") {
		d.openapi.ServeHTTP(w, r)
		return
	}

	// Dispatch to handler
	url := r.URL.Path
	for _, h := range d.handlers {
		if h.re.MatchString(url) && r.Method == h.method {
			ctx := context.WithValue(context.Background(), "uhppote", d.uhppote)
			ctx = context.WithValue(ctx, "log", d.log)
			ctx = context.WithValue(ctx, "compression", compression)
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

	if len(b) > 1024 && ctx.Value("compression") == "gzip" {
		w.Header().Set("Content-Encoding", "gzip")
		encoder := gzip.NewWriter(w)
		encoder.Write(b)
		encoder.Flush()
	} else {
		w.Write(b)
	}
}

func debug(ctx context.Context, serialNumber uint32, operation string, r *http.Request) {
	ctx.Value("log").(*log.Logger).Printf("DEBUG %-12d %-20s %v\n", serialNumber, operation, *r)
}

func warn(ctx context.Context, serialNumber uint32, operation string, err error) {
	ctx.Value("log").(*log.Logger).Printf("WARN  %-12d %-20s %v\n", serialNumber, operation, err)
}
