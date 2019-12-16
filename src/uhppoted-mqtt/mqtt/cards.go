package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"uhppoted"
)

func (m *MQTTD) deleteCard(impl *uhppoted.UHPPOTED, ctx context.Context, msg MQTT.Message) {
	body := struct {
		DeviceID   *uint32 `json:"device-id"`
		CardNumber *uint32 `json:"card-number"`
	}{}

	if err := json.Unmarshal(msg.Payload(), &body); err != nil {
		m.OnError(ctx, "delete-card", "Cannot parse request", uhppoted.StatusBadRequest, err)
	} else if body.DeviceID == nil {
		m.OnError(ctx, "delete-card", "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
	} else if *body.DeviceID == 0 {
		m.OnError(ctx, "delete-card", "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
	} else if body.CardNumber == nil {
		m.OnError(ctx, "delete-card", "Missing/invalid card number", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid card number '%s'", string(msg.Payload())))
	} else {
		rq := uhppoted.DeleteCardRequest{
			DeviceID:   *body.DeviceID,
			CardNumber: *body.CardNumber,
		}

		if response, status, err := impl.DeleteCard(ctx, rq); err != nil {
			m.OnError(ctx, "delete-card", "Error deleting card", status, err)
		} else if response != nil {
			m.Reply(ctx, response)
		}
	}
}
