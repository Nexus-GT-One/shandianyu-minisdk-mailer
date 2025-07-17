package excelutil

import (
	"bytes"
	"embed"
	"github.com/xuri/excelize/v2"
)

// 从embed文件流中读取excel
func ReadFromEmbedFS(file embed.FS, fileName string) [][]string {
	data, _ := file.ReadFile(fileName)
	reader := bytes.NewReader(data)
	f, _ := excelize.OpenReader(reader)
	content, _ := f.GetRows(f.GetSheetName(0))
	return content
}
