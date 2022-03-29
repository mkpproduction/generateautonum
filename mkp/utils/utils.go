package mkputils

import (
	"database/sql"
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

func DatetimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func DateNow() string {
	return time.Now().Format("2006-01-02")
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