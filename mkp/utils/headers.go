package mkputils

import (
	"github.com/dgrijalva/jwt-go"
	jtwgo "github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	echov4 "github.com/labstack/echo/v4"
)

func GetHeader(ctx echo.Context) Header {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return Header{
		UID:      claims["uid"].(float64),
		TID:      claims["tid"].(float64),
		PID:      claims["pid"].(float64),
		OID:      claims["oid"].(float64),
		RID:      claims["rid"].(float64),
		Username: claims["username"].(string),
	}
}

func GetHeaderV4(ctx echov4.Context) Header {
	user := ctx.Get("user").(*jtwgo.Token)
	claims := user.Claims.(jtwgo.MapClaims)

	return Header{
		UID:      claims["uid"].(float64),
		TID:      claims["tid"].(float64),
		PID:      claims["pid"].(float64),
		OID:      claims["oid"].(float64),
		RID:      claims["rid"].(float64),
		Username: claims["username"].(string),
	}
}
