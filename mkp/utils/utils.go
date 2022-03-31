package mkputils

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// DatetimeNow
func DatetimeNow() string {
	return time.Now().Format("20060102150405")
}

func DateNow() string {
	return time.Now().Format("20060102")
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func ValBlankOrNull(request interface{}, keyName ...string) error {
	var params interface{}
	_ = json.Unmarshal([]byte(ToString(request)), &params)
	paramsValue := params.(map[string]interface{})

	for idx := range keyName {
		name := keyName[idx]
		if len(strings.TrimSpace(paramsValue[name].(string))) == 0 {
			return errors.New(fmt.Sprintf("%s must be filled", name))
		}
	}

	return nil
}

func InArray(v interface{}, in interface{}) (ok bool, i int) {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for ; i < val.Len(); i++ {
			if ok = v == val.Index(i).Interface(); ok {
				return
			}
		}
	}
	return
}

func BindValidateStruct(ctx echo.Context, i interface{}) error {
	if err := ctx.Bind(i); err != nil {
		return err
	}

	if err := ctx.Validate(i); err != nil {
		return err
	}
	return nil
}

// Make hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Check hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func DBTransaction(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Rollback Panic
		} else if err != nil {
			tx.Rollback() // err is not nill
		} else {
			err = tx.Commit() // err is nil
		}
	}()
	err = txFunc(tx)
	return err
}

func ToString(i interface{}) string {
	log, _ := json.Marshal(i)
	logString := string(log)

	return logString
}

func ResponseJSON(success bool, code string, msg string, result interface{}, addMsg ...string) Response {
	tm := time.Now()
	response := Response{
		Meta: Meta{
			Code:          code,
			Success:       success,
			Message:       msg,
			AdditionalMsg: addMsg,
		},
		Result:           result,
		ResponseDatetime: tm,
	}

	return response
}

func CreateCredential(secret string, value string) (result string, err error) {

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(value))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	db, err := decodeHex([]byte(sha))
	if err != nil {
		fmt.Printf("failed to decode hex: %s", err)
		return
	}

	f := base64Encode(db)

	return string(f), err
}

func base64Encode(input []byte) []byte {
	eb := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
	base64.StdEncoding.Encode(eb, input)

	return eb
}

func decodeHex(input []byte) ([]byte, error) {
	db := make([]byte, hex.DecodedLen(len(input)))
	_, err := hex.Decode(db, input)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Base64ToHex(s string) string {
	p, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		// handle error
	}
	h := hex.EncodeToString(p)
	return h
}
