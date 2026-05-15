package excelReader_test

import (
	"testing"

	"github.com/aid297/aid/v2/excel/excelV3/excelReader"
)

func Test(t *testing.T) {
	data := make([][]string, 0)

	reader := excelReader.NewReader(
		excelReader.Filename("./2月.xlsx"),
		excelReader.OriginalColumn(6),
		excelReader.OriginalRow(4),
		excelReader.FinishedRow(4),
	).Read("月度汇总", func(rowNum int, cols []string) (err error) {
		data = append(data, cols)

		return
	})

	if err := reader.GetError(); err != nil {
		t.Errorf("读取文件失败：%v", err)
	}

	t.Logf("读取成功：%+v", data)
}
