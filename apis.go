package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/secondarykey/golib/http"
)

type Twitter struct {
	*http.OAuth1
}

func NewTwitter(c *http.Credential) *Twitter {

	oa := http.NewOAuth1(c,
		"https://api.twitter.com/oauth/request_token",
		"https://api.twitter.com/oauth/authorize",
		"https://api.twitter.com/oauth/access_token",
	)
	return &Twitter{
		OAuth1: oa,
	}
}

type TweetObject struct {
	Created_at              string
	Id                      int64
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

type TweetList []TweetObject

func (t TweetList) Len() int {
	return len(t)
}

func (t TweetList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TweetList) Less(i, j int) bool {
	return t[i].Id < t[j].Id
}

func (t *Twitter) GetTimeline(id int64) (TweetList, error) {

	count := "10"
	if id != 0 {
		count = "100"
	}

	arg := map[string]string{"count": count}
	if id != 0 {
		arg["since_id"] = fmt.Sprintf("%d", id)
	}

	resp, err := t.Get(
		"https://api.twitter.com/1.1/statuses/home_timeline.json",
		arg)
	if err != nil {
		return nil, fmt.Errorf("TwitterAPI Timeline Error:%s", err)
	}

	defer resp.Body.Close()

	bits, err := ioutil.ReadAll(resp.Body)
	var tweets TweetList
	err = json.Unmarshal(bits, &tweets)
	if err != nil {
		return nil, fmt.Errorf("JSON Unmarshal Error:%s", err)
	}

	return tweets, nil
}

func (t *Twitter) Update(status string) error {
	_, err := t.Post(
		"https://api.twitter.com/1.1/statuses/update.json",
		map[string]string{
			"status": status,
		})
	return fmt.Errorf("Twitter get Error:%s", err)
}

func (t *Twitter) Search(word string, maxId int) error {

	args := map[string]string{"count": "100"}
	if maxId != 0 {
		args["since_id"] = fmt.Sprintf("%d", maxId)
	}
	args["q"] = word

	resp, err := t.Get(
		"https://api.twitter.com/1.1/search/tweets.json", args)
	defer resp.Body.Close()

	bits, err := ioutil.ReadAll(resp.Body)

	var data interface{}
	err = json.Unmarshal(bits, &data)
	if err != nil {
		return fmt.Errorf("JSON Unmarshal Error:%s", err)
	}

	fmt.Println(data)
	return nil
}
