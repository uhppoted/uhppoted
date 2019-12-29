package uhppoted

import (
	"context"
	"fmt"
	"uhppote"
	"uhppote/types"
)

type GetCardsRequest struct {
	DeviceID uint32
}

type GetCardsResponse struct {
	DeviceID uint32   `json:"device-id"`
	Cards    []uint32 `json:"cards"`
}

func (u *UHPPOTED) GetCards(ctx context.Context, request GetCardsRequest) (*GetCardsResponse, int, error) {
	u.debug(ctx, "get-cards", fmt.Sprintf("request  %v", request))

	device := request.DeviceID

	N, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCards(device)
	if err != nil {
		return nil, StatusInternalServerError, fmt.Errorf("Error retrieving cards from %v", device)
	}

	cards := make([]uint32, 0)

	for index := uint32(0); index < N.Records; index++ {
		record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCardByIndex(device, index+1)
		if err != nil {
			return nil, StatusInternalServerError, fmt.Errorf("Error retrieving cards from %v", device)
		}

		cards = append(cards, record.CardNumber)
	}

	response := GetCardsResponse{
		DeviceID: device,
		Cards:    cards,
	}

	u.debug(ctx, "get-cards", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}

type DeleteCardsRequest struct {
	DeviceID uint32
}

type DeleteCardsResponse struct {
	DeviceID uint32 `json:"device-id"`
	Deleted  bool   `json:"deleted"`
}

func (u *UHPPOTED) DeleteCards(ctx context.Context, request DeleteCardsRequest) (*DeleteCardsResponse, int, error) {
	u.debug(ctx, "delete-cards", fmt.Sprintf("request  %v", request))

	device := request.DeviceID

	deleted, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).DeleteCards(device)
	if err != nil {
		return nil, StatusInternalServerError, fmt.Errorf("Error deleting cards from %v", device)
	}

	response := DeleteCardsResponse{
		DeviceID: device,
		Deleted:  deleted.Succeeded,
	}

	u.debug(ctx, "delete-cards", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}

type GetCardRequest struct {
	DeviceID   uint32
	CardNumber uint32
}

type GetCardResponse struct {
	DeviceID uint32     `json:"device-id"`
	Card     types.Card `json:"card"`
}

func (u *UHPPOTED) GetCard(ctx context.Context, request GetCardRequest) (*GetCardResponse, int, error) {
	u.debug(ctx, "get-card", fmt.Sprintf("request  %v", request))

	device := request.DeviceID
	cardID := request.CardNumber

	card, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCardByID(device, cardID)
	if err != nil {
		return nil, StatusInternalServerError, fmt.Errorf("Error retrieving card %v from %v", cardID, device)
	}

	if card == nil {
		return nil, StatusNotFound, fmt.Errorf("No record for card %v on %v", cardID, device)
	}

	response := GetCardResponse{
		DeviceID: device,
		Card:     *card,
	}

	u.debug(ctx, "get-card", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}

type PutCardRequest struct {
	DeviceID uint32
	Card     types.Card
}

type PutCardResponse struct {
	DeviceID uint32     `json:"device-id"`
	Card     types.Card `json:"card"`
}

func (u *UHPPOTED) PutCard(ctx context.Context, request PutCardRequest) (*PutCardResponse, int, error) {
	u.debug(ctx, "put-card", fmt.Sprintf("request  %v", request))

	device := request.DeviceID
	card := request.Card

	authorised, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).PutCard(device, card)
	if err != nil {
		return nil, StatusInternalServerError, fmt.Errorf("Error storing card %v to %v", card.CardNumber, device)
	}

	if !authorised.Succeeded {
		return nil, StatusInternalServerError, fmt.Errorf("Error storing card %v to %v", card.CardNumber, device)
	}

	response := PutCardResponse{
		DeviceID: device,
		Card:     card,
	}

	u.debug(ctx, "put-card", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}

type DeleteCardRequest struct {
	DeviceID   uint32
	CardNumber uint32
}

type DeleteCardResponse struct {
	DeviceID   uint32 `json:"device-id"`
	CardNumber uint32 `json:"card-number"`
	Deleted    bool   `json:"deleted"`
}

func (u *UHPPOTED) DeleteCard(ctx context.Context, request DeleteCardRequest) (*DeleteCardResponse, int, error) {
	u.debug(ctx, "delete-card", fmt.Sprintf("request  %v", request))

	device := request.DeviceID
	cardNo := request.CardNumber

	deleted, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).DeleteCard(device, cardNo)
	if err != nil {
		return nil, StatusInternalServerError, fmt.Errorf("Error deleting card %v from %v", cardNo, device)
	}

	response := DeleteCardResponse{
		DeviceID:   device,
		CardNumber: cardNo,
		Deleted:    deleted.Succeeded,
	}

	u.debug(ctx, "delete-card", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}
