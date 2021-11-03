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

//WeChatErr *
type WeChatErr struct {
	ErrCode int                    `json:"errcode,omitempty" xml:"errcode,omitempty"`
	ErrMsg  string                 `json:"errmsg,omitempty" xml:"errmsg,omitempty"`
	Hints   map[string]interface{} `json:"hints,omitempty" xml:"hints,omitempty"`
}

//WeChatUserInfo *
type WeChatUserInfo struct {
	WeChatErr
	OpenID     string   `json:"openid,omitempty"`
	Nickname   string   `json:"nickname,omitempty"`
	Sex        int      `json:"sex,omitempty"`
	Province   string   `json:"province,omitempty"`
	City       string   `json:"city,omitempty"`
	Country    string   `json:"country,omitempty"`
	Headimgurl string   `json:"headimgurl,omitempty"`
	Privilege  []string `json:"privilege,omitempty"`
	UnionID    string   `json:"unionid,omitempty"`
}

//WeChatAuthCodeURL *
func (c *Config) WeChatAuthCodeURL(state string, qr bool, opts ...AuthCodeOption) string {
	var buf bytes.Buffer

	authURL := ""

	if qr {
		authURL = c.Endpoint.QRAuthURL
	} else {
		authURL = c.Endpoint.AuthURL
	}
	buf.WriteString(authURL)
	v := url.Values{
		"response_type": {"code"},
		"appid":         {c.ClientID},
		"redirect_uri":  CondVal(c.RedirectURL),
		"scope":         CondVal(strings.Join(c.Scopes, " ")),
		"state":         CondVal(state),
	}

	for _, opt := range opts {
		opt.setValue(v)
	}

	if strings.Contains(authURL, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}

	buf.WriteString(v.Encode())
	// buf.WriteString("&_=" + fmt.Sprintf("%v", utils.RandInt64()))
	buf.WriteString("#wechat_redirect")
	return buf.String()
}

//WeChatExchange *
func (c *Config) WeChatExchange(ctx context.Context, code string, opts ...AuthCodeOption) (*Token, error) {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.TokenURL)
	v := url.Values{
		"grant_type": {"authorization_code"},
		"appid":      {c.ClientID},
		"secret":     {c.ClientSecret},
		"code":       {code},
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

	var token Token
	if err := json.Unmarshal(body, &token); nil != err {
		log.Printf("Token body err:%v\n", err.Error())
		return nil, err
	}

	t := &Token{
		WeChatErr: WeChatErr{
			ErrCode: token.ErrCode,
			ErrMsg:  token.ErrMsg,
		},
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		ExpiresIN:    token.ExpiresIN,
		OpenID:       token.OpenID,
		UnionID:      token.UnionID,
		SCOPE:        token.SCOPE,
	}

	return t, nil
}

//MPAccessToken *
func (c *Config) MPAccessToken(ctx context.Context, opts ...AuthCodeOption) (*Token, error) {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.MPTokenURL)
	v := url.Values{
		"grant_type": {"client_credential"},
		"appid":      {c.ClientID},
		"secret":     {c.ClientSecret},
	}

	for _, opt := range opts {
		opt.setValue(v)
	}

	if strings.Contains(c.Endpoint.MPTokenURL, "?") {
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

	var token Token
	if err := json.Unmarshal(body, &token); nil != err {
		log.Printf("Token body err:%v\n", err.Error())
		return nil, err
	}

	t := &Token{
		WeChatErr: WeChatErr{
			ErrCode: token.ErrCode,
			ErrMsg:  token.ErrMsg,
		},
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		ExpiresIN:    token.ExpiresIN,
		OpenID:       token.OpenID,
		UnionID:      token.UnionID,
		SCOPE:        token.SCOPE,
	}

	return t, nil
}

//MPJSAPITicket *
func (c *Config) MPJSAPITicket(ctx context.Context, accessToken string, opts ...AuthCodeOption) (*Token, error) {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.TicketURL)

	v := url.Values{
		"access_token": {accessToken},
		"type":         {"jsapi"},
	}

	for _, opt := range opts {
		opt.setValue(v)
	}

	if strings.Contains(c.Endpoint.TicketURL, "?") {
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

	var token Token
	if err := json.Unmarshal(body, &token); nil != err {
		log.Printf("Token body err:%v\n", err.Error())
		return nil, err
	}

	t := &Token{
		WeChatErr: WeChatErr{
			ErrCode: token.ErrCode,
			ErrMsg:  token.ErrMsg,
		},
		Ticket:    token.Ticket,
		ExpiresIN: token.ExpiresIN,
	}

	return t, nil
}

//WeChatUserInfo *
func (c *Config) WeChatUserInfo(ctx context.Context, token *Token, opts ...AuthCodeOption) (*WeChatUserInfo, error) {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.UserInfoURL)
	v := url.Values{
		"access_token": {token.AccessToken},
		"openid":       {token.OpenID},
		"lang":         {"zh-CN"},
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

	// if _, err := WeChatErrChk(body); nil != err {
	// 	return nil, err
	// }

	var userInfo WeChatUserInfo
	if err := json.Unmarshal(body, &userInfo); nil != err {
		return nil, err
	}

	u := &WeChatUserInfo{
		WeChatErr: WeChatErr{
			ErrCode: userInfo.ErrCode,
			ErrMsg:  userInfo.ErrMsg,
		},
		OpenID:     userInfo.OpenID,
		Nickname:   userInfo.Nickname,
		Sex:        userInfo.Sex,
		Province:   userInfo.Province,
		City:       userInfo.City,
		Country:    userInfo.Country,
		Headimgurl: userInfo.Headimgurl,
		Privilege:  userInfo.Privilege,
		UnionID:    userInfo.UnionID,
	}

	return u, nil
}
