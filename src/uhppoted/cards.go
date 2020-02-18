package uhppoted

import (
	"fmt"
	"uhppote/types"
)

type GetCardsRequest struct {
	DeviceID DeviceID
}

type GetCardsResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Cards    []uint32 `json:"cards"`
}

func (u *UHPPOTED) GetCards(request GetCardsRequest) (*GetCardsResponse, error) {
	u.debug("get-cards", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)

	N, err := u.Uhppote.GetCards(device)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error retrieving cards from %v (%w)", device, err))
	}

	cards := make([]uint32, 0)

	for index := uint32(0); index < N.Records; index++ {
		record, err := u.Uhppote.GetCardByIndex(device, index+1)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error retrieving cards from %v (%w)", device, err))
		}

		cards = append(cards, record.CardNumber)
	}

	response := GetCardsResponse{
		DeviceID: DeviceID(N.SerialNumber),
		Cards:    cards,
	}

	u.debug("get-cards", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type DeleteCardsRequest struct {
	DeviceID DeviceID
}

type DeleteCardsResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Deleted  bool     `json:"deleted"`
}

func (u *UHPPOTED) DeleteCards(request DeleteCardsRequest) (*DeleteCardsResponse, error) {
	u.debug("delete-cards", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)

	deleted, err := u.Uhppote.DeleteCards(device)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error deleting cards from %v (%w)", device, err))
	}

	response := DeleteCardsResponse{
		DeviceID: DeviceID(deleted.SerialNumber),
		Deleted:  deleted.Succeeded,
	}

	u.debug("delete-cards", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type GetCardRequest struct {
	DeviceID   DeviceID
	CardNumber uint32
}

type GetCardResponse struct {
	DeviceID DeviceID   `json:"device-id"`
	Card     types.Card `json:"card"`
}

func (u *UHPPOTED) GetCard(request GetCardRequest) (*GetCardResponse, error) {
	u.debug("get-card", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	cardID := request.CardNumber

	card, err := u.Uhppote.GetCardByID(device, cardID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Sprintf("Error retrieving card %v from %v (%w)", card.CardNumber, device, err))
	}

	if card == nil {
		return nil, fmt.Errorf("%w: %v", NotFound, fmt.Errorf("Error retrieving card %v from %v", card.CardNumber, device))
	}

	response := GetCardResponse{
		DeviceID: DeviceID(device),
		Card:     *card,
	}

	u.debug("get-card", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type PutCardRequest struct {
	DeviceID DeviceID
	Card     types.Card
}

type PutCardResponse struct {
	DeviceID DeviceID   `json:"device-id"`
	Card     types.Card `json:"card"`
}

func (u *UHPPOTED) PutCard(request PutCardRequest) (*PutCardResponse, error) {
	u.debug("put-card", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	card := request.Card

	authorised, err := u.Uhppote.PutCard(device, card)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error storing card %v to %v (%w)", card.CardNumber, device, err))
	}

	if !authorised.Succeeded {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error storing card %v to %v (%w)", card.CardNumber, device, err))
	}

	response := PutCardResponse{
		DeviceID: DeviceID(authorised.SerialNumber),
		Card:     card,
	}

	u.debug("put-card", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type DeleteCardRequest struct {
	DeviceID   DeviceID
	CardNumber uint32
}

type DeleteCardResponse struct {
	DeviceID   DeviceID `json:"device-id"`
	CardNumber uint32   `json:"card-number"`
	Deleted    bool     `json:"deleted"`
}

func (u *UHPPOTED) DeleteCard(request DeleteCardRequest) (*DeleteCardResponse, error) {
	u.debug("delete-card", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	cardNo := request.CardNumber

	deleted, err := u.Uhppote.DeleteCard(device, cardNo)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error deleting card %v from %v (%w)", cardNo, device, err))
	}

	response := DeleteCardResponse{
		DeviceID:   DeviceID(deleted.SerialNumber),
		CardNumber: cardNo,
		Deleted:    deleted.Succeeded,
	}

	u.debug("delete-card", fmt.Sprintf("response %+v", response))

	return &response, nil
}
