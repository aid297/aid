package v1HTTPService

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dromara/carbon/v2"

	"github.com/aid297/aid/v2/anySlice"
	"github.com/aid297/aid/v2/excel/excelV3/excelReader"
	"github.com/aid297/aid/v2/operation"
	"github.com/aid297/aid/v2/web-site/backend/aid-web-backend/src/module/httpModule/v1HTTPModule/request"
)

type (
	CheckingInService struct{}

	StandardDateKind string
	StandardDate     struct {
		Kind         StandardDateKind `json:"kind"`         // 日期类型
		OriginalText string           `json:"originalText"` // 原始内容
		Value        carbon.Carbon    `json:"value"`        // 日期
		Text         string           `json:"text"`         // 文字内容
	}
	StandardDateRet struct {
		OriginalTexts   []string                `json:"originalTexts"`
		StandardDateMap map[string]StandardDate `json:"standardDateMap"`
	}

	EverydayRet struct {
		ClockInTime1 string `json:"clockInTime1"` // 打卡时间1
		ClockInNote1 string `json:"clockInNote1"` // 打卡说明1
		ClockInOK1   bool   `json:"clockInOK1"`   // 打卡1是否正常
		ClockInTime2 string `json:"clockInTime2"` // 打卡时间2
		ClockInNote2 string `json:"clockInNote2"` // 打卡说明2
		ClockInOK2   bool   `json:"clockInOK2"`   // 打卡2是否正常
		EventNote    string `json:"eventNote"`    // 事件说明
	}
)

const (
	WORKDAY            = "WORKDAY"            // 工作日
	FORCE_WORKDAY      = "FORCE-WORKDAY"      // 调休
	FORCE_ANNUAL_LEAVE = "FORCE-ANNUAL-LEAVE" // 强制年假休息
	WEEKEND            = "WEEKEND"            // 周末
	HOLIDAY2           = "HOLIDAY2"           // 双薪
	HOLIDAY3           = "HOLIDAY3"           // 三薪
)

// Cal 计算考勤
func (my *CheckingInService) Cal(
	r excelReader.Reader,
	form *request.CheckingInCalRequest,
) (
	standardDate *StandardDateRet,
	everyday map[string]map[string]EverydayRet,
	monthly map[string]map[string]string,
	err error,
) {
	if standardDate, err = my.getDateHeader(r, form); err != nil {
		return
	}

	if everyday, err = getEveryday(r); err != nil {
		return
	}

	if monthly, err = getMonthly(r, form); err != nil {
		return
	}

	return
}

// getDateHeader 获取日期头
func (*CheckingInService) getDateHeader(r excelReader.Reader, form *request.CheckingInCalRequest) (*StandardDateRet, error) {
	var (
		err          error
		dates        = &StandardDateRet{}
		originalDate = carbon.Parse(form.OriginalDate)
	)

	if err = r.Read(
		"打卡时间",
		func(_ int, cols []string) (err error) {
			standardDateMap := make(map[string]StandardDate, len(cols))
			dates.OriginalTexts = make([]string, len(cols))

			for idx := range len(cols) {
				standardDateValue := *originalDate.AddDays(idx)
				standardDate := StandardDate{
					Kind:         getStandardDateKind(standardDateValue, form.ForceWorkdays, form.ForceAnnualLeaves, form.Holidays2, form.Holidays3),
					OriginalText: cols[idx],
					Value:        standardDateValue,
					Text:         standardDateValue.ToDateString(),
				}

				dates.OriginalTexts[idx] = standardDate.Text
				standardDateMap[standardDate.Text] = standardDate
			}

			dates.StandardDateMap = standardDateMap
			return
		},
		excelReader.OriginalRow(4),
		excelReader.FinishedRow(4),
		excelReader.OriginalColumn(6),
	).GetError(); err != nil {
		return nil, err
	}

	if len(dates.StandardDateMap) == 0 {
		return nil, errors.New("日期头不存在")
	}

	return dates, nil
}

// dateStandardDateKind 获取日期类型
func getStandardDateKind(
	c carbon.Carbon,
	forceWorkdays []string,
	forceAnnualLeave []string,
	holidays2 []string,
	holidays3 []string,
) StandardDateKind {
	switch {
	case anySlice.NewList(forceWorkdays).In(c.ToDateString()): // 调休
		return FORCE_WORKDAY
	case anySlice.NewList(forceAnnualLeave).In(c.ToDateString()): // 强制年假休休
		return FORCE_ANNUAL_LEAVE
	case anySlice.NewList(holidays2).In(c.ToDateString()): // 双薪
		return HOLIDAY2
	case anySlice.NewList(holidays3).In(c.ToDateString()): // 三薪
		return HOLIDAY3
	case c.IsWeekend(): // 周末
		return WEEKEND
	default: // 工作日
		return WORKDAY
	}
}

// getEveryday 获取每日打卡记录
func getEveryday(r excelReader.Reader) (map[string]map[string]EverydayRet, error) {
	data := make(map[string]map[string]EverydayRet)

	if err := r.Read(
		"每日统计",
		func(_ int, cols []string) (err error) {
			name := strings.ReplaceAll(cols[0], " ", "")
			date := getDate(cols[6])

			if _, ok := data[date]; !ok {
				data[date] = make(map[string]EverydayRet)
			}

			ret := EverydayRet{
				ClockInTime1: operation.NewTernary(operation.TrueFn(func() string { return cols[9] })).GetByValue(len(cols) >= 9),
				ClockInNote1: operation.NewTernary(operation.TrueFn(func() string { return cols[10] })).GetByValue(len(cols) >= 10),
				ClockInTime2: operation.NewTernary(operation.TrueFn(func() string { return cols[11] })).GetByValue(len(cols) >= 11),
				ClockInNote2: operation.NewTernary(operation.TrueFn(func() string { return cols[12] })).GetByValue(len(cols) >= 12),
				EventNote:    operation.NewTernary(operation.TrueFn(func() string { return cols[21] })).GetByValue(len(cols) >= 21),
			}

			ret.ClockInOK1 = carbon.Parse(fmt.Sprintf("%s %s", date, ret.ClockInTime1)).IsAM()
			if strings.Contains(ret.ClockInTime2, "次日 ") { // 次日打卡
				clockInTime2 := strings.Replace(ret.ClockInTime2, "次日 ", "", 1)
				ret.ClockInTime2 = carbon.Parse(date+" ", clockInTime2).AddDay().ToDateString()
				ret.ClockInOK2 = true
			} else { // 当日打卡
				ret.ClockInOK2 = carbon.Parse(fmt.Sprintf("%s %s", date, ret.ClockInTime2)).IsPM()
			}

			data[date][name] = ret

			return
		},
		excelReader.OriginalRow(5),
	).GetError(); err != nil {
		return nil, err
	}

	return data, nil
}

func getDate(datetime string) string {
	return carbon.ParseByLayout(
		strings.NewReplacer("星期一", "Mon", "星期二", "Tue", "星期三", "Wed", "星期四", "Thu", "星期五", "Fri", "星期六", "Sat", "星期日", "Sun").Replace(datetime),
		"06-01-02 Mon",
	).ToDateString()
}

// getMonthly 获取月度汇总数据
func getMonthly(r excelReader.Reader, form *request.CheckingInCalRequest) (map[string]map[string]string, error) {
	data := make(map[string]map[string]string)

	err := r.Read(
		"月度汇总",
		func(_ int, cols []string) (err error) {
			name := strings.ReplaceAll(cols[0], " ", "")
			for idx := range cols[6:] {
				date := carbon.Parse(form.OriginalDate).AddDays(idx).ToDateString()
				if _, ok := data[date]; !ok {
					data[date] = make(map[string]string)
				}

				data[date][name] = cols[idx+6]
			}

			return
		},
		excelReader.OriginalRow(5),
	).GetError()

	return data, err
}
