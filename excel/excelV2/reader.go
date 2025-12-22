package excelV2

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
)

type (
	Reader struct {
		Error          error
		lock           sync.RWMutex
		rawFile        *excelize.File
		filename       string
		sheetName      string
		originalRow    uint
		finishedRow    uint
		originalColNo  uint
		originalColTxt string
		finishedColNo  uint
		finishedColTxt string
		originalData   [][]string
		dataMap        map[string]string
		dataMaps       []map[string]string
	}
)

func (*Reader) New(attrs ...ReaderAttributer) *Reader {
	return (&Reader{lock: sync.RWMutex{}, sheetName: "Sheet 1", originalRow: 1, originalColNo: 1, originalData: make([][]string, 0), dataMap: make(map[string]string), dataMaps: make([]map[string]string, 0)}).setAttrs(attrs...)
}

func (my *Reader) setAttrs(attrs ...ReaderAttributer) *Reader {
	for i := range attrs {
		attrs[i].Register(my)
	}
	return my
}

func (my *Reader) SetAttrs(attrs ...ReaderAttributer) *Reader {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.setAttrs(attrs...)
}

func (my *Reader) read() *Reader {
	var (
		err  error
		rows = make([][]string, 0)
	)

	if my.filename == "" {
		my.Error = ErrFilenameRequired
		return my
	}

	if my.sheetName == "" {
		my.Error = ErrSheetNameRequired
		return my
	}

	if my.rawFile, err = excelize.OpenFile(fmt.Sprintf(my.filename)); err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrOpen, err)
		return my
	}

	defer func(r *Reader) {
		if err = r.rawFile.Close(); err != nil {
			r.Error = fmt.Errorf("%w：%w", ErrClose, err)
		}
	}(my)

	if rows, err = my.rawFile.GetRows(my.sheetName); err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrRead, err)
		return my
	}

	if len(rows) > 0 {
		if my.finishedRow == 0 {
			for rowNo := range rows[my.originalRow:] {
				if my.finishedColNo == 0 {
					my.originalData[rowNo] = rows[rowNo][my.originalColNo:]
				} else {
					my.originalData[rowNo] = rows[rowNo][my.originalColNo:my.finishedColNo]
				}
			}
		} else {
			for rowNo := range rows[my.originalRow:my.finishedRow] {
				if my.finishedColNo == 0 {
					my.originalData[rowNo] = rows[rowNo][my.originalColNo:]
				} else {
					my.originalData[rowNo] = rows[rowNo][my.originalColNo:my.finishedColNo]
				}
			}
		}
	}

	return my
}

func (my *Reader) getOriginalData() [][]string {
	if len(my.originalData) == 0 && my.Error == nil {
		my.read()
	}

	return my.originalData
}

func (my *Reader) GetOriginalData() [][]string {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.getOriginalData()
}

func (my *Reader) GetMaps() []map[string]string {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.getOriginalData()

	for row := range my.originalData {
		my.dataMaps[row] = make(map[string]string)
		for col := range my.originalData[row] {
			var colTxt string
			if colTxt, my.Error = ColumnNumberToText(col); my.Error != nil {
				my.Error = fmt.Errorf("%w：行 %d 列索引 %d 转换为列名称错误", ErrColumnNumber, row+1, col+1)
				return nil
			}

			my.dataMaps[row][colTxt] = my.originalData[row][col]
		}
	}

	return my.dataMaps
}

func (my *Reader) GetMap() map[string]string {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.getOriginalData()

	for row := range my.originalData {
		for col := range my.originalData[row] {
			var colTxt string
			if colTxt, my.Error = ColumnNumberToText(col); my.Error != nil {
				my.Error = fmt.Errorf("%w：行 %d 列索引 %d 转换为列名称错误", ErrColumnNumber, row+1, col+1)
				return nil
			}

			my.dataMap[fmt.Sprintf("%s%d", colTxt, row)] = my.originalData[row][col]
		}
	}

	return my.dataMap
}

func (my *Reader) ToStruct(title []string, ret any) {
	my.lock.Lock()
	defer my.lock.Unlock()

	// 确保已读取原始数据
	my.getOriginalData()
	if my.Error != nil {
		return
	}

	if ret == nil {
		my.Error = fmt.Errorf("ret 不能为空")
		return
	}

	rv := reflect.ValueOf(ret)
	if rv.Kind() != reflect.Ptr {
		my.Error = fmt.Errorf("ret 必须是指向切片的指针")
		return
	}

	sliceVal := rv.Elem()
	if sliceVal.Kind() != reflect.Slice {
		my.Error = fmt.Errorf("ret 必须是指向切片的指针")
		return
	}

	elemType := sliceVal.Type().Elem()
	isElemPtr := false
	structType := elemType
	if elemType.Kind() == reflect.Ptr {
		isElemPtr = true
		structType = elemType.Elem()
	}

	if structType.Kind() != reflect.Struct {
		my.Error = fmt.Errorf("切片元素类型必须是结构体或结构体指针")
		return
	}

	// 构建字段名映射（小写），支持字段名和 `json` tag
	fieldIndex := make(map[string]int)
	for i := 0; i < structType.NumField(); i++ {
		f := structType.Field(i)
		// 跳过未导出字段
		if f.PkgPath != "" {
			continue
		}

		// 优先使用 `excel` tag，其次 `json` tag，最后字段名
		// 支持 tag 值为 "-" 表示忽略该字段
		added := false

		if etag := f.Tag.Get("excel"); etag != "" {
			if etag == "-" {
				continue
			}
			parts := strings.Split(etag, ",")
			if parts[0] != "" {
				fieldIndex[strings.ToLower(parts[0])] = i
				added = true
			}
		}

		// 无论是否存在 excel tag，都尝试注册 json tag（允许同时使用两种 tag）
		if jtag := f.Tag.Get("json"); jtag != "" {
			if jtag == "-" {
				continue
			}
			parts := strings.Split(jtag, ",")
			if parts[0] != "" {
				fieldIndex[strings.ToLower(parts[0])] = i
				added = true
			}
		}

		if !added {
			name := strings.ToLower(f.Name)
			fieldIndex[name] = i
		}
	}

	// 遍历每一行并进行转换
	for rowIdx := range my.originalData {
		row := my.originalData[rowIdx]

		// 创建元素实例（总是先创建指针，然后根据需要 append 值或指针）
		newPtr := reflect.New(structType) // *T
		structVal := newPtr.Elem()

		for colIdx, colName := range title {
			if colIdx >= len(row) {
				continue
			}
			raw := row[colIdx]
			key := strings.ToLower(strings.TrimSpace(colName))
			if key == "" {
				continue
			}

			fi, ok := fieldIndex[key]
			if !ok {
				// 未找到字段，跳过
				continue
			}

			f := structVal.Field(fi)
			if !f.CanSet() {
				continue
			}

			// 处理指针字段
			targetType := f.Type()
			isPtrField := false
			if targetType.Kind() == reflect.Ptr {
				isPtrField = true
				targetType = targetType.Elem()
			}

			// 转换字符串到目标类型
			converted, err := convertStringToReflectValue(raw, targetType)
			if err != nil {
				my.Error = fmt.Errorf("第 %d 行列 %s 转换失败：%w", rowIdx+1, colName, err)
				return
			}

			if isPtrField {
				ptr := reflect.New(targetType)
				ptr.Elem().Set(converted)
				f.Set(ptr)
			} else {
				f.Set(converted)
			}
		}

		// append 到切片
		if isElemPtr {
			sliceVal.Set(reflect.Append(sliceVal, newPtr))
		} else {
			sliceVal.Set(reflect.Append(sliceVal, newPtr.Elem()))
		}
	}
}

// convertStringToReflectValue 将字符串转换为给定的 reflect.Kind 的值(非指针)
func convertStringToReflectValue(s string, t reflect.Type) (reflect.Value, error) {
	// 空字符串返回零值
	if s == "" {
		return reflect.Zero(t), nil
	}

	switch t.Kind() {
	case reflect.String:
		return reflect.ValueOf(s).Convert(t), nil
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(b).Convert(t), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, t.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i).Convert(t), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u, err := strconv.ParseUint(s, 10, t.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(u).Convert(t), nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, t.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(f).Convert(t), nil
	case reflect.Struct:
		// 支持 time.Time? 如果需要可以扩展
		return reflect.Zero(t), fmt.Errorf("不支持将字符串直接转换为 struct 类型 (%s)", t.String())
	default:
		return reflect.Zero(t), fmt.Errorf("不支持的目标类型: %s", t.Kind().String())
	}
}
