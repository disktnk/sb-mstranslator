package mstranslator

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	primaryKeyPath          = data.MustCompilePath("primary_key")
	clientIDPath            = data.MustCompilePath("client_id")
	clientSecretPath        = data.MustCompilePath("client_secret")
	scopePath               = data.MustCompilePath("scope")
	accessTokenURLPath      = data.MustCompilePath("access_token_url")
	translatorGrantTypePath = data.MustCompilePath("grant_type")
	translatorURLPath       = data.MustCompilePath("translator_url")
)

const (
	msTranslatorScope          = "http://api.microsofttranslator.com"
	msTranslatorAccessTokenURL = "https://datamarket.accesscontrol.windows.net/v2/OAuth2-13"
	msTranslatorGrantType      = "client_credentials"
	msTranslatorURL            = "http://api.microsofttranslator.com/V2/Http.svc/Translate"
)

func NewState(ctx *core.Context, params data.Map) (core.SharedState,
	error) {
	clientID := ""
	if cid, err := params.Get(clientIDPath); err != nil {
		return nil, err
	} else if clientID, err = data.AsString(cid); err != nil {
		return nil, err
	}
	clientSecret := ""
	if cs, err := params.Get(clientSecretPath); err != nil {
		return nil, err
	} else if clientSecret, err = data.AsString(cs); err != nil {
		return nil, err
	}
	scope := msTranslatorScope
	if scp, err := params.Get(scopePath); err == nil {
		if scope, err = data.AsString(scp); err != nil {
			return nil, err
		}
	}
	grantType := msTranslatorGrantType
	if gt, err := params.Get(translatorGrantTypePath); err == nil {
		if grantType, err = data.AsString(gt); err != nil {
			return nil, err
		}
	}
	accessTokenURL := msTranslatorAccessTokenURL
	if aturl, err := params.Get(accessTokenURLPath); err == nil {
		if accessTokenURL, err = data.AsString(aturl); err != nil {
			return nil, err
		}
	}
	translatorURL := msTranslatorURL
	if tlurl, err := params.Get(translatorURLPath); err == nil {
		if translatorURL, err = data.AsString(tlurl); err != nil {
			return nil, err
		}
	}
	return &accessToken{
		clientID:       clientID,
		clientSecret:   clientSecret,
		scope:          scope,
		grantType:      grantType,
		accessTokenURL: accessTokenURL,
		translatorURL:  translatorURL,
		updateTime:     time.Now(),
	}, nil
}

type accessToken struct {
	clientID       string
	clientSecret   string
	scope          string
	grantType      string
	accessTokenURL string

	translatorURL string

	accessTokenMsg accessTokenMsg
	updateTime     time.Time
}

func (t *accessToken) Terminate(ctx *core.Context) error {
	return nil
}

type accessTokenMsg struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	Scope       string `json:"scope"`
}

func (t *accessToken) isRequiredUpdate() bool {
	duration := int(time.Since(t.updateTime)) * int(time.Second)
	expire, err := strconv.Atoi(t.accessTokenMsg.ExpiresIn)
	if err != nil {
		expire = 0
	}
	if duration >= expire {
		return true
	}
	return false
}

func (t *accessToken) updateAccessToken() error {
	value := url.Values{
		"client_id":     {t.clientID},
		"client_secret": {t.clientSecret},
		"scope":         {t.scope},
		"grant_type":    {t.grantType},
	}
	resp, err := http.PostForm(t.accessTokenURL, value)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(&t.accessTokenMsg)
}

func lookupAccessToken(ctx *core.Context, name string) (*accessToken,
	error) {
	st, err := ctx.SharedStates.Get(name)
	if err != nil {
		return nil, err
	}
	if s, ok := st.(*accessToken); ok {
		return s, nil
	}
	return nil, fmt.Errorf(
		"state '%v' cannot be converted to translate.state", name)
}

func Translate(ctx *core.Context, tokenName, from, to, target string) (
	string, error) {
	token, err := lookupAccessToken(ctx, tokenName)
	if err != nil {
		return "", err
	}
	values := url.Values{
		"text": {target},
		"from": {from},
		"to":   {to},
	}
	query := values.Encode()

	if token.isRequiredUpdate() {
		if err := token.updateAccessToken(); err != nil {
			return "", err
		}
	}
	accessToken := token.accessTokenMsg.AccessToken

	client := &http.Client{}

	urlPath := token.translatorURL + "?" + query
	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	bodyStr := "<result>" + string(body) + "</result>"
	var ret translated
	if err := xml.Unmarshal([]byte(bodyStr), &ret); err != nil {
		return "", err
	}

	return ret.Result, nil
}

type translated struct {
	Result string `xml:"string"`
}
