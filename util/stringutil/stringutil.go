package stringutil

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// 首字母大写
func FirstUpperCase(str string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(str[:1]), str[1:])
}

// 首字母小写
func FirstLowerCase(str string) string {
	return fmt.Sprintf("%s%s", strings.ToLower(str[:1]), str[1:])
}

// 字符串模板替换
func TemplateParse(str string, data map[string]string) string {
	t, _ := template.New("example").Parse(str)
	var buf bytes.Buffer
	t.Execute(&buf, data)
	return buf.String()
}
