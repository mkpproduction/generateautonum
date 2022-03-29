package mkputils

import "time"

type Response struct {
	Meta             Meta        `json:"meta"`
	Result           interface{} `json:"result"`
	ResponseDatetime time.Time   `json:"responseDatetime"`
}

type Meta struct {
	Code          string   `json:"code"`
	Success       bool     `json:"success"`
	Message       string   `json:"message"`
	AdditionalMsg []string `json:"additionalMsg"`
}
