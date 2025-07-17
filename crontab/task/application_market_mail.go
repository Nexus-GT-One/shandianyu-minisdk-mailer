package task

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	htmlcharset "golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"mime"
	"regexp"
	"shandianyu-minisdk-mailer/mail_parser"
	"shandianyu-minisdk-mailer/service"
	"shandianyu-minisdk-mailer/thirdparty/feishu"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"shandianyu-minisdk-mailer/util/secretutil"
	"shandianyu-minisdk-mailer/util/stringutil"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

func init() {
	message.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return htmlcharset.NewReaderLabel(charset, input) // 使用 x/net/html/charset 提供的解码器
	}

	for {
		run()
		time.Sleep(1 * time.Minute)
	}
}

// 支持 GBK 和 UTF-8 的 CharsetReader
func CharsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(charset) {
	case "gbk", "gb2312", "gb18030":
		return transform.NewReader(input, simplifiedchinese.GBK.NewDecoder()), nil
	case "utf-8", "us-ascii":
		return input, nil
	default:
		// 默认不转换，直接返回，可能会乱码
		return input, nil
	}
}

// 解码邮件头，支持GBK编码等
func decodeMIMEHeader(header string) string {
	dec := &mime.WordDecoder{
		CharsetReader: func(charset string, input io.Reader) (io.Reader, error) {
			return htmlcharset.NewReaderLabel(charset, input)
		},
	}
	decoded, err := dec.DecodeHeader(header)
	if err != nil {
		return header // fallback
	}
	return decoded
}

// 解码正文，自动尝试 UTF-8 → GBK fallback
func decodeBodyWithFallback(body io.Reader) string {
	raw, err := ioutil.ReadAll(body)
	if err != nil {
		return "[读取失败]"
	}
	if utf8.Valid(raw) {
		return string(raw)
	}
	decoded, err := ioutil.ReadAll(transform.NewReader(
		bytes.NewReader(raw),
		simplifiedchinese.GBK.NewDecoder(),
	))
	if err == nil {
		return string(decoded)
	}
	return string(raw)
}

// extractPlainTextFromEntity 统一提取正文，自动 fallback 到 text/html，并识别编码
func extractPlainTextFromEntity(entity *message.Entity) (string, error) {
	if entity == nil {
		return "", fmt.Errorf("空邮件体")
	}

	// multipart 解析
	if mr := entity.MultipartReader(); mr != nil {
		var htmlFallback string
		var htmlErr error

		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", fmt.Errorf("读取邮件Part失败: %w", err)
			}

			contentType := "text/html"
			_, params, _ := part.Header.ContentType()
			charset := strings.ToLower(params["charset"])
			switch contentType {
			case "text/plain":
				return readBodyWithCharset(part.Body, charset)
			case "text/html":
				htmlBody, err := readBodyWithCharset(part.Body, charset)
				if err == nil {
					htmlFallback = htmlBody
				} else {
					htmlErr = err
				}
			}
		}

		// fallback 到 HTML 内容
		if htmlFallback != "" {
			return htmlFallback, nil
		} else if htmlErr != nil {
			return "", htmlErr
		}

		return "", fmt.Errorf("未找到可用的 text/plain 或 text/html")
	}

	// 非 multipart 的处理
	contentType, params, err := entity.Header.ContentType()
	if err != nil {
		contentType = "text/plain"
	}
	charset := strings.ToLower(params["charset"])

	switch contentType {
	case "text/plain":
		return readBodyWithCharset(entity.Body, charset)
	case "text/html":
		htmlBody, err := readBodyWithCharset(entity.Body, charset)
		if err != nil {
			return "", err
		}
		return "[HTML邮件内容]\n" + htmlBody, nil
	default:
		return "", fmt.Errorf("未知内容类型: %s", contentType)
	}
}

// 根据 charset 解码正文
func readBodyWithCharset(r io.Reader, charset string) (string, error) {
	return removeHTMLAndCSS(decodeBodyWithFallback(r)), nil
}

func removeHTMLAndCSS(input string) string {
	// 移除 <style> 标签及其内容
	reStyle := regexp.MustCompile(`(?is)<style[^>]*?>.*?</style>`)
	input = reStyle.ReplaceAllString(input, "")

	// 移除内联样式 style="..."
	reInlineStyle := regexp.MustCompile(`(?i)style\s*=\s*"(.*?)"`)
	input = reInlineStyle.ReplaceAllString(input, "")

	// 移除所有 HTML 标签
	reTags := regexp.MustCompile(`(?is)<[^>]+>`)
	input = reTags.ReplaceAllString(input, "")

	// 合并多个换行符和空白行为单个换行
	reNewlines := regexp.MustCompile(`[\s\r\n]{2,}`)
	input = reNewlines.ReplaceAllString(input, "\n")

	// 去除前后空白
	return strings.TrimSpace(input)
}

func run() {
	// 新浪邮箱
	serverName := "imap.sina.com"
	port := 993
	username := "fishflash_mail@sina.com"
	password := "f3da51a030dfb484"

	// 连接 IMAP 服务器（TLS）
	tlsConfig := &tls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS12,
		CipherSuites:       []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, tls.TLS_RSA_WITH_AES_128_GCM_SHA256},
	}

	c, err := client.DialTLS(fmt.Sprintf("%s:%d", serverName, port), tlsConfig)
	if err != nil {
		logger.Error("连接IMAP服务器失败: %v", err)
		return
	}
	defer c.Logout()

	if err := c.Login(username, password); err != nil {
		logger.Error("登录失败: %v", err)
		return
	}
	logger.Info("登录成功")

	// 选择收件箱
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		logger.Error("选择收件箱失败: %v", err)
		return
	}
	logger.Info("收件箱邮件数: %d\n", mbox.Messages)
	if mbox.Messages == 0 {
		return
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(1, mbox.Messages)
	criteria := []imap.FetchItem{imap.FetchEnvelope, imap.FetchInternalDate}
	messages := make(chan *imap.Message, mbox.Messages)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, criteria, messages)
	}()

	if err := <-done; err != nil {
		logger.Error("Fetch邮件元数据失败: %v", err)
		return
	}

	mails := make([]*imap.Message, 0)
	for msg := range messages {
		mails = append(mails, msg)
	}

	// 获取最新游标
	lastMailIndex := service.SystemService.GetLastMailIndex()

	// 有游标的时候倒序排序；没游标的时候就正序排序
	if lastMailIndex > 0 {
		sort.Slice(mails, func(i, j int) bool {
			return mails[i].InternalDate.After(mails[j].InternalDate)
		})
	} else {
		sort.Slice(mails, func(i, j int) bool {
			return mails[j].InternalDate.After(mails[i].InternalDate)
		})
	}

	pageSize := 50
	beijingLoc, _ := time.LoadLocation("Asia/Shanghai")
	total := len(mails)
	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	for page := 1; page <= totalPages; page++ {
		logger.Info("==== 第 %d 页 / 共 %d 页 ====\n", page, totalPages)
		start := (page - 1) * pageSize
		end := start + pageSize
		if end > total {
			end = total
		}
		for _, mailInfo := range mails[start:end] {
			seqset := new(imap.SeqSet)
			seqset.AddNum(mailInfo.SeqNum)

			section := &imap.BodySectionName{}
			msgs := make(chan *imap.Message, 1)
			done := make(chan error, 1)
			go func() {
				done <- c.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, msgs)
			}()

			msg := <-msgs
			if msg == nil {
				logger.Error("邮件 %d 拉取失败\n", mailInfo.SeqNum)
				continue
			}

			if err := <-done; err != nil {
				logger.Error("邮件 %d 内容读取失败: %v\n", mailInfo.SeqNum, err)
				continue
			}

			// 设置已读
			item := imap.FormatFlagsOp(imap.AddFlags, false)
			flags := []interface{}{imap.SeenFlag}
			c.Store(seqset, item, flags, nil)

			r := msg.GetBody(section)
			if r == nil {
				logger.Error("邮件 %d 内容为空\n", mailInfo.SeqNum)
				continue
			}

			entity, err := message.Read(r)
			if err != nil {
				logger.Error("邮件 %d MIME 解析失败: %v\n", mailInfo.SeqNum, err)
				continue
			}
			subject := decodeMIMEHeader(mailInfo.Envelope.Subject)
			dateInBeijing := mailInfo.InternalDate.In(beijingLoc)
			bodyText, err := extractPlainTextFromEntity(entity)
			if err != nil {
				bodyText = "正文解析失败: " + err.Error()
			}

			// 如果已经读取过的邮件就不需要再读取了
			if lastMailIndex >= dateInBeijing.UnixMilli() {
				logger.Info("邮件已经读完了，最新游标：%v", dateInBeijing.Format(time.DateTime))
				return
			}

			ms5String := secretutil.MD5(bodyText)
			existsGameMail := service.GameMailService.FindByMd5(ms5String)
			if existsGameMail != nil {
				return
			}

			// 打印一下邮件
			logger.Info("序号: %d\n日期: %s\n收件人: %s\n主题: %s\nmd5: %s\n正文:\n%s\n\n",
				mailInfo.SeqNum,
				dateInBeijing.Format(time.DateTime),
				strings.TrimSpace((*arrayutil.First(mailInfo.Envelope.To)).Address()),
				subject,
				ms5String,
				bodyText)

			// 解析邮件
			oneGame, newGameMail := mail_parser.ParseMail(subject, dateInBeijing.Format(time.DateTime), bodyText)

			// 记录读取游标
			service.SystemService.SaveLastMailIndex(dateInBeijing.UnixMilli())

			// 如果不是苹果提审相关的游戏就跳过
			if oneGame == nil || newGameMail == nil {
				continue
			}

			// 邮件写入数据库
			newGameMail.Developer = strings.TrimSpace((*arrayutil.First(mailInfo.Envelope.To)).Address())
			gameDb.InsertOne(*newGameMail)

			// 发送消息
			title := fmt.Sprintf("游戏 %s 的邮件消息", oneGame.Symbol)
			data := map[string]string{
				"status":     newGameMail.Status,
				"developer":  newGameMail.Developer,
				"appVersion": newGameMail.AppVersion,
				"time":       dateInBeijing.Format(time.DateTime),
				"title":      newGameMail.Title,
				"content":    newGameMail.Content,
			}
			content := stringutil.TemplateParse(`**游戏状态**：{{.status}}
**游戏版本**：{{.appVersion}}
**收件时间**：{{.time}}
**邮件标题**：{{.title}}
**邮件内容**：
{{.content}}`, data)
			feishu.MailRobot().SendRobotInteractive(title, content)
		}
	}
}
