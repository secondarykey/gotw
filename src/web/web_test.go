package web

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestWeb(t *testing.T) {

	Convey("Webアクセスをしてみる", t, func() {
		w1 := NewWeb()
		resp, err := w1.Get("http://www.yahoo.co.jp")
		Convey("Yahoo にアクセス", func() {
			So(resp.StatusCode, ShouldEqual, 200)
			So(err, ShouldBeNil)
		})
		w2 := NewWeb()
		resp2, err2 := w2.Get("http://www.test")
		Convey("存在しないWebにアクセス", func() {
			So(err2, ShouldNotBeNil)
			So(resp2, ShouldBeNil)
		})
	})
	Convey("Params Test", t, func() {
		params := NewParams()
		params.add("ccc", "Value")
		params.add("aaa", "Value")
		params.add("BBB", "Value")
		params.add("AAA", "Value")
		p := params.Keys()
		keys := []string{"AAA", "BBB", "aaa", "ccc"}
		So(p, ShouldResemble, keys)
	})

	Convey("Escape Test", t, func() {
		So(escape("日本語"), ShouldEqual, "%E6%97%A5%E6%9C%AC%E8%AA%9E")
		So(escape("="), ShouldEqual, "%3D")
	})
}
