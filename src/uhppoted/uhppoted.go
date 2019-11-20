package uhppoted

import (
	"context"
	"log"
	"net/http"
)

const (
	//	StatusOK              = http.StatusOK
	//	StatusBadRequest      = http.StatusBadRequest
	//	StatusUnauthorized    = http.StatusUnauthorized
	//	StatusForbidden       = http.StatusForbidden
	//	StatusNotFound        = http.StatusNotFound
	//	StatusRequestTimeout  = http.StatusRequestTimeout
	//	StatusTooManyRequests = http.StatusTooManyRequests

	StatusInternalServerError = http.StatusInternalServerError

//	StatusNotImplemented      = http.StatusNotImplemented
//	StatusBadGateway          = http.StatusBadGateway
//	StatusServiceUnavailable  = http.StatusServiceUnavailable
//	StatusGatewayTimeout      = http.StatusGatewayTimeout
)

type Service interface {
	Reply(ctx context.Context, response interface{})
	Oops(ctx context.Context, operation string, message string, errorCode int)
}

type Request interface {
}

type UHPPOTED struct {
	Service Service
}

func (u *UHPPOTED) debug(ctx context.Context, deviceId int, operation string, rq Request) {
	ctx.Value("log").(*log.Logger).Printf("DEBUG %-12d %-20s %#v\n", deviceId, operation, rq)
}

func (u *UHPPOTED) warn(ctx context.Context, deviceId uint32, operation string, err error) {
	ctx.Value("log").(*log.Logger).Printf("WARN  %-12d %-20s %v\n", deviceId, operation, err)
}

func (u *UHPPOTED) reply(ctx context.Context, response interface{}) {
	u.Service.Reply(ctx, response)
}

func (u *UHPPOTED) oops(ctx context.Context, operation string, message string, errorCode int) {
	u.Service.Oops(ctx, operation, message, errorCode)
}
