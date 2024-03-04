package mkputils

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	// Success Response Code
	Success         = 200
	StatusCreated   = 2001001
	StatusEdited    = 2001002
	StatusRemoved   = 2001003
	StatusFetchData = 2001004

	// Invalid Response Code
	InvalidFormat    = 400
	InvalidCreated   = 4001001
	InvalidEdited    = 4001002
	InvalidRemoved   = 4001003
	InvalidFetchData = 4001005

	StatusNotFound = 404
	StatusFound    = 302
)

var statusText = map[int]string{
	Success:         "successfully",
	StatusCreated:   "successfully created",
	StatusEdited:    "successfully edited",
	StatusRemoved:   "successfully removed",
	StatusFetchData: "successfully fetch data",

	InvalidFormat:    "invalid format",
	InvalidCreated:   "invalid created",
	InvalidEdited:    "invalid edited",
	InvalidRemoved:   "invalid removed",
	InvalidFetchData: "invalid fetch data",

	StatusNotFound: "not found",
	StatusFound:    "found",
}

func getResponseText(code int) string {
	return statusText[0]
}

type meta struct {
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Success         bool   `json:"success"`
	Message         string `json:"message"`
}

type (
	response struct {
		StatusCode       int         `json:"statusCode"`
		Meta             meta        `json:"meta"`
		Result           interface{} `json:"result"`
		ResponseDatetime time.Time   `json:"responseDatetime"`
	}
	optFunc func(r *response)
)

func defaultResponseOK() response {
	return response{
		StatusCode: http.StatusFound,
		Meta: meta{
			ResponseCode:    Success,
			ResponseMessage: getResponseText(Success),
			Success:         true,
			Message:         EMPTY_VALUE,
		},
		ResponseDatetime: time.Now(),
	}
}

func defaultResponseFail() response {
	return response{
		StatusCode: http.StatusBadRequest,
		Meta: meta{
			ResponseCode:    InvalidFormat,
			ResponseMessage: getResponseText(InvalidFormat),
			Success:         false,
			Message:         EMPTY_VALUE,
		},
		ResponseDatetime: time.Now(),
	}
}

func SetStatusCode(code int) optFunc {
	return func(r *response) {
		r.StatusCode = code
	}

}

func SetResponseCode(code int) optFunc {
	return func(r *response) {
		r.Meta.ResponseCode = code
		r.Meta.ResponseMessage = getResponseText(code)
	}
}

func SetMessage(message string) optFunc {
	return func(r *response) {
		r.Meta.Message = message
	}
}

func SetResult(result interface{}) optFunc {
	return func(r *response) {
		r.Result = result
	}
}

func ResponseOK(ctx echo.Context, opts ...optFunc) error {
	o := defaultResponseOK()

	for _, fn := range opts {
		fn(&o)
	}

	return ctx.JSON(o.StatusCode, &response{
		Meta: meta{
			ResponseCode:    o.Meta.ResponseCode,
			ResponseMessage: o.Meta.ResponseMessage,
			Success:         o.Meta.Success,
			Message:         o.Meta.Message,
		},
		Result:           o.Result,
		ResponseDatetime: o.ResponseDatetime,
	})
}

func ResponseFAIL(ctx echo.Context, opts ...optFunc) error {
	o := defaultResponseFail()

	for _, fn := range opts {
		fn(&o)
	}

	return ctx.JSON(o.StatusCode, &response{
		Meta: meta{
			ResponseCode:    o.Meta.ResponseCode,
			ResponseMessage: o.Meta.ResponseMessage,
			Success:         o.Meta.Success,
			Message:         o.Meta.Message,
		},
		Result:           o.Result,
		ResponseDatetime: o.ResponseDatetime,
	})
}
