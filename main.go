package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"oauth"
	"encoding/json"
	"./twitter"
)

func main() {

	var consumer oauth.AccessToken
	readToken(&consumer,"consumer.json")

	t := twitter.NewTwitter(consumer.Token,consumer.Secret)

	//var access oauth.AccessToken
	//err = readToken(&access,"access.json")
	//if err != nil {
		t.GetRequestTokenAndUrl("oob")

		fmt.Println("Go to: " + t.Url)
	
		verificationCode := ""
		fmt.Scanln(&verificationCode)

		t.GetAccessToken(verificationCode)

		//アクセストークンの保存
		/*
		err = writeToken(accessToken, "access.json")
		if err != nil {
			log.Fatal(err)
		}
		access = *accessToken
		*/
	//}

	tweets := t.GetTimeline()
	for _,tweet := range tweets {
		log.Print(tweet.User.Name)
		log.Print("----" + tweet.Text)
	}

	for {
		fmt.Print("> ")
		command := ""
		fmt.Scanln(&command)
		end := cmd(t,command)
		if ( end == false) {
			break
		}
	}
	fmt.Println("Bye!")
}

func cmd(t *twitter.Twitter,cmd string) bool {

	switch cmd {
		case "timeline":
			tweets := t.GetTimeline()
			for _,tweet := range tweets {
				log.Print(tweet.User.Name)
				log.Print("----" + tweet.Text)
			}
		case "q":
			return false
		default :
			fmt.Print("Sending status?[Y/n]:")
			ans := ""
			fmt.Scanln(&ans)
			if ( ans == "Y" ) {
				t.Update(cmd)
			}
	}

	return true
}

func readToken(token interface{}, filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, token)
}

func writeToken(token interface{}, filename string) error {
	b, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0666)
}

