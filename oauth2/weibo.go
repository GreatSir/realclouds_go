package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/shibingli/realclouds_go/utils"
)

//WeboErr *
type WeboErr struct {
	ErrCode int                    `json:"errcode,omitempty" xml:"errcode,omitempty"`
	ErrMsg  string                 `json:"errmsg,omitempty" xml:"errmsg,omitempty"`
	Hints   map[string]interface{} `json:"hints,omitempty" xml:"hints,omitempty"`
}

//WeboUserInfo *
type WeboUserInfo struct {
	WeboErr
	ID               int64               `json:"id,omitempty"`
	IDStr            string              `json:"idstr,omitempty"`
	Class            int                 `json:"class,omitempty"`
	ScreenName       string              `json:"screen_name,omitempty"`
	Name             string              `json:"name,omitempty"`
	Province         string              `json:"province,omitempty"`
	City             string              `json:"city,omitempty"`
	Location         string              `json:"location,omitempty"`
	Description      string              `json:"description,omitempty"`
	URL              string              `json:"url,omitempty"`
	ProfileImageURL  string              `json:"profile_image_url,omitempty"`
	CoverImagePhone  string              `json:"cover_image_phone,omitempty"`
	ProfileURL       string              `json:"profile_url,omitempty"`
	Domain           string              `json:"domain,omitempty"`
	Weihao           string              `json:"weihao,omitempty"`
	Gender           string              `json:"gender,omitempty"`
	FollowersCount   int64               `json:"followers_count,omitempty"`
	FriendsCount     int64               `json:"friends_count,omitempty"`
	PagefriendsCount int64               `json:"pagefriends_count,omitempty"`
	StatusesCount    int64               `json:"statuses_count,omitempty"`
	FavouritesCount  int64               `json:"favourites_count,omitempty"`
	CreatedAT        string              `json:"created_at,omitempty"`
	Following        bool                `json:"following,omitempty"`
	AllowAllActMsg   bool                `json:"allow_all_act_msg,omitempty"`
	GeoEnabled       bool                `json:"geo_enabled,omitempty"`
	Verified         bool                `json:"verified,omitempty"`
	VerifiedType     int                 `json:"verified_type,omitempty"`
	Remark           string              `json:"remark,omitempty"`
	Insecurity       WeiboUserInsecurity `json:"insecurity,omitempty"`
	Status           WeiboUserStatus     `json:"status,omitempty"`
	AllowAllComment  bool                `json:"allow_all_comment,omitempty"`
	AvatarLarge      string              `json:"avatar_large,omitempty"`
	VerifiedReason   string              `json:"verified_reason,omitempty"`
	FollowMe         bool                `json:"follow_me,omitempty"`
	OnlineStatus     int                 `json:"online_status,omitempty"`
	BiFollowersCount int                 `json:"bi_followers_count,omitempty"`
}

//WeiboUserInsecurity *
type WeiboUserInsecurity struct {
	SexualContent bool `json:"sexual_content,omitempty"`
}

//WeiboUserStatus *
type WeiboUserStatus struct {
	ID                  int64                            `json:"id,omitempty"`
	IDStr               string                           `json:"idstr,omitempty"`
	CreatedAT           string                           `json:"created_at,omitempty"`
	Text                string                           `json:"text,omitempty"`
	TextLength          int                              `json:"textLength,omitempty"`
	Source              string                           `json:"source,omitempty"`
	SourceAllowClick    int                              `json:"source_allowclick,omitempty"`
	SourceType          int                              `json:"source_type,omitempty"`
	Favorited           bool                             `json:"favorited,omitempty"`
	Truncated           bool                             `json:"truncated,omitempty"`
	InReplyToStatusID   string                           `json:"in_reply_to_status_id,omitempty"`
	InReplyToUserID     string                           `json:"in_reply_to_user_id,omitempty"`
	InReplyToScreenName string                           `json:"in_reply_to_screen_name,omitempty"`
	Geo                 string                           `json:"geo,omitempty"`
	Mid                 string                           `json:"mid,omitempty"`
	Annotations         []string                         `json:"annotations,omitempty"`
	RepostsCount        int                              `json:"reposts_count,omitempty"`
	CommentsCount       int                              `json:"comments_count,omitempty"`
	PicURLs             []string                         `json:"pic_urls,omitempty"`
	IsPaid              bool                             `json:"is_paid,omitempty"`
	MblogVIPType        int                              `json:"mblog_vip_type,omitempty"`
	AttitudesCount      int64                            `json:"attitudes_count,omitempty"`
	IsLongText          bool                             `json:"isLongText,omitempty"`
	Mlevel              int                              `json:"mlevel,omitempty"`
	Visible             WeiboUserStatusVisible           `json:"visible,omitempty"`
	BizFeature          int                              `json:"biz_feature,omitempty"`
	HasActionTypeCard   int                              `json:"hasActionTypeCard,omitempty"`
	DarwinTags          []string                         `json:"darwin_tags,omitempty"`
	HotWeiboTags        []string                         `json:"hot_weibo_tags,omitempty"`
	TextTagTips         []string                         `json:"text_tag_tips,omitempty"`
	UserType            int                              `json:"userType,omitempty"`
	MoreInfoType        int                              `json:"more_info_type,omitempty"`
	PositiveRecomFlag   int                              `json:"positive_recom_flag,omitempty"`
	GifIDs              string                           `json:"gif_ids,omitempty"`
	IsShowBulletin      int                              `json:"is_show_bulletin,omitempty"`
	CommentManageInfo   WeiboUserStatusCommentManageInfo `json:"comment_manage_info,omitempty"`
}

//WeiboUserStatusVisible *
type WeiboUserStatusVisible struct {
	Type   int `json:"type,omitempty"`
	ListID int `json:"list_id,omitempty"`
}

//WeiboUserStatusCommentManageInfo *
type WeiboUserStatusCommentManageInfo struct {
	CommentPermissionType int `json:"comment_permission_type,omitempty"`
}

//WeboExchange *
func (c *Config) WeboExchange(ctx context.Context, code string, opts ...AuthCodeOption) (*Token, error) {
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

	body, err := httpLib.POST(buf.String(), nil)
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
		AccessToken: token.AccessToken,
		ExpiresIN:   token.ExpiresIN,
		RemindIN:    token.RemindIN,
		UID:         token.UID,
	}
	return t, nil
}

//WeboUserInfo *
func (c *Config) WeboUserInfo(ctx context.Context, token *Token, opts ...AuthCodeOption) (*WeboUserInfo, error) {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.UserInfoURL)
	v := url.Values{
		"access_token": {token.AccessToken},
		"uid":          {token.UID},
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

	var userInfo WeboUserInfo
	if err := json.Unmarshal(body, &userInfo); nil != err {
		return nil, err
	}

	u := &WeboUserInfo{
		ID:               userInfo.ID,
		IDStr:            userInfo.IDStr,
		Class:            userInfo.Class,
		ScreenName:       userInfo.ScreenName,
		Name:             userInfo.Name,
		Province:         userInfo.Province,
		City:             userInfo.City,
		Location:         userInfo.Location,
		Description:      userInfo.Description,
		URL:              userInfo.URL,
		ProfileImageURL:  userInfo.ProfileImageURL,
		CoverImagePhone:  userInfo.CoverImagePhone,
		ProfileURL:       userInfo.ProfileURL,
		Domain:           userInfo.Domain,
		Weihao:           userInfo.Weihao,
		Gender:           userInfo.Gender,
		FollowersCount:   userInfo.FollowersCount,
		FriendsCount:     userInfo.FriendsCount,
		PagefriendsCount: userInfo.PagefriendsCount,
		StatusesCount:    userInfo.StatusesCount,
		FavouritesCount:  userInfo.FavouritesCount,
		CreatedAT:        userInfo.CreatedAT,
		Following:        userInfo.Following,
		AllowAllActMsg:   userInfo.AllowAllActMsg,
		GeoEnabled:       userInfo.GeoEnabled,
		Verified:         userInfo.Verified,
		VerifiedType:     userInfo.VerifiedType,
		Remark:           userInfo.Remark,
		Insecurity:       userInfo.Insecurity,
		Status:           userInfo.Status,
		AllowAllComment:  userInfo.AllowAllComment,
		AvatarLarge:      userInfo.AvatarLarge,
		VerifiedReason:   userInfo.VerifiedReason,
		FollowMe:         userInfo.FollowMe,
		OnlineStatus:     userInfo.OnlineStatus,
		BiFollowersCount: userInfo.BiFollowersCount,
	}

	return u, nil
}
