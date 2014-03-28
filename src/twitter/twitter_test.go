package twitter

import "testing"
import . "github.com/smartystreets/goconvey/convey"

func TestNewTwitter(t *testing.T) {
	Convey("NewTwitter", t, func() {
		//設定したKey、Secretの確認
		twt := NewTwitter("key", "secret")
		Convey("Twitter Object", func() {
			So(twt, ShouldNotBeNil)
		})
	})
}

func TestGetRequestTokenUrl(t *testing.T) {
}

func TestGetAccessToken(t *testing.T) {
}

func TestSetAccessToken(t *testing.T) {
}

func TestGetTimeline(t *testing.T) {
}

func TestUpdate(t *testing.T) {
}
