package mkputils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	echov4 "github.com/labstack/echo/v4"
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

func DatetimeConverter(dtm string, layout string) string {

	if dtm == "" {
		return ""
	}

	dt := dtm[:len(dtm)-6]
	dt2 := dtm[len(dtm)-6:]
	year, _ := strconv.Atoi(dt[:len(dt)-4])
	month, _ := strconv.Atoi(dt[4 : len(dt)-2])
	day, _ := strconv.Atoi(dt[len(dt)-2:])

	mm2 := dt2[1 : len(dt2)-2]

	hr, _ := strconv.Atoi(dt2[:len(dt2)-4])
	mm, _ := strconv.Atoi(mm2[len(mm2)-2:])
	ss, _ := strconv.Atoi(dt2[4 : len(dt2)-0])

	var time2 = time.Date(year, time.Month(month), day, hr, mm, ss, 0, time.UTC)

	return time2.Format(layout)
}

func DateConverter(dt string, layout string) string {
	if dt == "" {
		return ""
	}

	year, _ := strconv.Atoi(dt[:len(dt)-4])
	month, _ := strconv.Atoi(dt[4 : len(dt)-2])
	day, _ := strconv.Atoi(dt[len(dt)-2:])

	var tm = time.Date(year, time.Month(month), day, 00, 00, 0, 0, time.UTC)

	return tm.Format(layout)
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

func BindValidateStructV4(ctx echov4.Context, i interface{}) error {
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

// Decrypt from base64 to decrypted string
func Aes256Decrypt(cryptoText string, saltKey ...interface{}) (interface{}, error) {
	var result interface{}
	keyText := ""
	if len(saltKey) > 0 {
		keyText = saltKey[0].(string)
	}
	key := []byte(keyText)
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return result, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return result, err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	unMarshall := json.Unmarshal(ciphertext, &result)
	fmt.Println(unMarshall)
	return result, nil
}
