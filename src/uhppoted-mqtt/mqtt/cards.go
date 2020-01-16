package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"uhppote/types"
	"uhppoted"
)

func (m *MQTTD) getCards(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return nil, nil
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(request)))
		return nil, nil
	}

	rq := uhppoted.GetCardsRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.GetCards(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving cards", status, err)
		return nil, nil
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.GetCardsResponse
	}{
		metainfo:         meta,
		GetCardsResponse: *response,
	}, nil
}

func (m *MQTTD) deleteCards(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return nil, nil
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(request)))
		return nil, nil
	}

	rq := uhppoted.DeleteCardsRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.DeleteCards(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error deleting cards", status, err)
		return nil, nil
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.DeleteCardsResponse
	}{
		metainfo:            meta,
		DeleteCardsResponse: *response,
	}, nil
}

func (m *MQTTD) getCard(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID   *uhppoted.DeviceID `json:"device-id"`
		CardNumber *uint32            `json:"card-number"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return nil, nil
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(request)))
		return nil, nil
	}

	if body.CardNumber == nil {
		m.OnError(ctx, "Missing/invalid card number", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid card number '%s'", string(request)))
		return nil, nil
	}

	rq := uhppoted.GetCardRequest{
		DeviceID:   *body.DeviceID,
		CardNumber: *body.CardNumber,
	}

	response, status, err := impl.GetCard(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving card", status, err)
		return nil, nil
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.GetCardResponse
	}{
		metainfo:        meta,
		GetCardResponse: *response,
	}, nil
}

func (m *MQTTD) putCard(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Card     *types.Card        `json:"card"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return nil, nil
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(request)))
		return nil, nil
	}

	if body.Card == nil {
		m.OnError(ctx, "Missing/invalid card", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid card'%s'", string(request)))
		return nil, nil
	}

	rq := uhppoted.PutCardRequest{
		DeviceID: *body.DeviceID,
		Card:     *body.Card,
	}

	response, status, err := impl.PutCard(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error storing card", status, err)
		return nil, nil
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.PutCardResponse
	}{
		metainfo:        meta,
		PutCardResponse: *response,
	}, nil
}

func (m *MQTTD) deleteCard(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID   *uhppoted.DeviceID `json:"device-id"`
		CardNumber *uint32            `json:"card-number"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return nil, nil
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(request)))
		return nil, nil
	}

	if body.CardNumber == nil {
		m.OnError(ctx, "Missing/invalid card number", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid card number '%s'", string(request)))
		return nil, nil
	}

	rq := uhppoted.DeleteCardRequest{
		DeviceID:   *body.DeviceID,
		CardNumber: *body.CardNumber,
	}

	response, status, err := impl.DeleteCard(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error deleting card", status, err)
		return nil, nil
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.DeleteCardResponse
	}{
		metainfo:           meta,
		DeleteCardResponse: *response,
	}, nil
}
