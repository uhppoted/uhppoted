package mqtt

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"regexp"
	"uhppote"
	"uhppoted"
	"uhppoted-mqtt/auth"
)

type MQTTD struct {
	ServerID        string
	Broker          string
	TLS             *tls.Config
	Topic           string
	EventsTopic     string
	EventsKeyID     string
	HMAC            auth.HMAC
	Authentication  string
	HOTP            *auth.HOTP
	RSA             *auth.RSA
	Nonce           auth.Nonce
	Permissions     auth.Permissions
	EventMap        string
	SignOutgoing    bool
	EncryptOutgoing bool
	Debug           bool

	connection MQTT.Client
	interrupt  chan os.Signal
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
	topic    string
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

func (m *MQTTD) Run(u *uhppote.UHPPOTE, l *log.Logger) {
	MQTT.CRITICAL = l
	MQTT.ERROR = l
	MQTT.WARN = l

	if m.Debug {
		MQTT.DEBUG = l
	}

	api := uhppoted.UHPPOTED{
		Log: l,
	}

	d := dispatcher{
		mqttd:    m,
		uhppoted: &api,
		uhppote:  u,
		log:      l,
		topic:    m.Topic,
		table: map[string]fdispatch{
			m.Topic + "/devices:get":             fdispatch{"get-devices", (*MQTTD).getDevices},
			m.Topic + "/device:get":              fdispatch{"get-device", (*MQTTD).getDevice},
			m.Topic + "/device/status:get":       fdispatch{"get-status", (*MQTTD).getStatus},
			m.Topic + "/device/time:get":         fdispatch{"get-time", (*MQTTD).getTime},
			m.Topic + "/device/time:set":         fdispatch{"set-time", (*MQTTD).setTime},
			m.Topic + "/device/door/delay:get":   fdispatch{"get-door-delay", (*MQTTD).getDoorDelay},
			m.Topic + "/device/door/delay:set":   fdispatch{"set-door-delay", (*MQTTD).setDoorDelay},
			m.Topic + "/device/door/control:get": fdispatch{"get-door-control", (*MQTTD).getDoorControl},
			m.Topic + "/device/door/control:set": fdispatch{"set-door-control", (*MQTTD).setDoorControl},
			m.Topic + "/device/cards:get":        fdispatch{"get-cards", (*MQTTD).getCards},
			m.Topic + "/device/cards:delete":     fdispatch{"delete-cards", (*MQTTD).deleteCards},
			m.Topic + "/device/card:get":         fdispatch{"get-card", (*MQTTD).getCard},
			m.Topic + "/device/card:put":         fdispatch{"put-card", (*MQTTD).putCard},
			m.Topic + "/device/card:delete":      fdispatch{"delete-card", (*MQTTD).deleteCard},
			m.Topic + "/device/events:get":       fdispatch{"get-events", (*MQTTD).getEvents},
			m.Topic + "/device/event:get":        fdispatch{"get-event", (*MQTTD).getEvent},
		},
	}

	if err := m.subscribeAndServe(&d); err != nil {
		l.Printf("ERROR: Error connecting to '%s': %v", m.Broker, err)
		m.Close(l)
		return
	}

	log.Printf("... connected to %s\n", m.Broker)

	if err := m.listen(&api, u, l); err != nil {
		l.Printf("ERROR: Error binding to listen port '%d': %v", 12345, err)
		m.Close(l)
		return
	}
}

func (m *MQTTD) Close(l *log.Logger) {
	if m.interrupt != nil {
		close(m.interrupt)
	}

	if m.connection != nil {
		log.Printf("... closing connection to %s", m.Broker)
		token := m.connection.Unsubscribe(m.Topic + "/#")
		if token.Wait() && token.Error() != nil {
			l.Printf("WARN: Error unsubscribing from topic '%s': %v", m.Topic, token.Error())
		}

		m.connection.Disconnect(250)
	}

	m.connection = nil
}

func (m *MQTTD) subscribeAndServe(d *dispatcher) error {
	var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		d.dispatch(client, msg)
	}

	options := MQTT.NewClientOptions()

	options.AddBroker(m.Broker)
	options.SetClientID("twystd-uhppoted-mqttd")
	options.SetDefaultPublishHandler(f)
	options.SetTLSConfig(m.TLS)

	m.connection = MQTT.NewClient(options)
	token := m.connection.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	token = m.connection.Subscribe(m.Topic+"/#", 0, nil)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (m *MQTTD) listen(api *uhppoted.UHPPOTED, u *uhppote.UHPPOTE, l *log.Logger) error {
	log.Printf("... listening on %v", u.ListenAddress)

	ctx := context.WithValue(context.Background(), "uhppote", u)
	ctx = context.WithValue(ctx, "client", m.connection)
	ctx = context.WithValue(ctx, "log", l)
	ctx = context.WithValue(ctx, "topic", m.Topic)

	last := uhppoted.NewEventMap(m.EventMap)
	if err := last.Load(l); err != nil {
		l.Printf("WARN: Error loading event map [%v]", err)
	}

	handler := func(event uhppoted.EventMessage) {
		m.send("listen", &m.EventsKeyID, m.EventsTopic, event, msgEvent, l)
	}

	m.interrupt = make(chan os.Signal)

	go func() {
		api.Listen(ctx, handler, last, m.interrupt)
	}()

	return nil
}

func (d *dispatcher) dispatch(client MQTT.Client, msg MQTT.Message) {
	ctx := context.WithValue(context.Background(), "uhppote", d.uhppote)
	ctx = context.WithValue(ctx, "client", client)
	ctx = context.WithValue(ctx, "log", d.log)
	ctx = context.WithValue(ctx, "topic", d.topic)

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

			replyTo := d.mqttd.Topic + "/reply"

			if rq.ClientID != nil {
				replyTo = d.mqttd.Topic + "/" + *rq.ClientID
			}

			if rq.ReplyTo != nil {
				replyTo = *rq.ReplyTo
			}

			meta := metainfo{
				RequestID: rq.RequestID,
				ClientID:  rq.ClientID,
				ServerID:  d.mqttd.ServerID,
				Method:    fn.method,
				Nonce:     func() uint64 { return d.mqttd.Nonce.Next() },
			}

			reply, err := fn.f(d.mqttd, meta, d.uhppoted, ctx, rq.Request)

			if err != nil {
				d.log.Printf("DEBUG %-20s %s", fn.method, string(rq.Request))
				d.log.Printf("WARN  %-20s %v", fn.method, err)

				if errx, ok := err.(*errorx); ok {
					d.mqttd.send(fn.method, rq.ClientID, replyTo, errx, msgError, d.log)
				}
			}

			if reply != nil {
				d.mqttd.send(fn.method, rq.ClientID, replyTo, reply, msgReply, d.log)
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

func (mqttd *MQTTD) send(method string, destID *string, topic string, message interface{}, msgtype msgType, log *log.Logger) {
	if m, err := mqttd.wrap(msgtype, message, destID); err != nil {
		log.Printf("WARN  %-20s %v", method, err)
	} else if m != nil {
		mqttd.connection.Publish(topic, 0, false, string(m)).Wait()
	}
}

func isBase64(request []byte) bool {
	return regex.base64.Match(request)
}

func clean(s string) string {
	return regex.clean.ReplaceAllString(s, " ")
}
