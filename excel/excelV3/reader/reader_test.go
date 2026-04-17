package reader

import (
	`testing`

	`github.com/xuri/excelize/v2`
)

func Test(t *testing.T) {
	data := make([][]string, 0)

	reader := NewReader(
		Filename("./2月.xlsx"),
		OriginalRow(4),
		FinishedRow(4),
	).Read("月度汇总", func(rowNum int, rows *excelize.Rows) (err error) {
		var cols []string
		if cols, err = rows.Columns(); err != nil {
			return err
		}

		data = append(data, cols[6:])

		return
	})

	if err := reader.GetError(); err != nil {
		t.Errorf("读取文件失败：%v", err)
	}

	t.Logf("读取成功：%+v", data)
}
