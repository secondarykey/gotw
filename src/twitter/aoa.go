package twitter

import (
	"encoding/base64"
	"github.com/secondarykey/golib/web"
	"io/ioutil"
	"encoding/json"
)

//This file is Twitter Application only authrization
/*
  Application-only auth の戻り値
*/
type AccessTokenAOA struct {
	Token_Type   string
	Access_Token string
}

func (this *Twitter) GetAOAToken() (string,error) {

	url := "https://api.twitter.com/oauth2/token"
	key := this.oauth.ConsumerKey + ":" + this.oauth.ConsumerSecret
	base64 := base64.StdEncoding.EncodeToString([]byte(key))

	wb := web.NewWeb()
	wb.AddHeader("Authorization", "Basic "+base64)
	wb.AddParam("grant_type", "client_credentials")
	resp, err := wb.Post(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bits, err := ioutil.ReadAll(resp.Body)
	var at AccessTokenAOA
	err = json.Unmarshal(bits, &at)
	if err != nil {
		return "", err
	}
	return at.Access_Token,nil
}

func (this *Twitter) SearchAOA(token,word string) (*SearchObject, error) {

	url := "https://api.twitter.com/1.1/search/tweets.json"
	wb := web.NewWeb()
	wb.AddHeader("Authorization", "Bearer "+token)

	for _,key := range this.param.Keys() {
		wb.AddParam(key, this.param.Get(key))
	}
	//使い回しの為削除
	this.param = web.NewParameter()
	wb.AddParam("q", word)

	result, err := wb.Get(url)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	//検索用のオブジェクト
	var so SearchObject
	err = web.ReadJson(result, &so)
	if err != nil {
		return nil, err
	}
	return &so, nil
}

