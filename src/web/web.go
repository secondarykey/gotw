package web

import (
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"fmt"
)

/*
 */
type Web struct {
	contentType string
	header      http.Header
	params *parameter
}

/*
 */
type parameter struct {
	param map[string]string
	order []string
}

/*
 * URLエスケープ
 */
func escape(param string) string {
	return url.QueryEscape(param)
}

/*
 * HTTP エラーコード
 */
type HttpError struct {
	status     string
	statusCode int
}

/*
 * エラー実装
 */
func (self HttpError) Error() string {
	return strconv.Itoa(self.statusCode) + ":\n" +
			self.status
}

/*
 * Webインスタンスの生成
 */
func NewWeb() *Web {
	return &Web{
		params: NewParameter(),
		header: http.Header{},
		contentType:"",
	}
}

/*
 * パラメータの追加
 */
func (self *Web) AddParam(key, value string) {
	self.params.add(key, value)
}

/*
 * パラメータ用のType
 * paramは保存場所
 * orderは引数のソート用に使用する
 */
func NewParameter() *parameter {
	return &parameter{
		param: make(map[string]string),
		order: make([]string, 0),
	}
}

/*
 * 引数の追加を行う
 */
func (self *parameter) add(key, value string) {
	self.addUnEscape(key, escape(value))
}

/*
 * addのエスケープをしない場合の呼び出し
 * 通常copyの時しか使わない
 */
func (self *parameter) addUnEscape(key, value string) {
	if _, flag := self.param[key]; !flag {
		self.param[key] = value
		self.order = append(self.order, key)
	}
}

/*
 * 値の取得
 */
func (self *parameter) get(key string) string {
	return self.param[key]
}

/*
 * 新規にparameterを生成して、自身のコピーをする
 */
func (self *parameter) copy() *parameter {
	clone := NewParameter()
	for _, key := range self.keys() {
		clone.addUnEscape(key, self.get(key))
	}
	return clone
}

/*
 * キーの取り出しsort.Strings()でソートしてから返す
 */
func (self *parameter) keys() []string {
	sort.Strings(self.order)
	return self.order
}

func (self *Web) getQuery() string {
	params := self.params.keys()
	ret := ""
	sep := ""
	for _, key := range params {
		value := self.params.get(key)
		ret += sep+key+"="+value
		sep = "&"
	}
	return ret
}

/*
 */
func (self *Web) Get(url string) (*http.Response, error) {
	q := self.getQuery()
	if q != "" {
		q += "?" + q
	}
	return self.execute("GET", url+q, "")
}

/*
 */
func (self *Web) Post(url string) (*http.Response, error) {
	self.contentType = "application/x-www-form-urlencoded"
	return self.execute("POST", url, self.getQuery())
}

/*
 */
func (self *Web) execute(method string, url string, body string) (*http.Response, error) {

	req, reqErr := http.NewRequest(method, url, strings.NewReader(body))
	if reqErr != nil {
		return nil, reqErr
	}

	//ヘッダの設定
	req.Header = self.header
	fmt.Println(self.header)
	if self.contentType != "" {
		req.Header.Set("Content-Type", self.contentType)
	}
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	client := &http.Client{}
	resp, doErr := client.Do(req)
	if doErr != nil {
		return nil, doErr
	}

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		defer resp.Body.Close()
		return nil, HttpError{
			status:     resp.Status,
			statusCode: resp.StatusCode,
		}
	}
	return resp, doErr
}
