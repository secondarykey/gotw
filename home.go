package hello

import (
    "fmt"
    "html/template"
    "net/http"
)

func init() {
    http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	t,err:= template.ParseFiles("template/tweet.html")
	if err != nil{
		fmt.Println("error:", err)
	}
	t.Execute(w, nil)
}

