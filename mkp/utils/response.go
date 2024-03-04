package mkputils

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	// Success Response Code
	Success   = 200
	Created   = 1001
	Edited    = 1002
	Removed   = 1003
	FetchData = 1004

	Successfully    = 200200
	StatusCreated   = 2001001
	StatusEdited    = 2001002
	StatusRemoved   = 2001003
	StatusFetchData = 2001004

	// Invalid Response Code
	InvalidEntity    = 400400
	InvalidFormat    = 400
	InvalidCreated   = 4001001
	InvalidEdited    = 4001002
	InvalidRemoved   = 4001003
	InvalidFetchData = 4001005

	StatusNotFound = 404
	StatusFound    = 302
)

var statusText = map[int]string{
	Created:         "created",
	Edited:          "edited",
	Removed:         "removed",
	FetchData:       "fetch data",
	Successfully:    "successfully",
	StatusCreated:   "successfully created",
	StatusEdited:    "successfully edited",
	StatusRemoved:   "successfully removed",
	StatusFetchData: "successfully fetch data",

	InvalidEntity:    "unprocessable entity",
	InvalidFormat:    "invalid format",
	InvalidCreated:   "invalid created",
	InvalidEdited:    "invalid edited",
	InvalidRemoved:   "invalid removed",
	InvalidFetchData: "invalid fetch data",

	StatusNotFound: "not found",
	StatusFound:    "found",
}

func getResponseText(statusCode int, code int) string {
	str := fmt.Sprintf("%d%d", statusCode, code)
	strFtm, _ := strconv.Atoi(str)

	log.Println("strFmt:", strFtm)

	return statusText[strFtm]
}

type meta struct {
	Success         bool   `json:"success"`
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Message         string `json:"message"`
}

type (
	response struct {
		Meta             meta        `json:"meta"`
		Result           interface{} `json:"result"`
		ResponseDatetime time.Time   `json:"responseDatetime"`
	}
	optFunc func(r *response)
)

func defaultResponseOK() response {
	return response{
		Meta: meta{
			Success:      true,
			ResponseCode: Success,
			Message:      EMPTY_VALUE,
		},
		ResponseDatetime: time.Now(),
	}
}

func defaultResponseFail() response {
	return response{
		Meta: meta{
			Success:      false,
			ResponseCode: InvalidFormat,
			Message:      EMPTY_VALUE,
		},
		ResponseDatetime: time.Now(),
	}
}

func Code(code int) optFunc {
	return func(r *response) {
		r.Meta.ResponseCode = code
	}
}

func Message(message string) optFunc {
	return func(r *response) {
		r.Meta.Message = message
	}
}

func Result(result interface{}) optFunc {
	return func(r *response) {
		r.Result = result
	}
}

func ResponseOK(ctx echo.Context, opts ...optFunc) error {
	o := defaultResponseOK()

	for _, fn := range opts {
		fn(&o)
	}

	return ctx.JSON(http.StatusOK, &response{
		Meta: meta{
			Success:         o.Meta.Success,
			ResponseCode:    o.Meta.ResponseCode,
			ResponseMessage: getResponseText(http.StatusOK, o.Meta.ResponseCode),
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

	return ctx.JSON(http.StatusBadRequest, &response{
		Meta: meta{
			Success:         o.Meta.Success,
			ResponseCode:    o.Meta.ResponseCode,
			ResponseMessage: getResponseText(http.StatusBadRequest, o.Meta.ResponseCode),
			Message:         o.Meta.Message,
		},
		Result:           o.Result,
		ResponseDatetime: o.ResponseDatetime,
	})
}
