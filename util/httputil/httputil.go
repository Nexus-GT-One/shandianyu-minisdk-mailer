package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	loggerFactory "shandianyu-minisdk-mailer/provider/logger"
	"time"
)

var logger = loggerFactory.GetLogger()

// 发送 post 请求
func Post(url string, headers map[string]string, body []byte, target interface{}) error {
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

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read data error: %v", err)
		return err
	}

	logger.Info("send http POST request %s\nbody: %s\nresponse: %s", url, string(body), string(b))

	err = json.Unmarshal(b, &target)
	if err != nil {
		logger.Error("send http post response: %v", string(b))
		logger.Error("send http post response decode error: %v", err)
		return errors.Wrap(err, fmt.Sprintf("response %s", string(b)))
	}
	return nil
}
