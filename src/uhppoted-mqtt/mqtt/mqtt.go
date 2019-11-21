package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"uhppote"
	"uhppoted"
)

type MQTTD struct {
	Broker     string
	Topic      string
	connection MQTT.Client
}

type Request struct {
	Message MQTT.Message
	Device  struct {
		ID uint32 `json:"id"`
	} `json:"device"`
}

type dispatcher struct {
	uhppoted *uhppoted.UHPPOTED
	uhppote  *uhppote.UHPPOTE
	log      *log.Logger
	topic    string
}

func (m *MQTTD) Run(u *uhppote.UHPPOTE, l *log.Logger) {
	d := dispatcher{
		uhppoted: &uhppoted.UHPPOTED{
			Service: m,
		},
		uhppote: u,
		log:     l,
		topic:   m.Topic,
	}

	if err := m.listenAndServe(&d); err != nil {
		l.Printf("ERROR: Error connecting to '%s': %v", m.Broker, err)
		m.Close(l)
		return
	}

	log.Printf("... connected to %s\n", m.Broker)
}

func (m *MQTTD) Close(l *log.Logger) {
	if m.connection != nil {
		log.Printf("... closing connection to %s", m.Broker)
		token := m.connection.Unsubscribe(m.Topic + "/#")
		if token.Wait() && token.Error() != nil {
			l.Printf("WARN: Error unsubscribing from topic' %s': %v", "twystd-uhppote", token.Error())
		}

		m.connection.Disconnect(250)
	}

	m.connection = nil
}

func (m *MQTTD) listenAndServe(d *dispatcher) error {
	//	MQTT.DEBUG = log.New(os.Stdout, "", 0)
	//	MQTT.WARN = log.New(os.Stdout, "", 0)
	//	MQTT.ERROR = log.New(os.Stdout, "", 0)
	//	MQTT.CRITICAL = log.New(os.Stdout, "", 0)

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

	token = m.connection.Subscribe(m.Topic+"/gateway/#", 0, nil)
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
	request := Request{Message: msg}

	if err := json.Unmarshal(msg.Payload(), &request); err != nil {
		oops(ctx, "get-device", "Invalid message format", uhppoted.StatusBadRequest)
		return
	}

	switch msg.Topic() {
	case d.topic + "/gateway/ping":
		d.uhppoted.GetDevices(ctx, &request)

	case d.topic + "/gateway/device/ping":
		d.uhppoted.GetDevice(ctx, &request)

	case d.topic + "/gateway/device/status":
		d.uhppoted.GetStatus(ctx, &request)
	}
}

func (m *MQTTD) Reply(ctx context.Context, response interface{}) {
	b, err := json.Marshal(response)
	if err != nil {
		oops(ctx, "encoding/json", "Error generating response", uhppoted.StatusInternalServerError)
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

	token := client.Publish(topic+"/devices/ping", 0, false, string(b))
	token.Wait()
}

func (m *MQTTD) Oops(ctx context.Context, operation string, message string, errorCode int) {
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

func (rq *Request) DeviceId() (uint32, error) {
	if rq.Device.ID == 0 {
		return 0, errors.New("Missing device ID")
	}

	return rq.Device.ID, nil
}
