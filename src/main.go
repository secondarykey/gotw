package main

// 匿名フィールド(Methodの継承)
// interface{} Any型
// 組み込みフィールド

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"twitter"
	"github.com/secondarykey/golib/oauth"
	"strings"
	"net/http"
	"html/template"
)

type GotwError struct {
	Title   string
	Message string
}

func main() {

	/*
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
	*/

	//wait(twt)

	http.HandleFunc("/",handler)
	http.HandleFunc("/rain.json",rainHandler)

	http.Handle("/static/", http.StripPrefix("/static/",http.FileServer(http.Dir("static"))))
	http.ListenAndServe("localhost:4000", nil)
}

func rainHandler(w http.ResponseWriter,r *http.Request) {

	twt := getTwitter()
	tweets,err := twt.SearchAOA("雨")
	if err != nil {
		panic(err)
	}

	var jsonTweet []twitter.TweetObject
	for _, tweet := range tweets {
		if tweet.Geo.Coordinates[0] != 0.0 {
			jsonTweet = append(jsonTweet,tweet)
		}
	}

	bits, err := json.Marshal(jsonTweet);
	w.Write(bits)
}

func handler(w http.ResponseWriter,r *http.Request) {

	twt := getTwitter()
	setAccessToken(twt)

	t,err:= template.ParseFiles("template/index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

/*
 * consumer.jsonからTwitterオブジェクトを生成
 */
func getTwitter() *twitter.Twitter {
	var tokenSet oauth.TokenSet
	err := readJson(&tokenSet, "consumer.json")
	if err != nil {
		panic(err)
	}
	t := twitter.NewTwitter(tokenSet.Token, tokenSet.Secret)
	return t
}

/*
 * アクセストークンファイルがある場合は読み込んで設定
 * 存在しない場合はコードをリクエストトークンを生成して、取得しにいく
 */
func setAccessToken(t *twitter.Twitter) {
	var tokenSet oauth.TokenSet
	err := readJson(&tokenSet, "access.json")
	if err != nil {
		t.SetRequestTokenAndUrl("oob")

		fmt.Println("認証情報なしなので右にアクセス: " + t.GetAuthorizationUrl())
		verificationCode := ""

		fmt.Print(">")
		fmt.Scanln(&verificationCode)

		t.GetAccessToken(verificationCode)
		writeJson(t.GetToken(),"access.json")
	} else {
		t.SetAccessToken(&tokenSet)
	}
	return
}

/*
 * コマンドに応じた処理を行う
 */
func cmd(t *twitter.Twitter, cmd string) bool {

	if cmd == "" {
		return true
	}

	switch {
	case cmd == "timeline":
		tweets, err := t.GetTimeline()
		if err != nil {
			fmt.Println(err)
		} else {
			printTweet(tweets)
		}
	case cmd == "q":
		return false

	case strings.HasPrefix(cmd,"search ") == true:

		word := cmd[7:]
		tweets, err := t.SearchAOA(word)
		if err != nil {
			fmt.Println(err)
		} else {
			printTweet(tweets)
		}
	default:
		fmt.Print("Sending status?[Y/n]:")
		ans := ""
		fmt.Scanln(&ans)
		if ans == "Y" {
			err := t.Update(cmd)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return true
}

func printTweet(tweets []twitter.TweetObject) {
	for _, tweet := range tweets {
		fmt.Println(tweet.User.Name, tweet.User.Screen_name, tweet.Created_at, "------")
		fmt.Println(tweet.Text)
	}
}

/*
 * コマンドの永久ループ
 */
func wait(t *twitter.Twitter) {
	cmd(t, "timeline")
	for {
		fmt.Print("> ")
		command := ""
		fmt.Scanln(&command)
		end := cmd(t, command)
		if end == false {
			break
		}
	}
	fmt.Println("Bye!")
}

/*
 * Json形式のファイルをTypeに読み込む
 */
func readJson(token interface{}, filename string) error {
	if b, err := ioutil.ReadFile(filename); err != nil {
		return err
	} else {
		return json.Unmarshal(b, token)
	}
}

/*
 * TypeをJson形式で書き込む
 */
func writeJson(token interface{}, filename string) error {
	if b, err := json.Marshal(token); err != nil {
		return err
	} else {
		return ioutil.WriteFile(filename, b, 0666)
	}
}
