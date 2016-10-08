package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

/*
 */
type OAuth struct {
	Credential *Credential

	RequestTokenUrl   string
	AuthroizeTokenUrl string
	AccessTokenUrl    string

	RequestToken *TokenSet
	AuthroizeUrl string

	authParams *parameter
}

type Credential struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

/*
 */
const (
	OAUTH_VERSION    = "1.0"
	SIGNATURE_METHOD = "HMAC-SHA1"

	CALLBACK_PARAM         = "oauth_callback"
	CONSUMER_KEY_PARAM     = "oauth_consumer_key"
	NONCE_PARAM            = "oauth_nonce"
	SESSION_HANDLE_PARAM   = "oauth_session_handle"
	SIGNATURE_METHOD_PARAM = "oauth_signature_method"
	SIGNATURE_PARAM        = "oauth_signature"
	TIMESTAMP_PARAM        = "oauth_timestamp"
	TOKEN_PARAM            = "oauth_token"
	TOKEN_SECRET_PARAM     = "oauth_token_secret"
	VERIFIER_PARAM         = "oauth_verifier"
	VERSION_PARAM          = "oauth_version"
)

/*
 */
type TokenSet struct {
	Token  string
	Secret string
}

/*
 * リクエストトークンの取得を行う
 *
 */
func (o *OAuth) GetRequestToken(callback string) error {

	wb := NewWeb()
	o.addBaseParams()
	o.addParam(CALLBACK_PARAM, callback)

	key := escape(o.Credential.ConsumerSecret) + "&" + escape("")
	base := o.requestString("GET", o.RequestTokenUrl, o.authParams)
	sign := o.sign(base, key)

	o.addParam(SIGNATURE_PARAM, sign)

	data, err := o.getBody(wb, o.RequestTokenUrl)
	if err != nil {
		return err
	}

	token := data[TOKEN_PARAM]
	secret := data[TOKEN_SECRET_PARAM]

	o.RequestToken = &TokenSet{token[0], secret[0]}

	o.AuthroizeUrl = o.AuthroizeTokenUrl + "?" + TOKEN_PARAM + "=" + escape(token[0])
	return nil
}

func (o *OAuth) requestString(method string, url string, args *parameter) string {
	ret := method + "&" + escape(url)
	esp := "&"
	for _, key := range args.keys() {
		ret += esp
		ret += escape(key + "=" + args.get(key))
		esp = escape("&")
	}
	return ret
}

func (o *OAuth) sign(message, key string) string {
	hashfun := hmac.New(sha1.New, []byte(key))
	hashfun.Write([]byte(message))
	signature := hashfun.Sum(nil)
	base64sig := make([]byte, base64.StdEncoding.EncodedLen(len(signature)))
	base64.StdEncoding.Encode(base64sig, signature)
	return string(base64sig)
}

func (o *OAuth) addBaseParams() {

	o.authParams = NewParameter()

	clock := time.Now()
	ts := clock.Unix()
	nonce := rand.New(rand.NewSource(clock.UnixNano())).Int63()

	o.addParam(VERSION_PARAM, OAUTH_VERSION)
	o.addParam(SIGNATURE_METHOD_PARAM, SIGNATURE_METHOD)
	o.addParam(TIMESTAMP_PARAM, strconv.FormatInt(ts, 10))
	o.addParam(NONCE_PARAM, strconv.FormatInt(nonce, 10))
	o.addParam(CONSUMER_KEY_PARAM, o.Credential.ConsumerKey)

	return
}

func (o *OAuth) getBody(wb *Web, accessUrl string) (map[string][]string, error) {
	wb.header.Add("Authorization", o.getOAuthHeader())

	resp, err := wb.Get(accessUrl)
	if err != nil {
		return nil, err
	}

	bodyByte, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	//レスポンスを取得してコードを取得
	body := string(bodyByte)
	parts, err := url.ParseQuery(body)
	if err != nil {
		return nil, err
	}
	return parts, nil
}

func (o *OAuth) GetAccessToken(code string) error {

	wb := NewWeb()
	o.addBaseParams()

	o.addParam(VERIFIER_PARAM, code)
	o.addParam(TOKEN_PARAM, o.RequestToken.Token)

	key := escape(o.Credential.ConsumerSecret) + "&" + escape(o.RequestToken.Secret)
	base := o.requestString("GET", o.AccessTokenUrl, o.authParams)
	sign := o.sign(base, key)

	o.addParam(SIGNATURE_PARAM, sign)

	data, err := o.getBody(wb, o.AccessTokenUrl)
	if err != nil {
		return err
	}

	token := data[TOKEN_PARAM]
	secret := data[TOKEN_SECRET_PARAM]

	o.Credential.AccessToken = token[0]
	o.Credential.AccessSecret = secret[0]

	return nil
}

func (o *OAuth) Get(url string, args map[string]string) (*http.Response, error) {
	wb := o.createOAuthWeb("GET", url, args)
	return wb.Get(url)
}

func (o *OAuth) Post(url string, args map[string]string) (*http.Response, error) {
	wb := o.createOAuthWeb("POST", url, args)
	return wb.Post(url)
}

func (o *OAuth) createOAuthWeb(method string, url string, args map[string]string) *Web {

	o.addBaseParams()
	o.addParam(TOKEN_PARAM, o.Credential.AccessToken)

	param := o.authParams.copy()
	wb := NewWeb()
	if args != nil {
		for key, value := range args {
			//Web側の引数に追加
			wb.AddParam(key, value)
			//ベースの方にも追加
			param.add(key, value)
		}
	}

	key := escape(o.Credential.ConsumerSecret) + "&" + escape(o.Credential.AccessSecret)
	base := o.requestString(method, url, param)
	sign := o.sign(base, key)

	o.addParam(SIGNATURE_PARAM, sign)

	wb.header.Add("Authorization", o.getOAuthHeader())
	return wb
}

func (o *OAuth) addParam(key, value string) {
	o.authParams.add(key, value)
}

func (o *OAuth) getOAuthHeader() string {
	hdr := "OAuth "
	for pos, key := range o.authParams.keys() {
		if pos > 0 {
			hdr += ","
		}
		hdr += key + "=\"" + o.authParams.get(key) + "\""
	}
	return hdr
}