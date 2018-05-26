package middleware

import (
	"sync"

	"github.com/labstack/echo"
	"github.com/shibingli/realclouds_go/utils"
)

//DefaultSMS *
func DefaultSMS() *SMS {

	smsAppID := utils.GetENV("SMS_APPID")
	if len(smsAppID) == 0 {
		smsAppID = "1400059472"
	}

	smsAppKey := utils.GetENV("SMS_APPKEY")
	if len(smsAppKey) == 0 {
		smsAppKey = "ayerdudu"
	}

	smsVCodeTplID := utils.GetENV("SMS_VCODE_TPLID")
	if len(smsVCodeTplID) == 0 {
		smsVCodeTplID = "11111"
	}

	sms := &SMS{
		APPID:      smsAppID,
		APIKey:     smsAppKey,
		VCodeTplID: smsVCodeTplID,
	}

	return sms
}

//SMS *
type SMS struct {
	APPID      string
	APIKey     string
	VCodeTplID string
	Mutex      sync.RWMutex
}

//MwSMS SMS middleware
func (s *SMS) MwSMS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("sms", s)
		return next(c)
	}
}
