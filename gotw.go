package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

var maxId int64

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic!!!", err)
			os.Exit(-1)
		}
	}()

	twt, err := createTwitterInformation()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = setAccessToken(twt)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = drawTimeline(twt)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		ti := time.NewTicker(60 * time.Second)
		for {
			select {
			case <-ti.C:
				err := drawTimeline(twt)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print("> ")
			}
		}
		ti.Stop()
	}()

	wait(twt)
}

func createTwitterInformation() (*Twitter, error) {

	c := Credential{}
	err := read(&c)
	if err != nil {
		return nil, fmt.Errorf("設定ファイル読み込みエラー:%s", err.Error())
	}

	if c.ConsumerKey == "" || c.ConsumerSecret == "" {
		return nil, fmt.Errorf(".credentialにConsumerKeyとConsumerSecretを設定してください")
	}

	t := NewTwitter(&c)
	return t, nil
}

func setAccessToken(t *Twitter) error {

	if t.OAuth.Credential.AccessToken != "" {
		return nil
	}

	err := t.GetRequestToken("oob")
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

	return write(t.OAuth.Credential)
}

func drawTimeline(t *Twitter) error {

	tweets, err := t.GetTimeline(int(maxId))
	if err != nil {
		fmt.Println("Error:" + err.Error())
		return err
	}
	if len(tweets) <= 0 {
		return nil
	}

	sort.Sort(tweets)
	//30	黒
	//31	赤
	//32	緑
	//33	黄
	//34	青
	//35	マゼンダ
	//36	シアン
	//37	白
	for idx, tweet := range tweets {
		fmt.Printf("%s(@%s):\n", tweet.User.Name, tweet.User.Screen_name)
		color := "\x1b[%dm%s\x1b[0m\n"
		num := 32
		if (idx % 2) == 1 {
			num = 36
		}
		fmt.Printf(color, num, tweet.Text)
		fmt.Printf("--- %s\n", changeTime(tweet.Created_at))

		if tweet.Id > maxId {
			maxId = tweet.Id
		}
	}
	return nil
}

func cmd(t *Twitter, cmd string) bool {

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

func wait(t *Twitter) {
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

func read(c *Credential) error {
	if b, err := ioutil.ReadFile(".credential"); err != nil {
		return err
	} else {
		return json.Unmarshal(b, c)
	}
}

func write(c *Credential) error {
	if b, err := json.Marshal(c); err != nil {
		return err
	} else {
		return ioutil.WriteFile(".credential", b, 0666)
	}
}
