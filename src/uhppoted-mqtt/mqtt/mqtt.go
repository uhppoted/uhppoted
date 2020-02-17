package mqtt

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"regexp"
	"time"
	"uhppote"
	"uhppoted"
	"uhppoted-mqtt/auth"
)

type MQTTD struct {
	ServerID       string
	Connection     Connection
	TLS            *tls.Config
	Topics         Topics
	Alerts         Alerts
	HMAC           auth.HMAC
	Encryption     Encryption
	Authentication string
	Permissions    auth.Permissions
	EventMap       string
	Debug          bool

	client    paho.Client
	interrupt chan os.Signal
}

type Connection struct {
	Broker   string
	ClientID string
	UserName string
	Password string
}

type Topics struct {
	Requests string
	Replies  string
	Events   string
	System   string
}

type Alerts struct {
	QOS      byte
	Retained bool
}

type Encryption struct {
	SignOutgoing    bool
	EncryptOutgoing bool
	EventsKeyID     string
	SystemKeyID     string
	HOTP            *auth.HOTP
	RSA             *auth.RSA
	Nonce           auth.Nonce
}

type fdispatch struct {
	method string
	f      func(*MQTTD, metainfo, *uhppoted.UHPPOTED, context.Context, []byte) (interface{}, error)
}

type dispatcher struct {
	mqttd    *MQTTD
	uhppoted *uhppoted.UHPPOTED
	uhppote  *uhppote.UHPPOTE
	log      *log.Logger
	table    map[string]fdispatch
}

type request struct {
	ClientID  *string
	RequestID *string
	ReplyTo   *string
	Request   []byte
}

type metainfo struct {
	RequestID *string `json:"request-id,omitempty"`
	ClientID  *string `json:"client-id,omitempty"`
	ServerID  string  `json:"server-id,omitempty"`
	Method    string  `json:"method,omitempty"`
	Nonce     fnonce  `json:"nonce,omitempty"`
}

type errorx struct {
	Err     error  `json:"-"`
	Code    int    `json:"error-code"`
	Message string `json:"message"`
}

func (e *errorx) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

type fnonce func() uint64

func (f fnonce) MarshalJSON() ([]byte, error) {
	return json.Marshal(f())
}

var regex = struct {
	clean  *regexp.Regexp
	base64 *regexp.Regexp
}{
	clean:  regexp.MustCompile(`\s+`),
	base64: regexp.MustCompile(`^"[A-Za-z0-9+/]*[=]{0,2}"$`),
}

func (mqttd *MQTTD) Run(u *uhppote.UHPPOTE, log *log.Logger) {
	paho.CRITICAL = log
	paho.ERROR = log
	paho.WARN = log

	if mqttd.Debug {
		paho.DEBUG = log
	}

	api := uhppoted.UHPPOTED{
		Uhppote: u,
		Log:     log,
	}

	d := dispatcher{
		mqttd:    mqttd,
		uhppoted: &api,
		uhppote:  u,
		log:      log,
		table: map[string]fdispatch{
			mqttd.Topics.Requests + "/devices:get":             fdispatch{"get-devices", (*MQTTD).getDevices},
			mqttd.Topics.Requests + "/device:get":              fdispatch{"get-device", (*MQTTD).getDevice},
			mqttd.Topics.Requests + "/device/status:get":       fdispatch{"get-status", (*MQTTD).getStatus},
			mqttd.Topics.Requests + "/device/time:get":         fdispatch{"get-time", (*MQTTD).getTime},
			mqttd.Topics.Requests + "/device/time:set":         fdispatch{"set-time", (*MQTTD).setTime},
			mqttd.Topics.Requests + "/device/door/delay:get":   fdispatch{"get-door-delay", (*MQTTD).getDoorDelay},
			mqttd.Topics.Requests + "/device/door/delay:set":   fdispatch{"set-door-delay", (*MQTTD).setDoorDelay},
			mqttd.Topics.Requests + "/device/door/control:get": fdispatch{"get-door-control", (*MQTTD).getDoorControl},
			mqttd.Topics.Requests + "/device/door/control:set": fdispatch{"set-door-control", (*MQTTD).setDoorControl},
			mqttd.Topics.Requests + "/device/cards:get":        fdispatch{"get-cards", (*MQTTD).getCards},
			mqttd.Topics.Requests + "/device/cards:delete":     fdispatch{"delete-cards", (*MQTTD).deleteCards},
			mqttd.Topics.Requests + "/device/card:get":         fdispatch{"get-card", (*MQTTD).getCard},
			mqttd.Topics.Requests + "/device/card:put":         fdispatch{"put-card", (*MQTTD).putCard},
			mqttd.Topics.Requests + "/device/card:delete":      fdispatch{"delete-card", (*MQTTD).deleteCard},
			mqttd.Topics.Requests + "/device/events:get":       fdispatch{"get-events", (*MQTTD).getEvents},
			mqttd.Topics.Requests + "/device/event:get":        fdispatch{"get-event", (*MQTTD).getEvent},
		},
	}

	client, err := mqttd.subscribeAndServe(&d, log)
	if err != nil {
		log.Printf("ERROR: Error connecting to '%s': %v", mqttd.Connection.Broker, err)
		mqttd.Close(log)
		return
	}

	mqttd.client = client

	if err := mqttd.listen(&api, u, log); err != nil {
		log.Printf("ERROR: Error binding to listen port '%d': %v", 12345, err)
		mqttd.Close(log)
		return
	}
}

func (m *MQTTD) Close(log *log.Logger) {
	if m.interrupt != nil {
		close(m.interrupt)
	}

	if m.client != nil {
		log.Printf("INFO  closing connection to %s", m.Connection.Broker)
		m.client.Disconnect(250)
		log.Printf("INFO  closed connection to %s", m.Connection.Broker)
	}

	m.client = nil
}

func (m *MQTTD) subscribeAndServe(d *dispatcher, log *log.Logger) (paho.Client, error) {
	var handler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
		d.dispatch(client, msg)
	}

	var connected paho.OnConnectHandler = func(client paho.Client) {
		options := client.OptionsReader()
		servers := options.Servers()
		for _, url := range servers {
			log.Printf("%-5s %-12s %v", "INFO", "mqttd", fmt.Sprintf("Connected to %s", url))
		}

		token := m.client.Subscribe(m.Topics.Requests+"/#", 0, handler)
		if err := token.Error(); err != nil {
			log.Printf("ERROR unable to subscribe to %s (%v)", m.Topics.Requests, err)
			return
		}

		log.Printf("%-5s %-12s %v", "INFO", "mqttd", fmt.Sprintf("Subscribed to %s", m.Topics.Requests))
	}

	var disconnected paho.ConnectionLostHandler = func(client paho.Client, err error) {
		log.Printf("ERROR connection to MQTT broker lost (%v)", err)
	}

	options := paho.
		NewClientOptions().
		AddBroker(m.Connection.Broker).
		SetClientID(m.Connection.ClientID).
		SetTLSConfig(m.TLS).
		SetCleanSession(false).
		SetConnectRetry(true).
		SetConnectRetryInterval(30 * time.Second).
		SetOnConnectHandler(connected).
		SetConnectionLostHandler(disconnected)

	if m.Connection.UserName != "" {
		options.SetUsername(m.Connection.UserName)
		if m.Connection.Password != "" {
			options.SetPassword(m.Connection.Password)
		}
	}

	client := paho.NewClient(options)
	token := client.Connect()
	if err := token.Error(); err != nil {
		return nil, err
	}

	return client, nil
}

func (m *MQTTD) listen(api *uhppoted.UHPPOTED, u *uhppote.UHPPOTE, log *log.Logger) error {
	log.Printf("%-5s %-12s %v", "INFO", "mqttd", fmt.Sprintf("Listening on %v", u.ListenAddress))
	log.Printf("%-5s %-12s %v", "INFO", "mqttd", fmt.Sprintf("Publishing events to %s", m.Topics.Events))

	last := uhppoted.NewEventMap(m.EventMap)
	if err := last.Load(log); err != nil {
		log.Printf("WARN  Error loading event map [%v]", err)
	}

	handler := func(event uhppoted.EventMessage) {
		if err := m.send(&m.Encryption.EventsKeyID, m.Topics.Events, event, msgEvent, false); err != nil {
			log.Printf("WARN  %-20s %v", "listen", err)
		}
	}

	m.interrupt = make(chan os.Signal)

	go func() {
		api.Listen(handler, last, m.interrupt)
	}()

	return nil
}

func (d *dispatcher) dispatch(client paho.Client, msg paho.Message) {
	ctx := context.WithValue(context.Background(), "client", client)
	ctx = context.WithValue(ctx, "log", d.log)

	if fn, ok := d.table[msg.Topic()]; ok {
		msg.Ack()

		go func() {
			rq, err := d.mqttd.unwrap(msg.Payload())
			if err != nil {
				d.log.Printf("DEBUG %-20s %s", "dispatch", string(msg.Payload()))
				d.log.Printf("WARN  %-20s %v", "dispatch", err)
				return
			}

			if err := d.mqttd.authorise(rq.ClientID, msg.Topic()); err != nil {
				d.log.Printf("DEBUG %-20s %s", fn.method, string(rq.Request))
				d.log.Printf("WARN  %-20s %v", fn.method, fmt.Errorf("Error authorising request (%v)", err))
				return
			}

			ctx = context.WithValue(ctx, "request", rq.Request)
			ctx = context.WithValue(ctx, "method", fn.method)

			replyTo := d.mqttd.Topics.Replies

			if rq.ClientID != nil {
				replyTo = d.mqttd.Topics.Replies + "/" + *rq.ClientID
			}

			if rq.ReplyTo != nil {
				replyTo = *rq.ReplyTo
			}

			meta := metainfo{
				RequestID: rq.RequestID,
				ClientID:  rq.ClientID,
				ServerID:  d.mqttd.ServerID,
				Method:    fn.method,
				Nonce:     func() uint64 { return d.mqttd.Encryption.Nonce.Next() },
			}

			reply, err := fn.f(d.mqttd, meta, d.uhppoted, ctx, rq.Request)

			if err != nil {
				d.log.Printf("DEBUG %-20s %s", fn.method, string(rq.Request))
				d.log.Printf("WARN  %-20s %v", fn.method, err)

				if errx, ok := err.(*errorx); ok {
					if err := d.mqttd.send(rq.ClientID, replyTo, errx, msgError, false); err != nil {
						d.log.Printf("WARN  %-20s %v", fn.method, err)
					}
				}
			}

			if reply != nil {
				if err := d.mqttd.send(rq.ClientID, replyTo, reply, msgReply, false); err != nil {
					d.log.Printf("WARN  %-20s %v", fn.method, err)
				}
			}
		}()
	}
}

func (m *MQTTD) authorise(clientID *string, topic string) error {
	if m.Permissions.Enabled {
		if clientID == nil {
			return errors.New("Request without client-id")
		}

		match := regexp.MustCompile(`.*?/(\w+):(\w+)$`).FindStringSubmatch(topic)
		if len(match) != 3 {
			return fmt.Errorf("Invalid resource:action (%s)", topic)
		}

		return m.Permissions.Validate(*clientID, match[1], match[2])
	}

	return nil
}

// TODO: add callback for published/failed
func (mqttd *MQTTD) send(destID *string, topic string, message interface{}, msgtype msgType, critical bool) error {
	if mqttd.client == nil {
		return errors.New("No connection to MQTT broker")
	}

	m, err := mqttd.wrap(msgtype, message, destID)
	if err != nil {
		return err
	} else if m == nil {
		return errors.New("'wrap' failed to return a publishable message")
	}

	qos := byte(0)
	retained := false
	if critical {
		qos = mqttd.Alerts.QOS
		retained = mqttd.Alerts.Retained
	}

	token := mqttd.client.Publish(topic, qos, retained, string(m))
	if token.Error() != nil {
		return token.Error()
	}

	return nil
}

func isBase64(request []byte) bool {
	return regex.base64.Match(request)
}

func clean(s string) string {
	return regex.clean.ReplaceAllString(s, " ")
}
