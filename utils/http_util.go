package utils

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

func HttpGet(urlStr string, timeoutMS int) (data []byte, err error) {
	return HttpGetWithReferer(urlStr, "", timeoutMS)
}

func HttpGetWithRefererTryN(urlStr, referer string, timeoutMS int, n int) (data []byte, err error) {
	for i := 0; i < n; i++ {
		data, err = HttpGetWithReferer(urlStr, referer, timeoutMS)
		if err == nil {
			return
		}
		if !strings.Contains(err.Error(), "timeout") {
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
	return

}
func HttpGetWithReferer(urlStr, referer string, timeoutMS int) (data []byte, err error) {
	now := time.Now()
	client := HttpWithTimeOut(now, timeoutMS)
	req, err := http.NewRequest("GET", urlStr, nil)
	// User-Agent:

	if err != nil {
		return nil, err
	}
	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:39.0) Gecko/20100101 Firefox/39.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")

	response, err := client.Do(req)

	if err != nil {
		return
	}

	defer response.Body.Close()
	reader, err := charset.NewReader(response.Body, response.Header.Get("Content-Type")) // 处理字符集的问题,自动进行转码
	if err != nil {
		data, err = ioutil.ReadAll(response.Body)
	} else {
		data, err = ioutil.ReadAll(reader)
	}
	return
}
func HttpWithTimeOut(now time.Time, timeoutMillSeconds int) http.Client {
	timeoutDur := time.Millisecond * time.Duration(timeoutMillSeconds)
	// 在拨号回调中，使用DialTimeout来支持连接超时，当连接成功后，利用SetDeadline来让连接支持读写超时。
	fun := func(network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, timeoutDur)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(now.Add(timeoutDur))
		return conn, nil
	}
	transport := &http.Transport{Dial: fun, ResponseHeaderTimeout: timeoutDur}

	client := http.Client{
		Transport: transport,
	}
	return client
}
