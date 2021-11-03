package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/GreatSir/realclouds_go/utils"
)

//AuthCodeOption *
type AuthCodeOption interface {
	setValue(url.Values)
}

type setParam struct{ k, v string }

func (p setParam) setValue(m url.Values) { m.Set(p.k, p.v) }

//SetAuthURLParam *
func SetAuthURLParam(key, value string) AuthCodeOption {
	return setParam{key, value}
}

//Endpoint *
type Endpoint struct {
	TicketURL   string `json:"ticket_url,omitempty" xml:"ticket_url,omitempty"`
	QRAuthURL   string `json:"qr_auth_url,omitempty" xml:"qr_auth_url,omitempty"`
	AuthURL     string `json:"auth_url,omitempty" xml:"auth_url,omitempty"`
	TokenURL    string `json:"token_url,omitempty" xml:"token_url,omitempty"`
	OpenIDURL   string `json:"open_id_url,omitempty" xml:"open_id_url,omitempty"`
	MPTokenURL  string `json:"mp_token_url,omitempty" xml:"mp_token_url,omitempty"`
	UserInfoURL string `json:"user_info_url,omitempty" xml:"user_info_url,omitempty"`
}

//Config *
type Config struct {
	ClientID       string   `json:"client_id" xml:"client_id"`
	ClientSecret   string   `json:"client_secret" xml:"client_secret"`
	Token          string   `json:"token" xml:"token"`
	EncodingAESKey string   `json:"encoding_aes_key" xml:"encoding_aes_key"`
	Endpoint       Endpoint `json:"endpoint" xml:"endpoint"`
	RedirectURL    string   `json:"redirect_url" xml:"redirect_url"`
	Scopes         []string `json:"scopes" xml:"scopes"`
}

//Token *
type Token struct {
	WeChatErr
	Ticket       string    `json:"ticket,omitempty" xml:"ticket,omitempty"`
	AccessToken  string    `json:"access_token" xml:"access_token"`
	TokenType    string    `json:"token_type,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresIN    int       `json:"expires_in,omitempty"`
	RemindIN     string    `json:"remind_in,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
	UID          string    `json:"uid,omitempty"`
	OpenID       string    `json:"openid,omitempty"`
	UnionID      string    `json:"unionid,omitempty"`
	ClientID     string    `json:"client_id,omitempty"`
	SCOPE        string    `json:"scope,omitempty"`
}

//CondVal *
func CondVal(v string) []string {
	if v == "" {
		return nil
	}
	return []string{v}
}

//AuthCodeURL *
func (c *Config) AuthCodeURL(state string, opts ...AuthCodeOption) string {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.AuthURL)
	v := url.Values{
		"response_type": {"code"},
		"client_id":     {c.ClientID},
		"redirect_uri":  CondVal(c.RedirectURL),
		"scope":         CondVal(strings.Join(c.Scopes, " ")),
		"state":         CondVal(state),
	}

	for _, opt := range opts {
		opt.setValue(v)
	}

	if strings.Contains(c.Endpoint.AuthURL, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(v.Encode())
	buf.WriteString("&_=" + fmt.Sprintf("%v", utils.RandInt64()))
	return buf.String()
}

//Exchange *
func (c *Config) Exchange(ctx context.Context, code string, opts ...AuthCodeOption) (*Token, error) {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.TokenURL)
	v := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": CondVal(c.RedirectURL),
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
