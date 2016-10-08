package main

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

func (this *Twitter) GetTimeline() ([]TweetObject, error) {

	resp, err := this.oauth.Get(
		"https://api.twitter.com/1.1/statuses/home_timeline.json",
		map[string]string{"count": "10"})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bits, err := ioutil.ReadAll(resp.Body)
	var tweets []TweetObject
	err = json.Unmarshal(bits, &tweets)
	if err != nil {
		return nil, err
	}

	return tweets, nil
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
