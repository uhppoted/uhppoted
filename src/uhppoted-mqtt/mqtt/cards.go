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
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	rq := uhppoted.GetCardsRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.GetCards(rq)
	if err != nil {
		return nil, ferror(err, fmt.Sprintf("Error retrieving card list for %v", *body.DeviceID))
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
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	rq := uhppoted.DeleteCardsRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.DeleteCards(rq)
	if err != nil {
		return nil, ferror(err, fmt.Sprintf("Error deleting card list for %v", *body.DeviceID))
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
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	if body.CardNumber == nil {
		return nil, InvalidCardNumber
	}

	rq := uhppoted.GetCardRequest{
		DeviceID:   *body.DeviceID,
		CardNumber: *body.CardNumber,
	}

	response, err := impl.GetCard(rq)
	if err != nil {
		return nil, ferror(err, fmt.Sprintf("Error retrieving card %v", *body.CardNumber))
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
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	if body.Card == nil {
		return nil, InvalidCardNumber
	}

	rq := uhppoted.PutCardRequest{
		DeviceID: *body.DeviceID,
		Card:     *body.Card,
	}

	response, err := impl.PutCard(rq)
	if err != nil {
		return nil, ferror(err, fmt.Sprintf("Error storing card %v", body.Card.CardNumber))
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
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	if body.CardNumber == nil {
		return nil, InvalidCardNumber
	}

	rq := uhppoted.DeleteCardRequest{
		DeviceID:   *body.DeviceID,
		CardNumber: *body.CardNumber,
	}

	response, err := impl.DeleteCard(rq)
	if err != nil {
		return nil, ferror(err, fmt.Sprintf("Error deleting card %v", *body.CardNumber))
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
