package web

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"encoding/json"
)

func TestOAuth(t *testing.T) {

	var tokenSet TokenSet
	readJson(&tokenSet, "consumer.json")
	oa := NewOAuth1(
		tokenSet.Token,
		tokenSet.Secret,
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

/*
 *
 *
 *
 *
 */
func readJson(token interface{}, filename string) error {
	if b, err := ioutil.ReadFile(filename); err != nil {
		return err
	} else {
		return json.Unmarshal(b, token)
	}
}
