package tools

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	HttpClient = http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConnsPerHost: 50000,
			MaxIdleConns:        50000,
			IdleConnTimeout:     time.Second * 10,
		},
	}
)

func HttpPost(url string, body io.Reader, headers map[string]string) ([]byte, error) {
	return request(http.MethodPost, url, body, headers)
}
func HttpGet(url string, body io.Reader, headers map[string]string) ([]byte, error) {
	return request(http.MethodGet, url, body, headers)
}

// method: "POST", "GET", "PUT", "PATCH"
func request(method, url string, body io.Reader, header map[string]string) (bodyByte []byte, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// 设置头
	for k, v := range header {
		req.Header.Set(k, v)
	}
	// 记录request错误数量
	defer func() {
		if err != nil {
			atomic.AddInt64(&ErrCount, 1)
		}
	}()
	now := time.Now()
	resp, err := HttpClient.Do(req)
	// 记录延迟信息
	SetLatency(now)
	// 记录成功数
	atomic.AddInt64(&RequestCount, 1)
	if err != nil {
		return nil, err
	}
	// 记录错误码
	errCodes.Store(resp.StatusCode, 1)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("httpPost failed, status code: %v", resp.StatusCode)
		return nil, err
	}
	bodyByte, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return bodyByte, nil
}
