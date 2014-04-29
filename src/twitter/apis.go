package twitter

import (
	"github.com/secondarykey/golib/web"
)

/*
  Tweetの情報
*/
type TweetObject struct {
	Id                  int64
	Created_at              string
	Id_str                  string
	Text                    string
	Source                  string
	In_reply_to_user_id_str string
	User                    UserObject
	Geo                     Geo `json:"coordinates"`
	Place Places
}

type Places struct {
	Name string
	Bounding_box BoundingBox
}

type BoundingBox struct {
	Coordinates [][]Polygon `json:"coordinates,float"`
	Type        string
}

type Geo struct {
	Coordinates Point `json:"coordinates,float"`
	Type        string
}

type Polygon [4]CoordType
type Point   [2]CoordType
type CoordType float64

/*
 検索時のデータType
*/
type SearchObject struct {
	Statuses        []TweetObject
	Search_metadata SearchMetadata
}

/*
  検索APIのメタデータType
*/
type SearchMetadata struct {
	Max_id       int64
	Since_id     int64
	Refresh_url  string
	Next_results string
	Count        int64
	Completed_in float64
	Since_id_str string
	Query        string
	Max_id_str   string
}

/*
  ユーザ情報のType
*/
type UserObject struct {
	Id_str            string
	Name              string
	Screen_name       string
	Profile_image_url string
	Location string
}

/*
  タイムラインの取得
*/
func (this *Twitter) GetTimeline() ([]TweetObject, error) {

	args := this.getArgs()
	response, err := this.oauth.Get(
		"https://api.twitter.com/1.1/statuses/home_timeline.json",
		args)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var tweets []TweetObject
	err = web.ReadJson(response,tweets)
	if err != nil {
		return nil, err
	}
	return tweets, nil
}

func (this *Twitter) getArgs() map[string]string {
	args := make(map[string]string)
	for _,key := range this.param.Keys() {
		args[key] = this.param.Get(key)
	}
	return args
}

/*
  ステータス更新API呼び出し
*/
func (this *Twitter) Update(status string) error {

	this.AddParam("status",status)
	args := this.getArgs()
	resp, err := this.oauth.Post(
		"https://api.twitter.com/1.1/statuses/update.json",
		args)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var tweets TweetObject
	return web.ReadJson(resp,tweets)
}

