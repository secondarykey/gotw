package twitter

import (
	"encoding/json"
	"io/ioutil"
)

type TweetObject struct {
	Created_at              string
	Id_str                  string
	Text                    string
	Source                  string
	In_reply_to_user_id_str string
	User                    UserObject
}

type UserObject struct {
	Id_str      string
	Name        string
	Screen_name string
}

func (this *Twitter) GetTimeline() []TweetObject {
	response, err := this.oauth.Get(
		"https://api.twitter.com/1.1/statuses/home_timeline.json",
		map[string]string{"count": "10"})
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	bits, err := ioutil.ReadAll(response.Body)
	var tweets []TweetObject
	err2 := json.Unmarshal(bits, &tweets)
	if err2 != nil {
		panic(err2)
	}

	return tweets
}

func (this *Twitter) Update(status string) {
	_, err := this.oauth.Post(
		"https://api.twitter.com/1.1/statuses/update.json",
		map[string]string{
			"status": status,
		})
	if err != nil {
		panic(err)
	}
}
