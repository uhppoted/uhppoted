package mqtt

import (
	"context"
	"encoding/json"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"uhppote"
	"uhppoted"
)

type MQTTD struct {
	Broker     string
	Topic      string
	connection MQTT.Client
	interrupt  chan os.Signal
	Debug      bool
}

type fdispatch func(*uhppoted.UHPPOTED, context.Context, uhppoted.Request)
type fdispatchx func(*MQTTD, *uhppoted.UHPPOTED, context.Context, MQTT.Message)

type dispatcher struct {
	mqttd    *MQTTD
	uhppoted *uhppoted.UHPPOTED
	uhppote  *uhppote.UHPPOTE
	log      *log.Logger
	topic    string
	table    map[string]fdispatch
	tablex   map[string]fdispatchx
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
			m.Topic + "/devices:get":             (*uhppoted.UHPPOTED).GetDevices,
			m.Topic + "/device:get":              (*uhppoted.UHPPOTED).GetDevice,
			m.Topic + "/device/status:get":       (*uhppoted.UHPPOTED).GetStatus,
			m.Topic + "/device/time:get":         (*uhppoted.UHPPOTED).GetTime,
			m.Topic + "/device/time:set":         (*uhppoted.UHPPOTED).SetTime,
			m.Topic + "/device/door/delay:get":   (*uhppoted.UHPPOTED).GetDoorDelay,
			m.Topic + "/device/door/delay:set":   (*uhppoted.UHPPOTED).SetDoorDelay,
			m.Topic + "/device/door/control:get": (*uhppoted.UHPPOTED).GetDoorControl,
			m.Topic + "/device/door/control:set": (*uhppoted.UHPPOTED).SetDoorControl,
			m.Topic + "/device/cards:get":        (*uhppoted.UHPPOTED).GetCards,
			m.Topic + "/device/cards:delete":     (*uhppoted.UHPPOTED).DeleteCards,
			m.Topic + "/device/card:get":         (*uhppoted.UHPPOTED).GetCard,
			m.Topic + "/device/card:put":         (*uhppoted.UHPPOTED).PutCard,
			m.Topic + "/device/card:delete":      (*uhppoted.UHPPOTED).DeleteCard,
			m.Topic + "/device/event:get":        (*uhppoted.UHPPOTED).GetEvent,
		},
		tablex: map[string]fdispatchx{
			m.Topic + "/device/events:get": (*MQTTD).getEvents,
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
			l.Printf("WARN: Error unsubscribing from topic' %s': %v", m.Topic, token.Error())
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

	m.interrupt = make(chan os.Signal)

	go func() {
		api.Listen(ctx, m.interrupt)
	}()

	return nil
}

func (m *MQTTD) subscribeAndServe(d *dispatcher) error {
	var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		d.dispatch(client, msg)
	}

	options := MQTT.NewClientOptions().AddBroker(m.Broker)
	options.SetClientID("twystd-uhppoted-mqttd")
	options.SetDefaultPublishHandler(f)

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

	request := Request{
		Message: msg,
	}

	if err := json.Unmarshal(msg.Payload(), &request); err != nil {
		oops(ctx, "mqtt-dispatch", "Invalid message format", uhppoted.StatusBadRequest)
		return
	}

	if fn, ok := d.table[msg.Topic()]; ok {
		fn(d.uhppoted, ctx, &request)
	} else if fnx, ok := d.tablex[msg.Topic()]; ok {
		fnx(d.mqttd, d.uhppoted, ctx, msg)
	}
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
	b, err := json.Marshal(response)
	if err != nil {
		oops(ctx, "encoding/json", "Error encoding response", uhppoted.StatusInternalServerError)
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

	token := client.Publish(topic+"/reply", 0, false, string(b))
	token.Wait()
}

func (m *MQTTD) Oops(ctx context.Context, operation string, message string, errorCode int) {
	oops(ctx, operation, message, errorCode)
}

func (m *MQTTD) OnError(ctx context.Context, operation string, message string, errorCode int, err error) {
	// u.warn(ctx, 0, "get-events", err)
	// u.oops(ctx, "get-events", "Missing/invalid device ID", StatusBadRequest)
	oops(ctx, operation, message, errorCode)
}

func oops(ctx context.Context, operation string, message string, errorCode int) {
	response := struct {
		Operation string `json:"operation"`
		Error     struct {
			Message   string `json:"message"`
			ErrorCode int    `json:"error-code"`
		} `json:"error"`
	}{
		Operation: operation,
		Error: struct {
			Message   string `json:"message"`
			ErrorCode int    `json:"error-code"`
		}{
			Message:   message,
			ErrorCode: errorCode,
		},
	}

	b, err := json.Marshal(response)
	if err != nil {
		ctx.Value("log").(*log.Logger).Printf("ERROR: Error generating JSON response (%v)", err)
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

	token := client.Publish(topic+"/gateway/errors", 0, false, string(b))
	token.Wait()
}
