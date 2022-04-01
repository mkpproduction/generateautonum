package mkputils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func GetHeader(ctx echo.Context) Header {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return Header{
		UID:      claims["id"].(float64),
		TID:      claims["tid"].(float64),
		PID:      claims["pid"].(float64),
		OID:      claims["oid"].(float64),
		RID:      claims["rid"].(float64),
		Username: claims["username"].(string),
	}
}
