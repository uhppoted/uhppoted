package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"uhppote/types"
	"uhppoted"
)

func (m *MQTTD) getCards(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
		}
	}

	if body.DeviceID == nil {
		return nil, &errorx{
			Err:     errors.New("Missing device ID"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing device ID",
		}
	}

	rq := uhppoted.GetCardsRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.GetCards(ctx, rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: fmt.Sprintf("Error retrieving card list for %v", *body.DeviceID),
		}
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
		return nil, &errorx{
			Err:     err,
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
		}
	}

	if body.DeviceID == nil {
		return nil, &errorx{
			Err:     errors.New("Missing device ID"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing device ID",
		}
	}

	rq := uhppoted.DeleteCardsRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.DeleteCards(ctx, rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: fmt.Sprintf("Error deleting card list for %v", *body.DeviceID),
		}
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
		return nil, &errorx{
			Err:     err,
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
		}
	}

	if body.DeviceID == nil {
		return nil, &errorx{
			Err:     errors.New("Missing device ID"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing device ID",
		}
	}

	if body.CardNumber == nil {
		return nil, &errorx{
			Err:     errors.New("Missing/invalid card number"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid card number",
		}
	}

	rq := uhppoted.GetCardRequest{
		DeviceID:   *body.DeviceID,
		CardNumber: *body.CardNumber,
	}

	response, status, err := impl.GetCard(ctx, rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: fmt.Sprintf("Error retrieving card %v", *body.CardNumber),
		}
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
		return nil, &errorx{
			Err:     err,
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
		}
	}

	if body.DeviceID == nil {
		return nil, &errorx{
			Err:     errors.New("Missing device ID"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing device ID",
		}
	}

	if body.Card == nil {
		return nil, &errorx{
			Err:     errors.New("Missing/invalid card number"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid card number",
		}
	}

	rq := uhppoted.PutCardRequest{
		DeviceID: *body.DeviceID,
		Card:     *body.Card,
	}

	response, status, err := impl.PutCard(ctx, rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: fmt.Sprintf("Error storing card %v", body.Card.CardNumber),
		}
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
		return nil, &errorx{
			Err:     err,
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
		}
	}

	if body.DeviceID == nil {
		return nil, &errorx{
			Err:     errors.New("Missing device ID"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing device ID",
		}
	}

	if body.CardNumber == nil {
		return nil, &errorx{
			Err:     errors.New("Missing/invalid card number"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid card number",
		}
	}

	rq := uhppoted.DeleteCardRequest{
		DeviceID:   *body.DeviceID,
		CardNumber: *body.CardNumber,
	}

	response, status, err := impl.DeleteCard(ctx, rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: fmt.Sprintf("Error deleting card %v", *body.CardNumber),
		}
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
