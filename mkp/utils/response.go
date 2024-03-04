package mkputils

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type meta struct {
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Success         bool   `json:"success"`
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
			ResponseCode:    http.StatusOK,
			ResponseMessage: http.StatusText(http.StatusOK),
			Success:         true,
			Message:         EMPTY_VALUE,
		},
		ResponseDatetime: time.Now(),
	}
}

func defaultResponseFail() response {
	return response{
		Meta: meta{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: http.StatusText(http.StatusBadRequest),
			Success:         false,
			Message:         EMPTY_VALUE,
		},
		ResponseDatetime: time.Now(),
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

	return ctx.JSON(o.Meta.ResponseCode, &response{
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

	return ctx.JSON(o.Meta.ResponseCode, &response{
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
