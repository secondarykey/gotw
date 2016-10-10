package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type Web struct {
	contentType string
	header      http.Header
	params      *parameter
}

type parameter struct {
	param map[string]string
	order []string
}

func escape(param string) string {
	return url.QueryEscape(param)
}

func NewWeb() *Web {
	return &Web{
		params:      NewParameter(),
		header:      http.Header{},
		contentType: "",
	}
}

func (w *Web) AddParam(key, value string) {
	w.params.add(key, value)
}

func NewParameter() *parameter {
	return &parameter{
		param: make(map[string]string),
		order: make([]string, 0),
	}
}

func (p *parameter) add(key, value string) {
	p.addUnEscape(key, escape(value))
}

func (p *parameter) addUnEscape(key, value string) {
	if _, flag := p.param[key]; !flag {
		p.param[key] = value
		p.order = append(p.order, key)
	}
}

func (p *parameter) get(key string) string {
	return p.param[key]
}

func (p *parameter) copy() *parameter {
	clone := NewParameter()
	for _, key := range p.keys() {
		clone.addUnEscape(key, p.get(key))
	}
	return clone
}

func (p *parameter) keys() []string {
	sort.Strings(p.order)
	return p.order
}

func (w *Web) getQuery() string {
	params := w.params.keys()
	ret := ""
	sep := ""
	for _, key := range params {
		value := w.params.get(key)
		ret += sep + key + "=" + value
		sep = "&"
	}
	return ret
}

func (w *Web) Get(url string) (*http.Response, error) {
	q := w.getQuery()
	if q != "" {
		q = "?" + q
	}
	return w.execute("GET", url+q, "")
}

func (w *Web) Post(url string) (*http.Response, error) {
	w.contentType = "application/x-www-form-urlencoded"
	return w.execute("POST", url, w.getQuery())
}

func (w *Web) execute(method string, url string, body string) (*http.Response, error) {

	req, reqErr := http.NewRequest(method, url, strings.NewReader(body))
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header = w.header
	if w.contentType != "" {
		req.Header.Set("Content-Type", w.contentType)
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
		return nil, fmt.Errorf("[%d]%s", resp.StatusCode, resp.Status)
	}

	return resp, doErr
}
