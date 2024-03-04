package mkputils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"regexp"
	"strings"
)

func HandleSignatureMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") != echo.MIMEApplicationJSON {
			return next(c)
		}

		data, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}

		dt := make(map[string]interface{})
		if err := json.Unmarshal(data, &dt); err != nil {
			return err
		}

		if err := validateSignature(c, dt); err != nil {
			return err
		}

		c.Logger().Error(string(data))
		c.Request().Body = ioutil.NopCloser(bytes.NewReader(data))

		return next(c)
	}
}

func validateSignature(ctx echo.Context, body interface{}) error {
	mode := GetEnv("PRODUCTION_MODE", "PRODUCTION")

	if mode == "PRODUCTION" || mode == "STAGING" {
		type message struct {
			Authorization    string `json:"authorization"`
			Method           string `json:"method"`
			Path             string `json:"path"`
			URL              string `json:"url"`
			Timestamp        string `json:"timestamp"`
			RequestSignature string `json:"requestSignature"`
		}

		msg := message{
			Authorization:    ctx.Request().Header.Get("Authorization"),
			Method:           ctx.Request().Method,
			Path:             ctx.Request().RequestURI,
			Timestamp:        ctx.Request().Header.Get("X-TIMESTAMP"),
			RequestSignature: ctx.Request().Header.Get("X-SIGNATURE"),
		}

		basket := make(map[string]interface{})
		requestBodyStr, _ := json.Marshal(body)
		err := json.Unmarshal(requestBodyStr, &basket)
		if err != nil {
			return err
		}

		sKey1 := strings.ReplaceAll(ToString(basket), `"`, "")
		sKey2 := strings.ToLower(sKey1)
		sKey3 := strings.Trim(sKey2, " ")
		regx := regexp.MustCompile("[^a-zA-Z0-9{}:.,]")
		sKey4 := regx.ReplaceAllLiteralString(sKey3, "")

		if msg.Timestamp == EMPTY_VALUE {
			return errors.New("invalid timestamp")
		}

		if msg.RequestSignature == EMPTY_VALUE {
			return errors.New("invalid signature")
		}

		sValue := fmt.Sprintf("%s:%s:%s:%s", msg.Path, msg.Method, msg.Timestamp, sKey4)

		mac := hmac.New(sha512.New, []byte(msg.Authorization))
		mac.Write([]byte(sValue))
		mac.BlockSize()
		mac.Size()
		expectationMac := hex.EncodeToString(mac.Sum(nil))
		mac.Reset()

		if !hmac.Equal([]byte(msg.RequestSignature), []byte(expectationMac)) {
			return errors.New("invalid format signature")
		}
	}

	return nil
}
