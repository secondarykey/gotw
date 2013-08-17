package gotw

import (
    "fmt"
    "html/template"
    "net/http"
    "appengine"
    "appengine/user"
    "appengine/datastore"
)

type Config struct {
	Token string
	Secret string
}

type Agent struct {
	Id string
	Name string
	Email string
}

type Authentication struct {
	Id string
	SNSId string
	Name string
	Token string
	Secret string
}

func init() {
    http.HandleFunc("/", topHandler)
    http.HandleFunc("/logout", logoutHandler)
    http.HandleFunc("/home", homeHandler)
    http.HandleFunc("/settings/", settingsHandler)
    http.HandleFunc("/settings/wtite", writeAgentHandler)
    http.HandleFunc("/settings/twitter/add", addTwitterHandler)
    http.HandleFunc("/settings/twitter/callback", callbackHandler)
}

func initConfig(w http.ResponseWriter,r *http.Request) {
	c := appengine.NewContext(r)
	//設定情報がない場合
	q := datastore.NewQuery("Config")
	count ,err := q.Count(c)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	if count == 0 {
		config := Config {
			Token  : "ysTHKkYBW9PrHtgtYyElsA",
			Secret : "Ofl3NvzYGQKeNghBZ8KP1HMcZELxfv7dVnacjpDHvQ",
		}
		datastore.Put(c,datastore.NewIncompleteKey(c,"Config",nil),&config)
	}
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	handler(w,"template/top.html",nil)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	url, _ := user.LogoutURL(c, "/")
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	initConfig(w,r)

	c := appengine.NewContext(r)
	u := user.Current(c)
	//IDでDatastoreを検索
	q := datastore.NewQuery("Agent").Filter("Id =", u.ID)
	count ,err := q.Count(c)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.Redirect(w, r, "/settings/", http.StatusMovedPermanently)
	} else {
		//認証情報を取得


		handler(w,"template/home.html",u)
	}
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {

	//ユーザ情報から設定値を取得
	c := appengine.NewContext(r)
	u := user.Current(c)

	handler(w,"template/settings.html",u)
}
func writeAgentHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	u := user.Current(c)

	//IDでDatastoreを検索
	q := datastore.NewQuery("Agent").Filter("Id =", u.ID)
	count ,_ := q.Count(c)
	if ( count == 0 ) {
		agent := Agent {
			Name:r.FormValue("name"),
			Email:r.FormValue("email"),
			Id:u.ID,
		}
		datastore.Put(c,datastore.NewIncompleteKey(c,"Agent",nil),&agent)
	}

	http.Redirect(w, r, "/settings/", http.StatusMovedPermanently)
}

func addTwitterHandler(w http.ResponseWriter, r *http.Request) {

	//データの抜き出し

	//登録

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {

	//データの抜き出し

	//登録

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func handler(w http.ResponseWriter,templateName string,data interface{}) {
	t,err:= template.ParseFiles(templateName)
	if err != nil{
		fmt.Println("error:", err)
	}
	t.Execute(w, data)
}
