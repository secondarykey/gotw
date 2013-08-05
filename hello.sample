package hello

import (
    "flag"
    "fmt"
    "net/http"
    "io/ioutil"
    "log"
    "oauth"
    "encoding/json"
    "html/template"
)

func init() {
    http.HandleFunc("/", handler)
}

var (
	clientid     = flag.String("id", "your client id", "OAuth Client ID")
	clientsecret = flag.String("secret", "your secret id", "OAuth Client Secret")
	rtokenfile   = flag.String("request", "request.json", "Request token file name")
	atokenfile   = flag.String("access", "access.json", "Access token file name")
	code         = flag.String("code", "your code", "Verification code")
)


var provider = oauth.ServiceProvider{
	RequestTokenUrl:   "http://api.twitter.com/oauth/request_token",
	AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
	AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
}

// 取得したいパラメータをstructで記述
// 参考 https://dev.twitter.com/docs/api/1/get/statuses/mentions
type TweetObject struct{
	Created_at string
	Id_str string
	Text string
	Source string
	In_reply_to_user_id_str string
	User UserObject // JSONオブジェクト内のオブジェクトをこのように定義する。配列だと array []UserArrayみたいな感じ
}

// JSONオブジェク内のオブジェクト
type UserObject struct{
	Id_str string
	Name string
	Screen_name string
}


func handler(w http.ResponseWriter, r *http.Request) {
	flag.Parse()
	if *clientid == "" || *clientsecret == "" {
		flag.Usage()
		return
	}
	consumer := oauth.NewConsumer(*clientid, *clientsecret, provider)

	var atoken oauth.AccessToken
	err := readToken(&atoken, *atokenfile)
	if err != nil {
		log.Print("Couldn't read token:", err)

		var rtoken oauth.RequestToken
		err := readToken(&rtoken, *rtokenfile)
		if err != nil {
			log.Print("Couldn't read token:", err)
			log.Print("Getting Request Token")
			rtoken, url, err := consumer.GetRequestTokenAndUrl("oob")
			if err != nil {
				log.Fatal(err)
			}
			err = writeToken(rtoken, *rtokenfile)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Visit this URL:", url)
			fmt.Println("Then run this program again with -code=CODE")
			fmt.Println("where CODE is the verification PIN provided by Twitter.")
			return
		}

		log.Print("Getting Access Token")
		if *code == "" {
			fmt.Println("You must supply a -code parameter to get an Access Token.")
			return
		}
		tok, err := consumer.AuthorizeToken(&rtoken, *code)
		if err != nil {
			log.Fatal(err)
		}
		err = writeToken(tok, *atokenfile)
		if err != nil {
			log.Fatal(err)
		}
		atoken = *tok
	}

//	const url = "http://api.twitter.com/1/statuses/mentions.json"
	const url = "http://api.twitter.com/1/statuses/user_timeline.json"
	log.Print("GET ", url)

	resp, err := consumer.Get(url,nil, &atoken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return	 
	}
	
	w.Header().Add("Content-type","text/html charset=utf-8")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var tweets []TweetObject
	err2 := json.Unmarshal(body,&tweets) 
	if err2 != nil{
		fmt.Println("error:", err2)
	}

	t,err:= template.ParseFiles("template/tweet.html") 	// gotweet_change/template/tweet.html
	if err != nil{
		fmt.Println("error:", err)
	}
	t.Execute(w, tweets)

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


