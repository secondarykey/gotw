package twitter

import (
	"fmt"
	"github.com/mrjones/oauth"
)

type Twitter struct {
	consumer *oauth.Consumer
	Url      string
	rtoken   *oauth.RequestToken
	atoken   *oauth.AccessToken
}

func Usage() {
	fmt.Println("Usage:")
}

func NewTwitter(key string, secret string) *Twitter {
	c := oauth.NewConsumer(key, secret,
		oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	})
	return &Twitter{
		consumer: c,
	}

}

func (this *Twitter) GetRequestTokenAndUrl(callback string) *oauth.RequestToken {
	requestToken, url, err := this.consumer.GetRequestTokenAndUrl(callback)
	if err != nil {
		panic(err)
	}
	this.rtoken = requestToken
	this.Url = url
	return requestToken
}

func (this *Twitter) GetAccessToken(code string) *oauth.AccessToken {
	accessToken, err := this.consumer.AuthorizeToken(this.rtoken, code)
	if err != nil {
		panic(err)
	}
	this.atoken = accessToken
	return accessToken
}

func (this *Twitter) SetAccessToken(accessToken *oauth.AccessToken) {
	this.atoken = accessToken
}
