package uhppoted

import (
	"context"
	"fmt"
	"uhppote"
	"uhppote/types"
)

type ACL struct {
	DeviceID uint32         `json:"device-id"`
	From     types.Date     `json:"valid-from"`
	To       types.Date     `json:"valid-until"`
	Doors    map[uint8]bool `json:"doors"`
}

type GetCardsRequest struct {
	DeviceID uint32
}

type GetCardsResponse struct {
	Device struct {
		ID    uint32   `json:"id"`
		Cards []uint32 `json:"cards"`
	} `json:"device"`
}

type GetCardRequest struct {
	DeviceID   uint32
	CardNumber uint32
}

type GetCardResponse struct {
	Card struct {
		CardNumber uint32 `json:"card-number"`
		ACL        []ACL  `json:"acl"`
	} `json:"card"`
}

type PutCardRequest struct {
	DeviceID uint32     `json:"device"`
	Card     types.Card `json:"card"`
}

type PutCardResponse struct {
	Card struct {
		CardNumber uint32 `json:"card-number"`
		ACL        []ACL  `json:"acl"`
	} `json:"card"`
}

type DeleteCardRequest struct {
	DeviceID   uint32
	CardNumber uint32
}

type DeleteCardResponse struct {
	Device Device `json:"device"`
	Card   struct {
		CardNumber uint32 `json:"card-number"`
		Deleted    bool   `json:"deleted"`
	} `json:"card"`
}

type DeleteCardsResponse struct {
	Device struct {
		ID      uint32 `json:"id"`
		Deleted bool   `json:"deleted"`
	} `json:"device"`
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
		struct {
			ID    uint32   `json:"id"`
			Cards []uint32 `json:"cards"`
		}{
			ID:    device,
			Cards: cards,
		},
	}

	u.debug(ctx, "get-cards", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}

func (u *UHPPOTED) DeleteCards(ctx context.Context, rq Request) {
	u.debug(ctx, "delete-cards", rq)

	id, err := rq.DeviceID()
	if err != nil {
		u.warn(ctx, 0, "delete-cards", err)
		u.oops(ctx, "delete-cards", "Missing/invalid device ID)", StatusBadRequest)
		return
	}

	deleted, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).DeleteCards(*id)
	if err != nil {
		u.warn(ctx, *id, "delete-cards", err)
		u.oops(ctx, "delete-cards", "Error deleting cards", StatusInternalServerError)
		return
	}

	response := DeleteCardsResponse{
		struct {
			ID      uint32 `json:"id"`
			Deleted bool   `json:"deleted"`
		}{
			ID:      *id,
			Deleted: deleted.Succeeded,
		},
	}

	u.reply(ctx, response)
}

func (u *UHPPOTED) GetCard(ctx context.Context, request GetCardRequest) (*GetCardResponse, int, error) {
	u.debug(ctx, "get-card", fmt.Sprintf("request  %v", request))

	device := request.DeviceID
	cardID := request.CardNumber

	card, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCardById(device, cardID)
	if err != nil {
		return nil, StatusInternalServerError, fmt.Errorf("Error retrieving card %v from %v", cardID, device)
	}

	if card == nil {
		return nil, StatusNotFound, fmt.Errorf("No record for card %v on %v", cardID, device)
	}

	response := GetCardResponse{
		Card: struct {
			CardNumber uint32 `json:"card-number"`
			ACL        []ACL  `json:"acl"`
		}{
			CardNumber: card.CardNumber,
			ACL: []ACL{
				ACL{
					DeviceID: device,
					From:     card.From,
					To:       card.To,
					Doors: map[uint8]bool{
						1: card.Doors[0],
						2: card.Doors[1],
						3: card.Doors[2],
						4: card.Doors[3],
					},
				},
			},
		},
	}

	u.debug(ctx, "get-card", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
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
		Card: struct {
			CardNumber uint32 `json:"card-number"`
			ACL        []ACL  `json:"acl"`
		}{
			CardNumber: card.CardNumber,
			ACL: []ACL{
				ACL{
					DeviceID: device,
					From:     card.From,
					To:       card.To,
					Doors: map[uint8]bool{
						1: card.Doors[0],
						2: card.Doors[1],
						3: card.Doors[2],
						4: card.Doors[3],
					},
				},
			},
		},
	}

	u.debug(ctx, "put-card", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
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
		Device: Device{device},
		Card: struct {
			CardNumber uint32 `json:"card-number"`
			Deleted    bool   `json:"deleted"`
		}{
			CardNumber: cardNo,
			Deleted:    deleted.Succeeded,
		},
	}

	u.debug(ctx, "delete-card", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}
