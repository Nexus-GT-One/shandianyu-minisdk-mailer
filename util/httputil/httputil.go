package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	loggerFactory "shandianyu-minisdk-mailer/provider/logger"
	"strconv"
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
	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("send http post request error: %v", err)
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(target)
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
		b, _ := json.Marshal(map[string]any{"code": 1, "message": fmt.Sprintf("http.NewRequest error: %v", err)})
		return b
	}

	// 设置请求头
	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("send http get request error: %v+", err)
		b, _ := json.Marshal(map[string]any{"code": 1, "message": err.Error()})
		return b
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read data error: %v", err)
		b, _ = json.Marshal(map[string]any{"code": 1, "message": fmt.Sprintf("read data error: %v", err)})
		return b
	}

	return b
}

// 下载文件
//
// fileURL-文件下载地址; fileDist-文件本地保存地址
func Download(fileURL, fileDist string) error {
	// 发起 HTTP 请求
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 获取文件总大小
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	// 创建进度条
	bar := progressbar.DefaultBytes(int64(size), "downloading ...")

	// 创建文件
	file, err := os.Create(fileDist)
	if err != nil {
		return err
	}
	defer file.Close()

	// 读取数据并写入文件，同时更新进度条
	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	return err
}
