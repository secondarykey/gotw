package oauth

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"math"
)

func TestWeb(t *testing.T) {
	Convey("Webアクセスをしてみる", t, func() {
		web := NewWeb()
		resp, err := web.Get("http://www.yahoo.co.jp")
		Convey("Yahoo にアクセス", func() {
			So(resp.StatusCode, ShouldEqual, 200)
			So(err, ShouldBeNil)
		})

		web2 := NewWeb()
		resp2, err2 := web2.Get("http://www.test")
		Convey("存在しないWebにアクセス", func() {
			So(resp2.StatusCode, ShouldEqual, 404)
			So(err2, ShouldNotBeNil)
		})
	})
}
