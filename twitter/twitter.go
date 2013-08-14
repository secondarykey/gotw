package twitter

import (
	"fmt"
	"io/ioutil"
	"log"
	"oauth"
	"encoding/json"
)

type Twitter struct {
	consumer *oauth.Consumer
	Url string
	rtoken *oauth.RequestToken
	atoken *oauth.AccessToken
}

type TweetObject struct{
	Created_at string
	Id_str string
	Text string
	Source string
	In_reply_to_user_id_str string
	User UserObject
}

type UserObject struct{
	Id_str string
	Name string
	Screen_name string
}


func Usage() {
	fmt.Println("Usage:")
	fmt.Print("go run examples/twitter/twitter.go")
	fmt.Print("  --consumerkey <consumerkey>")
	fmt.Println("  --consumersecret <consumersecret>")
	fmt.Println("")
	fmt.Println("In order to get your consumerkey and consumersecret, you must register an 'app' at twitter.com:")
	fmt.Println("https://dev.twitter.com/apps/new")
}

func NewTwitter(key string,secret string) *Twitter {
	c := oauth.NewConsumer(key, secret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "http://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	})
	return &Twitter{
		consumer:c,
	}

}

func (this *Twitter) GetRequestTokenAndUrl(callback string) {
	requestToken, url, err := this.consumer.GetRequestTokenAndUrl(callback)
	if ( err != nil ) {
	}
	this.rtoken = requestToken
	this.Url = url
	return
}

func (this *Twitter) GetAccessToken(code string) {
	accessToken, err := this.consumer.AuthorizeToken(this.rtoken, code)
	if ( err != nil ) {
	}
	this.atoken = accessToken
	return
}

func (this *Twitter) GetTimeline() []TweetObject {
	response, err := this.consumer.Get(
		"http://api.twitter.com/1.1/statuses/home_timeline.json",
		map[string]string{"count": "10"},
		this.atoken)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	bits, err := ioutil.ReadAll(response.Body)
	var tweets []TweetObject
	err2 := json.Unmarshal(bits,&tweets)
	if err2 != nil{
		log.Print("error:", err2)
	}

	return tweets
}

func (this *Twitter) Update(status string) {
	_, err := this.consumer.Post(
		"http://api.twitter.com/1.1/statuses/update.json",
		map[string]string{
			"status": status,
		},
		this.atoken)
	if err != nil {
	}
}
