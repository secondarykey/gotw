package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/secondarykey/golib/http"
	"github.com/secondarykey/golib/util"
)

var maxId int64
var idBox map[string]*TweetObject

const CREDENTIAL_FILE = ".credential"

func init() {
	rand.Seed(time.Now().UnixNano())
	idBox = make(map[string]*TweetObject)
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic!!!", err)
			os.Exit(-1)
		}
	}()

	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}

func run() error {

	twt, err := createTwitterInformation()
	if err != nil {
		return err
	}

	err = setAccessToken(twt)
	if err != nil {
		return err
	}

	err = drawTimeline(twt)
	if err != nil {
		return err
	}

	go func() {
		ti := time.NewTicker(65 * time.Second)
		for {
			select {
			case <-ti.C:

				fmt.Println("[reload...]")
				err := drawTimeline(twt)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print("> ")
			}
		}
		ti.Stop()
	}()

	return wait(twt)
}

func createTwitterInformation() (*Twitter, error) {

	c := http.Credential{}

	_, err := os.Stat(CREDENTIAL_FILE)
	if err != nil {
		fmt.Print("[Consumer Key]=")
		fmt.Scanln(&c.ConsumerKey)
		fmt.Print("[Consumer Secret]=")
		fmt.Scanln(&c.ConsumerSecret)
	} else {
		err := util.ReadJsonFile(&c, CREDENTIAL_FILE)
		if err != nil {
			return nil, fmt.Errorf("設定ファイル読み込みエラー:%s", err.Error())
		}
	}

	if c.ConsumerKey == "" || c.ConsumerSecret == "" {
		return nil, fmt.Errorf("ConsumerKeyとConsumerSecretを設定してください")
	}

	t := NewTwitter(&c)
	return t, nil
}

func setAccessToken(t *Twitter) error {

	if t.OAuth1.Credential.AccessToken != "" {
		return nil
	}

	err := t.GetRequestToken("oob")
	if err != nil {
		return err
	}

	err = util.WriteJsonFile(t.OAuth1.Credential, CREDENTIAL_FILE)
	if err != nil {
		return err
	}

	fmt.Println("認可をおこなってください: " + t.GetAuthorizeURL())
	verificationCode := ""

	fmt.Print("[PIN]=")
	fmt.Scanln(&verificationCode)

	err = t.GetAccessToken(verificationCode)
	if err != nil {
		return err
	}

	return util.WriteJsonFile(t.OAuth1.Credential, CREDENTIAL_FILE)
}

func drawTimeline(t *Twitter) error {

	tweets, err := t.GetTimeline(maxId)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		return err
	}
	if len(tweets) <= 0 {
		return nil
	}

	sort.Sort(tweets)

	//30黒 31赤 32緑 33黄 34青 35マゼンダ 36シアン 37白

	for idx, tweet := range tweets {

		t := changeTime(tweet.Created_at)
		id := randId(&tweet)

		l := fmt.Sprintf("%s[%s]: %s(@%s)\n", t, id, tweet.User.Name, tweet.User.Screen_name)

		num := 32
		if (idx % 2) == 1 {
			num = 36
		}

		fmt.Print(color(num, l))

		// 自分への返信を赤表示

		fmt.Println(tweet.Text)

		if tweet.Id > maxId {
			maxId = tweet.Id
		}
	}
	return nil
}

func color(num int, txt string) string {
	c := "\x1b[%dm%s\x1b[0m"
	return fmt.Sprintf(c, num, txt)
}

func cmd(t *Twitter, cmd string) bool {

	// TODO add command
	// m -> mentions ... msg???
	//

	switch cmd {
	case "s":
		t.Search("#NowPlaying", 0)
	case "q":
		return false
	case "":
		return true
	default:
		fmt.Print("Sending status?[Y/n]:")
		ans := ""
		fmt.Scanln(&ans)
		if ans != "Y" {
			return true
		}

		err := t.Update(cmd)
		if err != nil {
			fmt.Println("Error:" + err.Error())
		}
	}
	return true
}

func changeTime(t string) string {

	ti, err := time.Parse(time.RubyDate, t)
	if err != nil {
		return err.Error()
	}
	jst, _ := time.LoadLocation("Asia/Tokyo")
	jt := ti.In(jst)
	return jt.Format("2006/01/02 15:04:05")
}

func wait(t *Twitter) error {

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
	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randId(tweet *TweetObject) string {
	b := make([]rune, 4)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	rId := string(b)

	idBox[rId] = tweet
	return rId
}
