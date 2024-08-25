package oauth

import (
	"fmt"
	"net"
	"net/http"
	"pandax/kit/model"
	"time"

	"golang.org/x/oauth2"
)

const (
	GithubAuthType = "github"
	WeChatAuthType = "wechat"
	EmptyAuthType  = "nil"
)

var (
	defaultHttpClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
)

func IsEmptyAuthType(authType string) bool {
	return authType == "" || authType == EmptyAuthType
}

type OAuthConfig struct {
	AuthType     string `yaml:"authType"`
	ClientId     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
}

type UserInfo struct {
	ID          string
	Url         string
	AuthType    string
	Username    string
	DisplayName string
	Email       string
	AvatarUrl   string
}

func (ui *UserInfo) User() *model.OauthAccount {
	return &model.OauthAccount{
		Name:   ui.Username,
		Email:  ui.Email,
		Avatar: ui.AvatarUrl,
		AuthInfos: []model.AuthInfo{
			{
				AuthType: ui.AuthType,
				AuthId:   ui.ID,
				Url:      ui.Url,
			},
		},
	}
}

type OAuthManager struct {
	conf map[string]OAuthConfig
}

func NewOAuthManager(conf map[string]OAuthConfig) *OAuthManager {
	return &OAuthManager{
		conf: conf,
	}
}

func (m *OAuthManager) GetAuthProvider(authType string) (AuthProvider, error) {
	var provider AuthProvider
	conf, ok := m.conf[authType]
	if !ok {
		return nil, fmt.Errorf("auth type %s not found in config", authType)
	}
	switch authType {
	case GithubAuthType:
		provider = NewGithubAuth(conf.ClientId, conf.ClientSecret)
	case WeChatAuthType:
		provider = NewWeChatAuth(conf.ClientId, conf.ClientSecret)
	default:
		return nil, fmt.Errorf("unknown auth type: %s", authType)
	}

	return provider, nil
}

type AuthProvider interface {
	GetToken(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*UserInfo, error)
}
