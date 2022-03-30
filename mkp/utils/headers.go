package mkputils

import "github.com/labstack/echo"

func GetHeader(ctx echo.Context) Header {
	return Header{
		ID:       10,
		TenantId: 10,
		Username: "admin",
	}
}
