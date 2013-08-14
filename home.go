package gotw

import (
    "fmt"
    "html/template"
    "net/http"
)

func init() {
    http.HandleFunc("/", topHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/home", homeHandler)
    http.HandleFunc("/callback", callbackHandler)
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	handler("template/top.html")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	//データ存在判定
		//存在する場合
		//ホームにリダイレクト
		//return

	//Twitter登録のURLに飛ぶ
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	//tweetを取得

	handler("template/home.html")
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	//データの抜き出し
	//登録
	//homeにリダイレクト
}

func handler(w http.ResponseWriter,templateName string) {
	t,err:= template.ParseFiles(templateName)
	if err != nil{
		fmt.Println("error:", err)
	}
	t.Execute(w, nil)
}
