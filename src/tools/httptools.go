package tools

import (
	"net/http"
	"net/http/httptrace"
	"time"
	"io/ioutil"
	"crypto/tls"
)

type Result struct {
	Err           error
	StatusCode    int
	Duration      time.Duration
	ConnDuration  time.Duration // connection setup(DNS lookup + Dial up) duration
	DnsDuration   time.Duration // dns lookup duration
	ReqDuration   time.Duration // request "write" duration
	ResDuration   time.Duration // response "read" duration
	DelayDuration time.Duration // delay between response and request
	ContentLength int64
	Body          []byte
}

func SimpleHttpClient() *http.Client {
	tr := &http.Transport{
		DisableKeepAlives: true,
		MaxIdleConns:1,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Proxy:              http.ProxyURL(nil),
	}

	client := &http.Client{Transport: tr, Timeout: time.Duration(8) * time.Second}
	return client
}

func MakeHttpRequest(client *http.Client, req *http.Request) *Result {
	s := time.Now()
	var dnsStart, connStart, resStart, reqStart, delayStart time.Time
	var dnsDuration, connDuration, resDuration, reqDuration, delayDuration time.Duration
	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			dnsDuration = time.Now().Sub(dnsStart)
		},
		GetConn: func(h string) {
			connStart = time.Now()
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			connDuration = time.Now().Sub(connStart)
			reqStart = time.Now()
		},
		WroteRequest: func(w httptrace.WroteRequestInfo) {
			reqDuration = time.Now().Sub(reqStart)
			delayStart = time.Now()
		},
		GotFirstResponseByte: func() {
			delayDuration = time.Now().Sub(delayStart)
			resStart = time.Now()
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	result := &Result{
		Err:           err,
		ConnDuration:  connDuration,
		DnsDuration:   dnsDuration,
		ReqDuration:   reqDuration,
		ResDuration:   resDuration,
		DelayDuration: delayDuration,
	}
	if err == nil {
		result.ContentLength = resp.ContentLength
		result.StatusCode = resp.StatusCode
		//io.Copy(ioutil.Discard, resp.Body)
		result.Body,err=ioutil.ReadAll(resp.Body)
	}
	if err != nil {
		result.Err = err
	}
	t := time.Now()
	result.ResDuration = t.Sub(resStart)
	result.Duration = t.Sub(s)
	return result
}
