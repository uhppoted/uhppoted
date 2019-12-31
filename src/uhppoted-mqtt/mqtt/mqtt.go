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
	Broker      string
	TLS         *tls.Config
	Topic       string
	HOTP        auth.HOTP
	Permissions auth.Permissions
	EventMap    string
	Debug       bool

	connection MQTT.Client
	interrupt  chan os.Signal
}

type fdispatch func(*uhppoted.UHPPOTED, context.Context, uhppoted.Request)
type fdispatchx struct {
	operation string
	f         func(*MQTTD, *uhppoted.UHPPOTED, context.Context, MQTT.Message)
}

type dispatcher struct {
	mqttd    *MQTTD
	uhppoted *uhppoted.UHPPOTED
	uhppote  *uhppote.UHPPOTE
	log      *log.Logger
	topic    string
	table    map[string]fdispatch
	tablex   map[string]fdispatchx
}

type request struct {
	RequestID *string `json:"request-id"`
	ReplyTo   *string `json:"reply-to"`
	ClientID  *string `json:"client-id"`
	HOTP      *string `json:"hotp"`
}

type metainfo struct {
	RequestID string `json:"request-id,omitempty"`
	ClientID  string `json:"client-id,omitempty"`
	Operation string `json:"operation,omitempty"`
}

func (m *MQTTD) Run(u *uhppote.UHPPOTE, l *log.Logger) {
	MQTT.CRITICAL = l
	MQTT.ERROR = l
	MQTT.WARN = l

	if m.Debug {
		MQTT.DEBUG = l
	}

	api := uhppoted.UHPPOTED{
		Service: m,
	}

	d := dispatcher{
		mqttd:    m,
		uhppoted: &api,
		uhppote:  u,
		log:      l,
		topic:    m.Topic,
		table: map[string]fdispatch{
			m.Topic + "/device/door/control:set": (*uhppoted.UHPPOTED).SetDoorControl,
		},
		tablex: map[string]fdispatchx{
			m.Topic + "/devices:get":             fdispatchx{"get-devices", (*MQTTD).getDevices},
			m.Topic + "/device:get":              fdispatchx{"get-device", (*MQTTD).getDevice},
			m.Topic + "/device/status:get":       fdispatchx{"get-status", (*MQTTD).getStatus},
			m.Topic + "/device/time:get":         fdispatchx{"get-time", (*MQTTD).getTime},
			m.Topic + "/device/time:set":         fdispatchx{"set-time", (*MQTTD).setTime},
			m.Topic + "/device/door/delay:get":   fdispatchx{"get-door-delay", (*MQTTD).getDoorDelay},
			m.Topic + "/device/door/delay:set":   fdispatchx{"set-door-delay", (*MQTTD).setDoorDelay},
			m.Topic + "/device/door/control:get": fdispatchx{"get-door-control", (*MQTTD).getDoorControl},

			m.Topic + "/device/cards:get":    fdispatchx{"get-cards", (*MQTTD).getCards},
			m.Topic + "/device/cards:delete": fdispatchx{"delete-cards", (*MQTTD).deleteCards},
			m.Topic + "/device/card:get":     fdispatchx{"get-card", (*MQTTD).getCard},
			m.Topic + "/device/card:put":     fdispatchx{"put-card", (*MQTTD).putCard},
			m.Topic + "/device/card:delete":  fdispatchx{"delete-card", (*MQTTD).deleteCard},
			m.Topic + "/device/events:get":   fdispatchx{"get-events", (*MQTTD).getEvents},
			m.Topic + "/device/event:get":    fdispatchx{"get-event", (*MQTTD).getEvent},
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

	m.interrupt = make(chan os.Signal)

	go func() {
		api.Listen(ctx, last, m.interrupt)
	}()

	return nil
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

func (d *dispatcher) dispatch(client MQTT.Client, msg MQTT.Message) {
	ctx := context.WithValue(context.Background(), "uhppote", d.uhppote)
	ctx = context.WithValue(ctx, "client", client)
	ctx = context.WithValue(ctx, "log", d.log)
	ctx = context.WithValue(ctx, "topic", d.topic)

	// TODO (in progress) replace with local wrapper functions
	if fn, ok := d.table[msg.Topic()]; ok {
		rq := Request{
			Message: msg,
		}

		if err := json.Unmarshal(msg.Payload(), &rq); err != nil {
			d.log.Printf("WARN  %-20s %v\n", "dispatch", fmt.Errorf("Invalid message: %v <%s>", err, string(msg.Payload())))
			return
		}

		fn(d.uhppoted, ctx, &rq)
		return
	}
	// ----

	if fnx, ok := d.tablex[msg.Topic()]; ok {
		body := struct {
			Request request `json:"request"`
		}{}

		if err := json.Unmarshal(msg.Payload(), &body); err != nil {
			d.log.Printf("WARN  %-20s %v %s\n", "dispatch", err, string(msg.Payload()))
			return
		}

		if err := d.mqttd.authenticate(body.Request); err != nil {
			d.log.Printf("WARN  %-20s %v %s\n", "dispatch", err, string(msg.Payload()))
			return
		}

		if err := d.mqttd.authorise(body.Request, msg.Topic()); err != nil {
			d.log.Printf("WARN  %-20s %v %s\n", "dispatch", err, string(msg.Payload()))
			return
		}

		ctx = context.WithValue(ctx, "request", body.Request)
		ctx = context.WithValue(ctx, "operation", fnx.operation)

		fnx.f(d.mqttd, d.uhppoted, ctx, msg)
	}
}

func (m *MQTTD) authenticate(rq request) error {
	if m.HOTP.Enabled {
		if rq.ClientID == nil {
			return errors.New("Request without client-id")
		} else if rq.HOTP == nil {
			return errors.New("Request without HOTP")
		}

		return m.HOTP.Validate(*rq.ClientID, *rq.HOTP)
	}

	return nil
}

func (m *MQTTD) authorise(rq request, topic string) error {
	if m.Permissions.Enabled {
		if rq.ClientID == nil {
			return errors.New("Request without client-id")
		}

		match := regexp.MustCompile(`.*?/(\w+):(\w+)$`).FindStringSubmatch(topic)
		if len(match) != 3 {
			return fmt.Errorf("Invalid resource:action (%s)", topic)
		}

		return m.Permissions.Validate(*rq.ClientID, match[1], match[2])
	}

	return nil
}

func (m *MQTTD) Send(ctx context.Context, message interface{}) {
	b, err := json.Marshal(message)
	if err != nil {
		oops(ctx, "encoding/json", "Error encoding message", uhppoted.StatusInternalServerError)
		return
	}

	client, ok := ctx.Value("client").(MQTT.Client)
	if !ok {
		panic("MQTT client not included in context")
	}

	topic, ok := ctx.Value("topic").(string)
	if !ok {
		panic("MQTT root topic not included in context")
	}

	token := client.Publish(topic+"/events", 0, false, string(b))
	token.Wait()
}

func (m *MQTTD) Reply(ctx context.Context, response interface{}) {
	client, ok := ctx.Value("client").(MQTT.Client)
	if !ok {
		panic("MQTT client not included in context")
	}

	topic, ok := ctx.Value("topic").(string)
	if !ok {
		panic("MQTT root topic not included in context")
	}

	replyTo := "reply"
	if rq, ok := ctx.Value("request").(request); ok {
		if rq.ReplyTo != nil {
			replyTo = *rq.ReplyTo
		}
	}

	b, err := json.Marshal(response)
	if err != nil {
		oops(ctx, "encoding/json", "Error encoding response", uhppoted.StatusInternalServerError)
		return
	}

	token := client.Publish(topic+"/"+replyTo, 0, false, string(b))
	token.Wait()
}

func getMetaInfo(ctx context.Context) *metainfo {
	metainfo := metainfo{
		RequestID: "",
		ClientID:  "",
		Operation: "",
	}

	if operation, ok := ctx.Value("operation").(string); ok {
		metainfo.Operation = operation
	}

	if rq, ok := ctx.Value("request").(request); ok {
		if rq.RequestID != nil {
			metainfo.RequestID = *rq.RequestID
		}

		if rq.ClientID != nil {
			metainfo.ClientID = *rq.ClientID
		}
	}

	return &metainfo
}

func (m *MQTTD) Oops(ctx context.Context, operation string, message string, errorCode int) {
	oops(ctx, operation, message, errorCode)
}

func (m *MQTTD) OnError(ctx context.Context, message string, errorCode int, err error) {

	if operation, ok := ctx.Value("operation").(string); ok {
		ctx.Value("log").(*log.Logger).Printf("WARN  %-20s %v\n", operation, err)
		oops(ctx, operation, message, errorCode)
		return
	}

	ctx.Value("log").(*log.Logger).Printf("WARN  %v\n", err)
	oops(ctx, "???", message, errorCode)
}

func oops(ctx context.Context, operation string, msg string, errorCode int) {
	client, ok := ctx.Value("client").(MQTT.Client)
	if !ok {
		panic("MQTT client not included in context")
	}

	topic, ok := ctx.Value("topic").(string)
	if !ok {
		panic("MQTT root topic not included in context")
	}

	requestID := ""
	replyTo := "errors"

	rq, ok := ctx.Value("request").(request)
	if ok {
		if rq.RequestID != nil {
			requestID = *rq.RequestID
		}

		if rq.ReplyTo != nil {
			replyTo = *rq.ReplyTo + "/errors"
		}
	}

	response := struct {
		Meta struct {
			RequestID string `json:"request-id,omitempty"`
		} `json:"meta-info"`
		Operation string `json:"operation"`
		Error     struct {
			Message   string `json:"message"`
			ErrorCode int    `json:"error-code"`
		} `json:"error"`
	}{
		Meta: struct {
			RequestID string `json:"request-id,omitempty"`
		}{
			RequestID: requestID,
		},
		Operation: operation,
		Error: struct {
			Message   string `json:"message"`
			ErrorCode int    `json:"error-code"`
		}{
			Message:   msg,
			ErrorCode: errorCode,
		},
	}

	b, err := json.Marshal(response)
	if err != nil {
		ctx.Value("log").(*log.Logger).Printf("ERROR: Error generating JSON response (%v)", err)
		return
	}

	token := client.Publish(topic+"/"+replyTo, 0, false, string(b))
	token.Wait()
}

func debug(ctx context.Context, operation string, msg interface{}) {
	ctx.Value("log").(*log.Logger).Printf("DEBUG %-20s %v\n", operation, msg)
}
