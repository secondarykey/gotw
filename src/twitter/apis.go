package twitter

import (
	"encoding/base64"
	"encoding/json"
	"github.com/secondarykey/golib/web"
	"io/ioutil"
)

/*
  Tweetの情報
*/
type TweetObject struct {
	Created_at              string
	Id_str                  string
	Text                    string
	Source                  string
	In_reply_to_user_id_str string
	User                    UserObject
	Geo                     Geo `json:"coordinates"`
}

type Geo struct {
	Coordinates Coordinate `json:"coordinates,float"`
	Type        string
}

type Coordinate [2]CoordType
type CoordType float64

/*
 検索時のデータType
*/
type SearchObjects struct {
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
	Id_str      string
	Name        string
	Screen_name string
	Profile_image_url string
}

/*
  タイムラインの取得
*/
func (this *Twitter) GetTimeline() ([]TweetObject, error) {
	response, err := this.oauth.Get(
		"https://api.twitter.com/1.1/statuses/home_timeline.json",
		map[string]string{"count": "10"})

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bits, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var tweets []TweetObject
	err = json.Unmarshal(bits, &tweets)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

/*
  ステータス更新API呼び出し
*/
func (this *Twitter) Update(status string) error {
	_, err := this.oauth.Post(
		"https://api.twitter.com/1.1/statuses/update.json",
		map[string]string{
			"status": status,
		})
	if err != nil {
		return err
	}
	return nil
}

/*
  Application-only auth の戻り値
*/
type AccessTokenAOA struct {
	Token_Type   string
	Access_Token string
}

/*
  Application-only Auth 用の検索
*/
func (this *Twitter) SearchAOA(word string) ([]TweetObject, error) {

	url := "https://api.twitter.com/oauth2/token"

	key := this.oauth.ConsumerKey + ":" + this.oauth.ConsumerSecret
	base64 := base64.StdEncoding.EncodeToString([]byte(key))

	wb := web.NewWeb()
	wb.AddHeader("Authorization", "Basic "+base64)
	wb.AddParam("grant_type", "client_credentials")

	resp, err := wb.Post(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bits, err := ioutil.ReadAll(resp.Body)
	at := AccessTokenAOA{}
	err = json.Unmarshal(bits, &at)

	url = "https://api.twitter.com/1.1/search/tweets.json"
	wb = web.NewWeb()
	wb.AddHeader("Authorization", "Bearer "+at.Access_Token)
	wb.AddParam("count", "100")
	//wb.AddParam("screen_name", "secondarykey")
	wb.AddParam("q", word)

	wb.AddParam("geocode", "38,138,500km")
	result, err := wb.Get(url)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	ioutil.WriteFile("searchData", data, 0666)

	//検索用のオブジェクト
	var tweets SearchObjects
	//取得する
	err = json.Unmarshal(data, &tweets)
	if err != nil {
		return nil, err
	}

	return tweets.Statuses, nil
}
