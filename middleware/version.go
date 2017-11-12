package middleware

import "github.com/labstack/echo"

//MwVersion 服务器版本
func MwVersion(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "RealClouds/4.1")
		return next(c)
	}
}
