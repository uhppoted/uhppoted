package mqtt

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"uhppote"
)

type MQTTD struct {
	Server     string
	connection MQTT.Client
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func (m *MQTTD) Run(u *uhppote.UHPPOTE, l *log.Logger) {
	if err := m.listenAndServe(); err != nil {
		l.Printf("ERROR: Error connecting to '%s': %v", m.Server, err)
		m.Close(l)
		return
	}

	log.Printf("... connected to %s\n", m.Server)
}

func (m *MQTTD) Close(l *log.Logger) {
	if m.connection != nil {
		log.Printf("... closing connection to %s", m.Server)
		token := m.connection.Unsubscribe("twystd-uhppote")
		if token.Wait() && token.Error() != nil {
			l.Printf("WARN: Error unsubscribing from topic' %s': %v", "twystd-uhppote", token.Error())
		}

		m.connection.Disconnect(250)
	}
}

func (m *MQTTD) listenAndServe() error {
	//	MQTT.DEBUG = log.New(os.Stdout, "", 0)
	//	MQTT.WARN = log.New(os.Stdout, "", 0)
	//	MQTT.ERROR = log.New(os.Stdout, "", 0)
	//	MQTT.CRITICAL = log.New(os.Stdout, "", 0)

	options := MQTT.NewClientOptions().AddBroker(m.Server)
	options.SetClientID("twystd-uhppoted-mqttd")
	options.SetDefaultPublishHandler(f)

	m.connection = MQTT.NewClient(options)
	token := m.connection.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	token = m.connection.Subscribe("twystd-uhppote", 0, nil)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
