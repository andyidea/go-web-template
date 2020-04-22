package util

import (
	"bufio"
	"bytes"

	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
)

//GenXLSX 生成单一sheet的xlsx文件
//  sheetName 表格名称
//  header 标题行(可选)
//	headerLen 列宽(可选) 0表示默认，默认值大约为8
//  data xlsx数据，二维数组
func GenXLSX(sheetName string, header []string, headerLen []float64, data [][]string) ([]byte, error) {
	var err error
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	if data == nil {
		return nil, errors.New("data 不可以为nil")
	}
	//data二维切片中 中如果存在值为nil的一维切片元素，那就挂了吧

	//创建Excel文件
	file = xlsx.NewFile()
	sheet, err = file.AddSheet(sheetName)
	if err != nil {
		return nil, errors.New("创建Excel sheet失败")
	}

	//写入标题行
	if header != nil {
		row = sheet.AddRow()
		for _, v := range header {
			cell = row.AddCell()
			cell.GetStyle().Fill.PatternType = "solid"
			cell.GetStyle().Fill.FgColor = "00ABABAB"
			cell.GetStyle().Alignment.Horizontal = "center"
			cell.Value = v
		}
	}

	//设置列宽
	for i := 0; headerLen != nil && i < len(header) && i < len(headerLen); i++ {
		sheet.Col(i).Width = headerLen[i]
	}

	//写入数据
	for _, r := range data {
		row = sheet.AddRow()
		for _, v := range r {
			cell = row.AddCell()
			cell.Value = v
		}
	}

	//写入buffer
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	err = file.Write(writer)
	if err != nil {
		return nil, errors.New("Excel文件写入buffer失败")
	}
	err = writer.Flush()
	if err != nil {
		return nil, errors.New("Buffer flush失败")
	}

	return buf.Bytes(), nil
}
