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

type Header struct {
	UID      float64 `json:"uid"`
	TID      float64 `json:"tid"`
	PID      float64 `json:"pid"`
	OID      float64 `json:"oid"`
	RID      float64 `json:"rid"`
	Username string  `json:"username"`
}
