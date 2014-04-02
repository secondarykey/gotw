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
	header http.Header
	params *parameter
}

/*
 */
type parameter struct {
	param map[string]string
	order []string
}

/*
 */
func escape(param string) string {
	return url.QueryEscape(param)
}

/*
 */
type HttpError struct {
	status     string
	statusCode int
}

/*
 */
func (self HttpError) Error() string {
	return strconv.Itoa(self.statusCode) + ":\n" +
		self.status
}

/*
 */
func NewWeb() *Web {
	return &Web{
		params: NewParams(),
		header: http.Header{},
		contentType:"",
	}
}


/*
 */
func (self *Web) AddParam(key, value string) {
	self.params.add(key, value)
}

/*
 */
func NewParams() *parameter {
	return &parameter{
		param: make(map[string]string),
		order: make([]string, 0),
	}
}

/*
 */
func (self *parameter) add(key, value string) {
	self.addUnEscape(key,escape(value))
}

func (self *parameter) addUnEscape(key, value string) {
	self.param[key] = value
	self.order = append(self.order, key)
}

func (self *parameter) Get(key string) string {
	return self.param[key]
}

func (self *parameter) Copy() *parameter {
	clone := NewParams()
	for _,key := range self.Keys() {
		clone.addUnEscape(key,self.Get(key))
	}
	return clone
}

/*
 */
func (self *parameter) Keys() []string {
	sort.Strings(self.order)
	return self.order
}

func (self *Web) getQuery() string {
	params := self.params.Keys()
	ret := ""
	sep := ""
	for _, key := range params {
		value := self.params.Get(key)
		ret += sep + escape(key) + "=" + escape(value)
		sep = "&"
	}
	return ret
}

/*
 */
func (self *Web) Get(url string) (*http.Response, error) {
	return self.execute("GET", url + "?" + self.getQuery(),"")
}

/*
 */
func (self *Web) Post(url string) (*http.Response, error) {
	self.contentType = "application/x-www-form-urlencoded"
	return self.execute("POST", url,self.getQuery())
}

/*
 */
func (self *Web) execute(method string, url string,body string) (*http.Response, error) {

	req, reqErr := http.NewRequest(method, url, strings.NewReader(body))
	if reqErr != nil {
		return nil, reqErr
	}

	//ヘッダの設定
	req.Header = self.header
	fmt.Println(self.header)
	if self.contentType != "" {
		req.Header.Set("Content-Type",self.contentType)
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
