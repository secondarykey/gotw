package main

import (
	"fmt"
)

type Twitter struct {
	oauth *OAuth
}

func Usage() {
	fmt.Println("Usage:")
}

func NewTwitter(c *Credential) *Twitter {

	oa := &OAuth{
		Credential:        c,
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AuthroizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	}

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

func (this *Twitter) Exchange(code string) error {
	return this.oauth.GetAccessToken(code)
}
