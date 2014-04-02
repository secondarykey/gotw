package web

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
type OAuth1 struct {
	ConsumerKey    string
	ConsumerSecret string

	RequestTokenUrl   string
	AuthroizeTokenUrl string
	AccessTokenUrl    string

	RequestToken *TokenSet
	AuthroizeUrl string

	AccessToken *TokenSet

	authParams *parameter
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
 */
func NewOAuth1(key, secret, requestTokenUrl, authroizeTokenUrl, accessTokenUrl string) *OAuth1 {
	return &OAuth1{
		ConsumerKey:       key,
		ConsumerSecret:    secret,
		RequestTokenUrl:   requestTokenUrl,
		AuthroizeTokenUrl: authroizeTokenUrl,
		AccessTokenUrl:    accessTokenUrl,
	}
}

/*
 * リクエストトークンの取得を行う
 *
 */
func (self *OAuth1) GetRequestToken(callback string) {

	wb := NewWeb()
	self.addBaseParams()
	self.addParam(CALLBACK_PARAM, callback)

	key := escape(self.ConsumerSecret) + "&" + escape("")
	base := self.requestString("GET", self.RequestTokenUrl, self.authParams)
	sign := self.sign(base, key)

	self.addParam(SIGNATURE_PARAM, sign)

	data, err := self.getBody(wb, self.RequestTokenUrl)
	if err != nil {
		panic(err)
	}

	token := data[TOKEN_PARAM]
	secret := data[TOKEN_SECRET_PARAM]

	self.RequestToken = &TokenSet{token[0], secret[0]}

	self.AuthroizeUrl = self.AuthroizeTokenUrl+"?"+TOKEN_PARAM+"="+escape(token[0])
	return
}

func (self *OAuth1) requestString(method string, url string, args *parameter) string {
	ret := method + "&" + escape(url)
	esp := "&"
	for _, key := range args.Keys() {
		ret += esp
		ret += escape(key+"="+args.Get(key))
		esp = escape("&")
	}
	return ret
}

func (self *OAuth1) sign(message, key string) string {
	hashfun := hmac.New(sha1.New, []byte(key))
	hashfun.Write([]byte(message))
	signature := hashfun.Sum(nil)
	base64sig := make([]byte, base64.StdEncoding.EncodedLen(len(signature)))
	base64.StdEncoding.Encode(base64sig, signature)
	return string(base64sig)
}

func (self *OAuth1) addBaseParams() {

	self.authParams = NewParams()

	clock := time.Now()
	ts := clock.Unix()
	nonce := rand.New(rand.NewSource(clock.UnixNano())).Int63()

	self.addParam(VERSION_PARAM, OAUTH_VERSION)
	self.addParam(SIGNATURE_METHOD_PARAM, SIGNATURE_METHOD)
	self.addParam(TIMESTAMP_PARAM, strconv.FormatInt(ts, 10))
	self.addParam(NONCE_PARAM, strconv.FormatInt(nonce, 10))
	self.addParam(CONSUMER_KEY_PARAM, self.ConsumerKey)

	return
}

func (self *OAuth1) getBody(wb *Web, accessUrl string) (map[string][]string, error) {
	wb.header.Add("Authorization", self.getOAuthHeader())

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

func (self *OAuth1) GetAccessToken(code string) {

	wb := NewWeb()
	self.addBaseParams()

	self.addParam(VERIFIER_PARAM, code)
	self.addParam(TOKEN_PARAM, self.RequestToken.Token)

	key := escape(self.ConsumerSecret) + "&" + escape(self.RequestToken.Secret)
	base := self.requestString("GET", self.AccessTokenUrl, self.authParams)
	sign := self.sign(base, key)

	self.addParam(SIGNATURE_PARAM, sign)

	data, err := self.getBody(wb, self.AccessTokenUrl)
	if err != nil {
		panic(err)
	}

	token := data[TOKEN_PARAM]
	secret := data[TOKEN_SECRET_PARAM]

	self.AccessToken = &TokenSet{token[0], secret[0]}

	return
}

func (self *OAuth1) Get(url string, args map[string]string) (*http.Response, error) {

	wb := self.createOAuthWeb("GET", url, args)

	return wb.Get(url)
}

func (self *OAuth1) Post(url string, args map[string]string) (*http.Response, error) {
	wb := self.createOAuthWeb("POST", url, args)
	return wb.Post(url)
}

func (self *OAuth1) createOAuthWeb(method string, url string, args map[string]string) *Web {

	self.addBaseParams()
	self.addParam(TOKEN_PARAM, self.AccessToken.Token)

	param := self.authParams.Copy()

	wb := NewWeb()
	if args != nil {
		for key, value := range args {
			wb.AddParam(key, value)
			param.add(key, value)
		}
	}

	key := escape(self.ConsumerSecret) + "&" + escape(self.AccessToken.Secret)
	base := self.requestString(method, url, param)
	sign := self.sign(base, key)

	self.addParam(SIGNATURE_PARAM, sign)

	wb.header.Add("Authorization", self.getOAuthHeader())
	return wb
}

func (self *OAuth1) addParam(key, value string) {
	self.authParams.add(key, value)
}

func (self *OAuth1) getOAuthHeader() string {
	hdr := "OAuth "
	for pos, key := range self.authParams.Keys() {
		if pos > 0 {
			hdr += ","
		}
		hdr += key+"=\""+self.authParams.Get(key)+"\""
	}
	return hdr
}
