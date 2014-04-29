package main

// 匿名フィールド(Methodの継承)
// interface{} Any型
// 組み込みフィールド

import (
	"encoding/json"
	"fmt"
	"github.com/secondarykey/golib/oauth"
	"io/ioutil"
	"net/http"
	"html/template"
	"twitter"
	"strconv"
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

	http.HandleFunc("/", handler)
	http.HandleFunc("/rain.json", rainHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe("localhost:4000", nil)
}

func rainHandler(w http.ResponseWriter, r *http.Request) {

	twt := getTwitter()
	token,err := twt.GetAOAToken()
	if err != nil {
		panic(err)
	}

	word:="アメ"
	twt.AddParam("count", "100")
	so, err := twt.SearchAOA(token,word)
	if err != nil {
		panic(err)
	}

	var jsonTweet []twitter.TweetObject
	maxid := int64(0)
	tCnt := 0
	pCnt := 0
	for {
		fmt.Println("ループ " + strconv.FormatInt(int64(len(jsonTweet)),10))
		if so == nil || len(so.Statuses) <= 0 {
			fmt.Println("終了")
			break
		}
		for _, tweet := range so.Statuses {
			tCnt++

			if maxid==0 || tweet.Id < maxid {
				maxid = tweet.Id
			}

			if tweet.Geo.Coordinates[0] != 0.0 {
				jsonTweet = append(jsonTweet, tweet)
			}

			if tweet.User.Location != "" {
				pCnt++
				fmt.Println(tweet.User.Location)
			}
		}

		if len(jsonTweet) >= 100 {
			break
		}

		fmt.Printf("T:%d,Location=%d",tCnt,pCnt)

		twt.AddParam("max_id", strconv.FormatInt(maxid,10))
		twt.AddParam("count", "100")
		so, err = twt.SearchAOA(token,word)
	}

	bits, err := json.Marshal(jsonTweet);
	if err != nil {
		panic(err)
	}
	w.Write(bits)
}

func handler(w http.ResponseWriter, r *http.Request) {

	twt := getTwitter()
	setAccessToken(twt)

	t, err := template.ParseFiles("template/index.html")
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
		writeJson(t.GetToken(), "access.json")
	} else {
		t.SetAccessToken(&tokenSet)
	}
	return
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
