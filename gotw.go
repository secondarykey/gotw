package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {

	/*
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
	*/

	twt, err := getTwitter()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = setAccessToken(twt)
	if err != nil {
		fmt.Println(err)
		return
	}

	wait(twt)
}

func getTwitter() (*Twitter, error) {

	c := Credential{}
	err := read(&c)
	if err != nil {
		return nil, fmt.Errorf("設定ファイル読み込みエラー:%s", err.Error())
	}
	t := NewTwitter(&c)
	return t, nil
}

func setAccessToken(t *Twitter) error {

	if t.oauth.Credential.AccessToken != "" {
		return nil
	}

	t.SetRequestTokenAndUrl("oob")

	fmt.Println("認可をおこなってください: " + t.GetAuthorizationUrl())
	verificationCode := ""

	fmt.Print("[PIN]=")
	fmt.Scanln(&verificationCode)

	err := t.Exchange(verificationCode)
	if err != nil {
		return err
	}
	return write(t.oauth.Credential)
}

func cmd(t *Twitter, cmd string) bool {

	switch cmd {
	case "t":
		tweets, err := t.GetTimeline()
		if err != nil {
			fmt.Println("Error:" + err.Error())
			return true
		}

		//30	黒
		//31	赤
		//32	緑
		//33	黄
		//34	青
		//35	マゼンダ
		//36	シアン
		//37	白
		for _, tweet := range tweets {
			fmt.Println(tweet.User.Name, tweet.User.Screen_name, tweet.Created_at, "------")
			color := "\x1b[32m%s\x1b[0m\n"
			fmt.Printf(color, tweet.Text)
		}
	case "q":
		return false
	case "":
		return true
	default:
		fmt.Print("Sending status?[Y/n]:")
		ans := ""
		fmt.Scanln(&ans)
		if ans == "Y" {
			t.Update(cmd)
		}
	}
	return true
}

func wait(t *Twitter) {
	cmd(t, "t")
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
