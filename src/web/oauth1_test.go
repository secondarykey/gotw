package web

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOAuth(t *testing.T) {

	oa := NewOAuth1(
		"ysTHKkYBW9PrHtgtYyElsA",
		"Ofl3NvzYGQKeNghBZ8KP1HMcZELxfv7dVnacjpDHvQ",
		"https://api.twitter.com/oauth/request_token",
		"https://api.twitter.com/oauth/authorize",
		"https://api.twitter.com/oauth/access_token")
	oa.GetRequestToken("oob")

	Convey("OAuth Test",t ,func(){
		Convey("RequestTokenの取得",func(){
			So(oa.AuthroizeTokenUrl,ShouldNotBeNil)
			if oa.RequestToken != nil {
				So(oa.RequestToken.Token,ShouldNotEqual,"")
				So(oa.RequestToken.Secret,ShouldNotEqual,"")
			} else {
				So(oa.RequestToken,ShouldNotBeNil)
			}
		})
	})


}
