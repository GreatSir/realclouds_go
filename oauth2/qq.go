package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Greatsir/realclouds_go/utils"
)

//QQErr *
type QQErr struct {
	Ret int    `json:"ret,omitempty" xml:"ret,omitempty"`
	Msg string `json:"msg,omitempty" xml:"msg,omitempty"`
}

//QQUserInfo *
type QQUserInfo struct {
	QQErr
	IsLost          int    `json:"is_lost,omitempty"`
	Nickname        string `json:"nickname,omitempty"`
	Gender          string `json:"gender,omitempty"`
	Province        string `json:"province,omitempty"`
	City            string `json:"city,omitempty"`
	Year            string `json:"year,omitempty"`
	FigureURL       string `json:"figureurl,omitempty"`
	FigureURL1      string `json:"figureurl_1,omitempty"`
	FigureURL2      string `json:"figureurl_2,omitempty"`
	FigureURLQQ1    string `json:"figureurl_qq_1,omitempty"`
	FigureURLQQ2    string `json:"figureurl_qq_2,omitempty"`
	IsYellowVIP     string `json:"is_yellow_vip,omitempty"`
	VIP             string `json:"vip,omitempty"`
	YellowVIPLevel  string `json:"yellow_vip_level,omitempty"`
	Level           string `json:"level,omitempty"`
	IsYellowYearVIP string `json:"is_yellow_year_vip,omitempty"`
}

//QQExchange *
func (c *Config) QQExchange(ctx context.Context, code string, opts ...AuthCodeOption) (*Token, error) {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.TokenURL)
	v := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"code":          {code},
		"redirect_uri":  {c.RedirectURL},
	}

	for _, opt := range opts {
		opt.setValue(v)
	}

	if strings.Contains(c.Endpoint.TokenURL, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	uri := buf.String()
	httpLib := utils.NewHTTPLib(uri)

	buf.Reset()

	path := v.Encode()
	buf.WriteString(path)
	buf.WriteString("&_=" + fmt.Sprintf("%v", utils.RandInt64()))

	body, err := httpLib.GET(buf.String(), nil)
	if nil != err {
		return nil, err
	}

	log.Printf("Token:%v\n", string(body))

	values, err := url.ParseQuery(string(body))
	if nil != err {
		return nil, err
	}

	accessToken := values.Get("access_token")
	expiresIN := values.Get("expires_in")
	refreshToken := values.Get("refresh_token")

	ei, err := utils.StringUtils(expiresIN).Int()
	if nil != err {
		return nil, err
	}

	retMap, err := c.QQOpenID(context.Background(), accessToken)
	if nil != err {
		return nil, err
	}

	openID, openIDOK := retMap["openid"]
	clientID, clientIDOK := retMap["client_id"]
	if !openIDOK || !clientIDOK {
		return nil, fmt.Errorf("%v", "Invalid QQ OpenID")
	}

	t := &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		OpenID:       openID,
		ClientID:     clientID,
		ExpiresIN:    ei,
	}

	return t, nil
}

//QQOpenID *
func (c *Config) QQOpenID(ctx context.Context, accessToken string) (map[string]string, error) {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.OpenIDURL)
	v := url.Values{
		"access_token": {accessToken},
	}

	if strings.Contains(c.Endpoint.OpenIDURL, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	uri := buf.String()
	httpLib := utils.NewHTTPLib(uri)

	buf.Reset()

	path := v.Encode()
	buf.WriteString(path)
	buf.WriteString("&_=" + fmt.Sprintf("%v", utils.RandInt64()))

	body, err := httpLib.GET(buf.String(), nil)
	if nil != err {
		return nil, err
	}

	sri := bytes.Index(body, []byte("{"))
	eri := bytes.Index(body, []byte("}"))

	if -1 == sri || -1 == eri {
		return nil, fmt.Errorf("%v", "Token body err:")
	}

	var retMap map[string]string
	if err := json.Unmarshal(body[sri:eri+1], &retMap); nil != err {
		log.Printf("Token body err:%v\n", err.Error())
		return nil, err
	}

	return retMap, nil
}

//QQUserInfo *
func (c *Config) QQUserInfo(ctx context.Context, token *Token, opts ...AuthCodeOption) (*QQUserInfo, error) {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.UserInfoURL)
	v := url.Values{
		"access_token": {token.AccessToken},
		"openid":       {token.OpenID},
	}

	for _, opt := range opts {
		opt.setValue(v)
	}

	if strings.Contains(c.Endpoint.UserInfoURL, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}

	httpLib := utils.NewHTTPLib(buf.String())

	buf.Reset()

	buf.WriteString(v.Encode())
	buf.WriteString("&_=" + fmt.Sprintf("%v", utils.RandInt64()))

	body, err := httpLib.GET(buf.String(), nil)
	if nil != err {
		return nil, err
	}

	log.Printf("QQ Info: %v\n", string(body))

	var userInfo QQUserInfo
	if err := json.Unmarshal(body, &userInfo); nil != err {
		return nil, err
	}

	u := &QQUserInfo{
		IsLost:          userInfo.IsLost,
		Nickname:        userInfo.Nickname,
		Gender:          userInfo.Gender,
		Province:        userInfo.Province,
		City:            userInfo.City,
		Year:            userInfo.Year,
		FigureURL:       userInfo.FigureURL,
		FigureURL1:      userInfo.FigureURL1,
		FigureURL2:      userInfo.FigureURL2,
		FigureURLQQ1:    userInfo.FigureURLQQ1,
		FigureURLQQ2:    userInfo.FigureURLQQ2,
		IsYellowVIP:     userInfo.IsYellowVIP,
		VIP:             userInfo.VIP,
		YellowVIPLevel:  userInfo.YellowVIPLevel,
		Level:           userInfo.Level,
		IsYellowYearVIP: userInfo.IsYellowYearVIP,
	}

	return u, nil
}
