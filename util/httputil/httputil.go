package httputil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	loggerFactory "shandianyu-minisdk-mailer/provider/logger"
	"time"
)

// 发送 post 请求
func Post(url string, headers map[string]string, body []byte, target interface{}) error {
	var logger = loggerFactory.GetLogger()
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		logger.Error("http.NewRequest error: %v", err)
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("Content-Type", "application/json")
	logger.Info("send POST request to %s", url)
	logger.Info("request data: %s", string(body))
	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("send http post request error: %v", err)
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(target)
	b, _ := json.Marshal(target)
	logger.Info("response data: %s", string(b))
	if err != nil {
		return err
	}
	return nil
}

// 发送 get 请求
func Get(url string, headers map[string]string) []byte {
	var logger = loggerFactory.GetLogger()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("http.NewRequest error: %v", err)
		return make([]byte, 0)
	}

	// 设置请求头
	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	req.Header.Set("Content-Type", "application/json")

	logger.Info("send GET request to %s", url)
	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("send http get request error: %v+", err)
		return make([]byte, 0)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read data error: %v+", err)
		return make([]byte, 0)
	}
	return bodyBytes
}
