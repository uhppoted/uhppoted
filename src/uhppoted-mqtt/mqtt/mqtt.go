package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
	"uhppote"
	"uhppote/types"
	"uhppoted"
)

type MQTTD struct {
	Broker     string
	Topic      string
	connection MQTT.Client
}

type Request struct {
	Message MQTT.Message
}

type fdispatch func(*uhppoted.UHPPOTED, context.Context, uhppoted.Request)

type dispatcher struct {
	uhppoted *uhppoted.UHPPOTED
	uhppote  *uhppote.UHPPOTE
	log      *log.Logger
	topic    string
	table    map[string]fdispatch
}

func (m *MQTTD) Run(u *uhppote.UHPPOTE, l *log.Logger) {
	d := dispatcher{
		uhppoted: &uhppoted.UHPPOTED{
			Service: m,
		},
		uhppote: u,
		log:     l,
		topic:   m.Topic,
		table: map[string]fdispatch{
			m.Topic + "/devices:get":           (*uhppoted.UHPPOTED).GetDevices,
			m.Topic + "/device:get":            (*uhppoted.UHPPOTED).GetDevice,
			m.Topic + "/device/status:get":     (*uhppoted.UHPPOTED).GetStatus,
			m.Topic + "/device/time:get":       (*uhppoted.UHPPOTED).GetTime,
			m.Topic + "/device/time:set":       (*uhppoted.UHPPOTED).SetTime,
			m.Topic + "/device/door/delay:get": (*uhppoted.UHPPOTED).GetDoorDelay,
			m.Topic + "/device/door/delay:set": (*uhppoted.UHPPOTED).SetDoorDelay,
		},
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
			l.Printf("WARN: Error unsubscribing from topic' %s': %v", m.Topic, token.Error())
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

func (rq Request) String() string {
	return rq.Message.Topic() + "  " + string(rq.Message.Payload())
}

func (rq *Request) DeviceId() (*uint32, error) {
	body := struct {
		DeviceID *uint32 `json:"device-id"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, err
	} else if body.DeviceID == nil {
		return nil, errors.New("Missing device ID")
	} else if *body.DeviceID == 0 {
		return nil, errors.New("Missing device ID")
	}

	return body.DeviceID, nil
}

func (rq *Request) DateTime() (*time.Time, error) {
	body := struct {
		DateTime *types.DateTime `json:"datetime"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, err
	} else if body.DateTime == nil {
		return nil, errors.New("Missing date/time")
	}

	return (*time.Time)(body.DateTime), nil
}

func (rq *Request) Door() (*uint8, error) {
	body := struct {
		Door *uint8 `json:"door"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, err
	} else if body.Door == nil {
		return nil, errors.New("Invalid door")
	} else if *body.Door < 1 || *body.Door > 4 {
		return nil, errors.New("Invalid door")
	}

	return body.Door, nil
}
func (rq *Request) DoorDelay() (*uint8, *uint8, error) {
	body := struct {
		Door  *uint8 `json:"door"`
		Delay *uint8 `json:"delay"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, nil, err
	} else if body.Door == nil {
		return nil, nil, errors.New("Invalid door")
	} else if *body.Door < 1 || *body.Door > 4 {
		return nil, nil, errors.New("Invalid door")
	} else if body.Delay == nil {
		return nil, nil, errors.New("Invalid door delay")
	} else if *body.Delay == 0 || *body.Delay > 60 {
		return nil, nil, errors.New("Invalid door delay")
	}

	return body.Door, body.Delay, nil
}
