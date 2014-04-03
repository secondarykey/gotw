package twitter

import (
	"fmt"
	"web"
)

type Twitter struct {
	oauth *web.OAuth1
}

func Usage() {
	fmt.Println("Usage:")
}

func NewTwitter(key string, secret string) *Twitter {
	oa := web.NewOAuth1(
		key, secret,
		"https://api.twitter.com/oauth/request_token",
		"https://api.twitter.com/oauth/authorize",
		"https://api.twitter.com/oauth/access_token",
	)
	return &Twitter{
		oauth: oa,
	}
}

func (this *Twitter) SetRequestTokenAndUrl(callback string) {
	this.oauth.GetRequestToken(callback)
	return
}

func (this *Twitter) GetAuthorizationUrl() string {
	return this.oauth.AuthroizeUrl
}

func (this *Twitter) GetAccessToken(code string) {
	this.oauth.GetAccessToken(code)
	return
}

func (this *Twitter) SetAccessToken(tokenSet *web.TokenSet) {
	this.oauth.AccessToken = tokenSet
	return
}


