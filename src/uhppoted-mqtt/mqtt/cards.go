package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"uhppote/types"
	"uhppoted"
)

func (m *MQTTD) getCards(impl *uhppoted.UHPPOTED, ctx context.Context, msg MQTT.Message) {
	operation := "get-cards"
	body := struct {
		DeviceID *uint32 `json:"device-id"`
	}{}

	if err := json.Unmarshal(msg.Payload(), &body); err != nil {
		m.OnError(ctx, operation, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return
	}

	if body.DeviceID == nil || *body.DeviceID == 0 {
		m.OnError(ctx, operation, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
		return
	}

	rq := uhppoted.GetCardsRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.GetCards(ctx, rq)
	if err != nil {
		m.OnError(ctx, operation, "Error retrieving cards", status, err)
		return
	}

	if response != nil {
		reply := struct {
			MetaInfo *metainfo `json:"meta-info,omitempty"`
			uhppoted.GetCardsResponse
		}{
			MetaInfo:         getMetaInfo(ctx, operation),
			GetCardsResponse: *response,
		}

		m.Reply(ctx, reply)
	}
}

func (m *MQTTD) deleteCards(impl *uhppoted.UHPPOTED, ctx context.Context, msg MQTT.Message) {
	operation := "delete-cards"
	body := struct {
		DeviceID *uint32 `json:"device-id"`
	}{}

	if err := json.Unmarshal(msg.Payload(), &body); err != nil {
		m.OnError(ctx, operation, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return
	}

	if body.DeviceID == nil || *body.DeviceID == 0 {
		m.OnError(ctx, operation, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
		return
	}

	rq := uhppoted.DeleteCardsRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.DeleteCards(ctx, rq)
	if err != nil {
		m.OnError(ctx, operation, "Error deleting cards", status, err)
		return
	}

	if response != nil {
		reply := struct {
			MetaInfo *metainfo `json:"meta-info,omitempty"`
			uhppoted.DeleteCardsResponse
		}{
			MetaInfo:            getMetaInfo(ctx, operation),
			DeleteCardsResponse: *response,
		}

		m.Reply(ctx, reply)
	}
}

func (m *MQTTD) getCard(impl *uhppoted.UHPPOTED, ctx context.Context, msg MQTT.Message) {
	operation := "get-card"
	body := struct {
		DeviceID   *uint32 `json:"device-id"`
		CardNumber *uint32 `json:"card-number"`
	}{}

	if err := json.Unmarshal(msg.Payload(), &body); err != nil {
		m.OnError(ctx, operation, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return
	}

	if body.DeviceID == nil || *body.DeviceID == 0 {
		m.OnError(ctx, operation, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
		return
	}

	if body.CardNumber == nil {
		m.OnError(ctx, operation, "Missing/invalid card number", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid card number '%s'", string(msg.Payload())))
		return
	}

	rq := uhppoted.GetCardRequest{
		DeviceID:   *body.DeviceID,
		CardNumber: *body.CardNumber,
	}

	response, status, err := impl.GetCard(ctx, rq)
	if err != nil {
		m.OnError(ctx, operation, "Error retrieving card", status, err)
		return
	}

	if response != nil {
		reply := struct {
			MetaInfo *metainfo `json:"meta-info,omitempty"`
			uhppoted.GetCardResponse
		}{
			MetaInfo:        getMetaInfo(ctx, operation),
			GetCardResponse: *response,
		}

		m.Reply(ctx, reply)
	}
}

func (m *MQTTD) putCard(impl *uhppoted.UHPPOTED, ctx context.Context, msg MQTT.Message) {
	operation := "put-card"
	body := struct {
		DeviceID *uint32     `json:"device-id"`
		Card     *types.Card `json:"card"`
	}{}

	if err := json.Unmarshal(msg.Payload(), &body); err != nil {
		m.OnError(ctx, operation, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return
	}

	if body.DeviceID == nil || *body.DeviceID == 0 {
		m.OnError(ctx, operation, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
		return
	}

	if body.Card == nil {
		m.OnError(ctx, operation, "Missing/invalid card", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid card'%s'", string(msg.Payload())))
		return
	}

	rq := uhppoted.PutCardRequest{
		DeviceID: *body.DeviceID,
		Card:     *body.Card,
	}

	response, status, err := impl.PutCard(ctx, rq)
	if err != nil {
		m.OnError(ctx, operation, "Error storing card", status, err)
	}

	if response != nil {
		reply := struct {
			MetaInfo *metainfo `json:"meta-info,omitempty"`
			uhppoted.PutCardResponse
		}{
			MetaInfo:        getMetaInfo(ctx, operation),
			PutCardResponse: *response,
		}

		m.Reply(ctx, reply)
	}
}

func (m *MQTTD) deleteCard(impl *uhppoted.UHPPOTED, ctx context.Context, msg MQTT.Message) {
	operation := "delete-card"
	body := struct {
		DeviceID   *uint32 `json:"device-id"`
		CardNumber *uint32 `json:"card-number"`
	}{}

	if err := json.Unmarshal(msg.Payload(), &body); err != nil {
		m.OnError(ctx, operation, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return
	}

	if body.DeviceID == nil || *body.DeviceID == 0 {
		m.OnError(ctx, operation, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
		return
	}

	if body.CardNumber == nil {
		m.OnError(ctx, operation, "Missing/invalid card number", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid card number '%s'", string(msg.Payload())))
		return
	}
	rq := uhppoted.DeleteCardRequest{
		DeviceID:   *body.DeviceID,
		CardNumber: *body.CardNumber,
	}

	response, status, err := impl.DeleteCard(ctx, rq)
	if err != nil {
		m.OnError(ctx, operation, "Error deleting card", status, err)
		return
	}

	if response != nil {
		reply := struct {
			MetaInfo *metainfo `json:"meta-info,omitempty"`
			uhppoted.DeleteCardResponse
		}{
			MetaInfo:           getMetaInfo(ctx, operation),
			DeleteCardResponse: *response,
		}

		m.Reply(ctx, reply)
	}
}
