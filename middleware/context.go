package middleware

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	session "github.com/ipfans/echo-session"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/shibingli/realclouds_go/utils"
)

//MwContext Context middleware
func MwContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		mc := &Context{c}
		return next(mc)
	}
}

//Context 自定义 Context
type Context struct {
	echo.Context
}

//GetSession 获取Session
func (c *Context) GetSession(key string) interface{} {
	session := session.Default(c.Context)
	val := session.Get(strings.TrimSpace(key))
	return val
}

//SetSession 添加Session
func (c *Context) SetSession(key string, val interface{}) error {
	session := session.Default(c.Context)
	session.Set(strings.TrimSpace(key), val)
	return session.Save()
}

//RemoveSession 删除Session
func (c *Context) RemoveSession(key string) error {
	session := session.Default(c.Context)
	session.Delete(strings.TrimSpace(key))
	return session.Save()
}

//IsAjax 判断请求是否为Ajax请求
func (c *Context) IsAjax() bool {
	if "" != c.Request().Header.Get("X-Requested-With") {
		return true
	}
	return false
}

//FormValue 获取表单参数
func (c *Context) FormValue(key string) string {
	val := c.Context.FormValue(strings.TrimSpace(key))
	return strings.TrimSpace(val)
}

//PathValue 获取路径参数
func (c *Context) PathValue(key string) string {
	val := c.Context.Param(strings.TrimSpace(key))
	return strings.TrimSpace(val)
}

//ToHTML 根据模板名称输出HTML
func (c *Context) ToHTML(tpl string, data interface{}) error {
	resultMap := make(map[string]interface{})
	resultMap["Data"] = data
	return c.Render(http.StatusOK, tpl, resultMap)
}

//ToJSON 输出 JSON
func (c *Context) ToJSON(data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

//ToXML 输出 XML
func (c *Context) ToXML(data interface{}) error {
	return c.XML(http.StatusOK, data)
}

//ToString 输出 String
func (c *Context) ToString(val string) error {
	return c.String(http.StatusOK, val)
}

//JSONBind 绑定JSON
func (c *Context) JSONBind(val interface{}) error {
	body := c.Request().Body
	defer body.Close()

	byteVal, err := ioutil.ReadAll(body)
	if nil != err {
		return err
	}

	if err := json.Unmarshal(byteVal, val); nil != err {
		return err
	}

	return nil
}

//PermanentRedirect 永久跳转 HttpStatusCode 308
func (c *Context) PermanentRedirect(path string) error {
	path = utils.StringUtils(path).RandURL()
	return c.Redirect(http.StatusPermanentRedirect, path)
}

//TemporaryRedirect 临时跳转 HttpStatusCode 307
func (c *Context) TemporaryRedirect(path string) error {
	path = utils.StringUtils(path).RandURL()
	return c.Redirect(http.StatusTemporaryRedirect, path)
}

//MySQL 获取 MySQL driver
func (c *Context) MySQL() *gorm.DB {
	return c.Get("mysql").(*gorm.DB)
}

//Redis 获取 Redis pool
func (c *Context) Redis() *Redis {
	return c.Get("redis").(*Redis)
}

//SMS 获取 SMS info
func (c *Context) SMS() *SMS {
	return c.Get("sms").(*SMS)
}

//DrityWord 获取 Drity word
func (c *Context) DrityWord() *DrityWord {
	return c.Get("drityword").(*DrityWord)
}

//NewCtx 获取 WebContext
func NewCtx(c echo.Context) *Context {
	return c.(*Context)
}
