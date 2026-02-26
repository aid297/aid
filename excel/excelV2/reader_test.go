package excelV2

import (
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestReader1(t *testing.T) {
	reader := NewReader().
		SetFilename(Filename("./test-by-filename.xlsx")).
		SetOpenFile(UnzipXMLSizeLimit(10*1024*1024), UnzipSizeLimit(10<<30)).
		Read(
			"Sheet B",
			func(rowNum int, rows *excelize.Rows) (err error) {
				var cols []string
				if cols, err = rows.Columns(); err != nil {
					return
				}

				for colNum := range cols {
					value := cols[colNum]
					t.Logf("列：%v\t", value)
				}

				return
			},
			OriginalRow(5),
			FinishedRow(10),
		)
	_ = reader.GetRawExcel()
}
