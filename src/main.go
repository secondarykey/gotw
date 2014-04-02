package main

//
// 匿名フィールド(Methodの継承)
// interface{} Any型
// 組み込みフィールド

import (
	"twitter"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"web"
)

type GotwError struct {
	Title   string
	Message string
}

func (this *GotwError) Error() string {
	return fmt.Sprintf("%s:\n%s", this.Title, this.Message)
}

func NewError(title, message string) GotwError {
	return GotwError{title, message}
}

func main() {

	/*
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	*/

	twt := getTwitter()
	setAccessToken(twt)

	wait(twt)
}

func getTwitter() *twitter.Twitter {
	var tokenSet web.TokenSet
	err := readJson(&tokenSet, "consumer.json")
	if err != nil {
		panic(NewError("consumer.json読み込みエラー", err.Error()))
	}
	t := twitter.NewTwitter(tokenSet.Token, tokenSet.Secret)
	return t
}

func setAccessToken(t *twitter.Twitter) {
	var tokenSet web.TokenSet
	err := readJson(&tokenSet, "access.json")
	if err != nil {
		t.SetRequestTokenAndUrl("oob")

		fmt.Println("認証情報なしなので右にアクセス: " + t.GetAuthorizationUrl())
		verificationCode := ""

		fmt.Print(">")
		fmt.Scanln(&verificationCode)

		t.GetAccessToken(verificationCode)
	} else {
		t.SetAccessToken(&tokenSet)
	}
	return
}

func cmd(t *twitter.Twitter, cmd string) bool {

	switch cmd {
	case "timeline":
		tweets := t.GetTimeline()
		for _, tweet := range tweets {
			fmt.Println(tweet.User.Name, tweet.User.Screen_name, tweet.Created_at, "------")
			fmt.Println(tweet.Text)
		}
	case "q":
		return false
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
 *
 *
 *
 *
 */
func readJson(token interface{}, filename string) error {
	if b, err := ioutil.ReadFile(filename); err != nil {
		return err
	} else {
		return json.Unmarshal(b, token)
	}
}

/*
 *
 *
 *
 *
 *
 */
func writeJson(token interface{}, filename string) error {
	if b, err := json.Marshal(token); err != nil {
		return err
	} else {
		return ioutil.WriteFile(filename, b, 0666)
	}
}
