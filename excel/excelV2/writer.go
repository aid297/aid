package excelV2

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"

	"github.com/aid297/aid/str"
)

type Writer struct {
	Error      error
	lock       sync.RWMutex
	rawFile    *excelize.File
	filename   string
	sheetName  string
	isSheetSet bool
}

func (*Writer) New(attrs ...WriterAttributer) *Writer {
	return (&Writer{lock: sync.RWMutex{}, rawFile: excelize.NewFile()}).setAttrs(attrs...)
}

func (my *Writer) setAttrs(attrs ...WriterAttributer) *Writer {
	for i := range attrs {
		attrs[i].Register(my)
	}
	return my
}

func (my *Writer) SetAttrs(attrs ...WriterAttributer) *Writer {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.setAttrs(attrs...)
}

func (my *Writer) getFilename() string { return my.filename }
func (my *Writer) GetFilename() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getFilename()
}
func (my *Writer) getSheetName() string { return my.sheetName }
func (my *Writer) GetSheetName() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getSheetName()
}

// setCell 设置 cells 值和格式
func (my *Writer) setCell(cell *Cell) *Writer {
	return my.setSheet().setCellValue(cell).setCellStyle(cell)
}

// setCellValue 设置 cells 值
func (my *Writer) setCellValue(cell *Cell) *Writer {
	var err error

	switch cell.GetContentType() {
	case CellContentTypeFormula:
		if err = my.rawFile.SetCellFormula(my.sheetName, cell.GetCoordinate(), cast.ToString(cell.GetContent())); err != nil {
			my.Error = fmt.Errorf("%w：%s %s %w", ErrWriteCellFormula, cell.GetCoordinate(), cell.GetContent(), err)
			return my
		}
	case CellContentTypeInt:
		if err = my.rawFile.SetCellInt(my.sheetName, cell.GetCoordinate(), cast.ToInt(cell.GetContent())); err != nil {
			my.Error = fmt.Errorf("%w：%s %s %w", ErrWriteCellInt, cell.GetCoordinate(), cell.GetContent(), err)
			return my
		}
	case CellContentTypeFloat:
		if err = my.rawFile.SetCellFloat(my.sheetName, cell.GetCoordinate(), cast.ToFloat64(cell.GetContent()), 2, 64); err != nil {
			my.Error = fmt.Errorf("%w：%s %s %w", ErrWriteCellFloat, cell.GetCoordinate(), cell.GetContent(), err)
			return my
		}
	case CellContentTypeBool:
		if err = my.rawFile.SetCellBool(my.sheetName, cell.GetCoordinate(), cast.ToBool(cell.GetContent())); err != nil {
			my.Error = fmt.Errorf("%w：%s %s %w", ErrWriteCellBool, cell.GetCoordinate(), cell.GetContent(), err)
			return my
		}
	case CellContentTypeTime:
		if err = my.rawFile.SetCellValue(my.sheetName, cell.GetCoordinate(), cast.ToTime(cell.GetContent())); err != nil {
			my.Error = fmt.Errorf("%w：%s %s %w", ErrWriteCellTime, cell.GetCoordinate(), cell.GetContent(), err)
			return my
		}
	default:
		if err = my.rawFile.SetCellValue(my.sheetName, cell.GetCoordinate(), cell.GetContent()); err != nil {
			my.Error = fmt.Errorf("%w：%s %s %w", ErrWriteCellAny, cell.GetCoordinate(), cell.GetContent(), err)
			return my
		}
	}

	return my
}

// setCellStyle 设置 cells 格式
func (my *Writer) setCellStyle(cell *Cell) *Writer {
	fill := excelize.Fill{Type: "pattern", Pattern: 0, Color: []string{""}}
	if cell.GetPatternRGB() != "" {
		fill.Pattern = 1
		fill.Color[0] = cell.GetPatternRGB()
	}

	var borders = make([]excelize.Border, 0)
	if cell.GetBorder().LengthNotEmpty() > 0 {
		cell.GetBorder().Each(func(_ int, item border) {
			borders = append(borders, excelize.Border{
				Type:  item.Type,
				Color: item.RGB,
				Style: item.Style,
			})
		})
	}

	if style, err := my.rawFile.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   cell.GetFontBold(),
			Italic: cell.GetFontItalic(),
			Family: cell.GetFontFamily(),
			Size:   cell.GetFontSize(),
			Color:  cell.GetFontRGB(),
		},
		Alignment: &excelize.Alignment{WrapText: cell.GetWrapText()},
		Fill:      fill,
		Border:    borders,
	}); err != nil {
		my.Error = fmt.Errorf("%w：%s", ErrSetFont, cell.GetCoordinate())
	} else {
		my.Error = my.rawFile.SetCellStyle(my.sheetName, cell.GetCoordinate(), cell.GetCoordinate(), style)
	}

	return my
}

func (my *Writer) createSheet(sheetName string) *Writer {
	if sheetName == "" {
		my.Error = ErrSheetNameRequired
		return my
	}
	sheetIndex, err := my.rawFile.NewSheet(sheetName)
	if err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrCreateSheet, err)
		return my
	}

	my.rawFile.SetActiveSheet(sheetIndex)
	my.sheetName = my.rawFile.GetSheetName(sheetIndex)

	return my
}

// CreateSheet 创建工作表
func (my *Writer) CreateSheet(sheetName string) *Writer {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.createSheet(sheetName)
}

func (my *Writer) setSheet() *Writer {
	var (
		sheetIndex int
		err        error
	)

	if my.isSheetSet {
		return my
	}

	if my.sheetName == "" {
		my.Error = ErrSheetNameRequired
		return my
	}

	if sheetIndex, err = my.rawFile.GetSheetIndex(my.sheetName); err != nil {
		my.Error = fmt.Errorf("%w：%w %s", ErrSetSheet, err, my.sheetName)
		return my
	}

	if sheetIndex == -1 {
		// sheet 不存在，创建sheet
		my.createSheet(my.sheetName)
	} else {
		// sheet 存在，设置为活动sheet
		my.rawFile.SetActiveSheet(sheetIndex)
	}

	my.isSheetSet = true
	return my
}

// Save 保存文件
func (my *Writer) Save() *Writer {
	var err error

	my.lock.Lock()

	if my.filename == "" {
		my.Error = ErrFilenameRequired
		return my
	}

	if err = my.rawFile.SaveAs(my.filename); err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrSave, err)
		return my
	}

	my.lock.Unlock()
	return my
}

func (my *Writer) FromStruct(data any, title []string, offset int, attrs ...CellAttributer) *Writer {
	my.lock.Lock()
	defer my.lock.Unlock()
	if defaultOffset := 1; offset <= 0 {
		offset = defaultOffset
	}
	if data == nil {
		my.Error = fmt.Errorf("data 不能为空")
		return my
	}

	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Ptr {
		my.Error = fmt.Errorf("data 必须是指向切片的指针")
		return my
	}

	sv := rv.Elem()
	if sv.Kind() != reflect.Slice {
		my.Error = fmt.Errorf("data 必须是指向切片的指针")
		return my
	}

	elemType := sv.Type().Elem()
	isElemPtr := false
	structType := elemType
	if elemType.Kind() == reflect.Ptr {
		isElemPtr = true
		structType = elemType.Elem()
	}

	if structType.Kind() != reflect.Struct {
		my.Error = fmt.Errorf("切片元素必须为结构体或结构体指针")
		return my
	}

	// 构建字段映射：优先 excel tag, 然后 json tag, 最后字段名（小写）
	fieldIndex := make(map[string]int)
	for i := 0; i < structType.NumField(); i++ {
		f := structType.Field(i)
		if f.PkgPath != "" { // unexported
			continue
		}
		if et := f.Tag.Get("excel"); et != "" {
			if et == "-" {
				continue
			}
			parts := strings.Split(et, ",")
			if parts[0] != "" {
				fieldIndex[strings.ToLower(parts[0])] = i
			}
		}
		if jt := f.Tag.Get("json"); jt != "" {
			if jt == "-" {
				continue
			}
			parts := strings.Split(jt, ",")
			if parts[0] != "" {
				fieldIndex[strings.ToLower(parts[0])] = i
			}
		}
		name := strings.ToLower(f.Name)
		fieldIndex[name] = i
	}

	// 写标题（如果提供）
	// writeCell 使用 rn/cn 为 0-based index
	rn := offset - 1
	if len(title) > 0 {
		for cn, t := range title {
			cell := APP.Cell.New(APP.CellAttr.Content.Set(t), APP.CellAttr.ContentType.Set(CellContentTypeAny))
			writeCell(cell, rn, cn, my)
		}
		rn++
	}

	// 遍历切片元素并写入
	for i := 0; i < sv.Len(); i++ {
		item := sv.Index(i)
		if isElemPtr {
			if item.IsNil() {
				// skip nil pointer element
				continue
			}
			item = item.Elem()
		}

		// 对于每一列，根据 title 找到字段并写入
		for cn, t := range title {
			key := strings.ToLower(strings.TrimSpace(t))
			if key == "" {
				continue
			}
			fi, ok := fieldIndex[key]
			if !ok {
				continue
			}
			f := item.Field(fi)
			sf := structType.Field(fi)
			// 未导出字段或不可设置跳过
			if !f.IsValid() || !f.CanInterface() {
				continue
			}

			// prepare per-field attributes from struct tags
			localAttrs := make([]CellAttributer, 0)
			if v := strings.TrimSpace(sf.Tag.Get("excel-font-size")); v != "" {
				if fv, err := strconv.ParseFloat(v, 64); err == nil {
					localAttrs = append(localAttrs, AttrCellFontSize{}.Set(fv))
				}
			}
			if v := strings.TrimSpace(sf.Tag.Get("excel-font-rgb")); v != "" {
				localAttrs = append(localAttrs, AttrCellFontRGB{}.Set(v))
			}
			if v := strings.TrimSpace(sf.Tag.Get("excel-pattern-rgb")); v != "" {
				localAttrs = append(localAttrs, AttrCellPatternRGB{}.Set(v))
			}
			if v := strings.TrimSpace(sf.Tag.Get("excel-font-bold")); v != "" {
				if b, err := strconv.ParseBool(v); err == nil {
					if b {
						localAttrs = append(localAttrs, AttrCellFontBold{}.SetTrue())
					} else {
						localAttrs = append(localAttrs, AttrCellFontBold{}.SetFalse())
					}
				}
			}
			if v := strings.TrimSpace(sf.Tag.Get("excel-font-italic")); v != "" {
				if b, err := strconv.ParseBool(v); err == nil {
					if b {
						localAttrs = append(localAttrs, AttrCellFontItalic{}.SetTrue())
					} else {
						localAttrs = append(localAttrs, AttrCellFontItalic{}.SetFalse())
					}
				}
			}
			if v := strings.TrimSpace(sf.Tag.Get("excel-border-rgb")); v != "" {
				parts := strings.Split(v, ",")
				for len(parts) < 4 {
					parts = append(parts, "")
				}
				localAttrs = append(localAttrs, AttrCellBorderRGB{}.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2]), strings.TrimSpace(parts[3])))
			}
			if v := strings.TrimSpace(sf.Tag.Get("excel-border-style")); v != "" {
				parts := strings.Split(v, ",")
				ints := make([]int, 4)
				for idx := 0; idx < 4 && idx < len(parts); idx++ {
					if iv, err := strconv.Atoi(strings.TrimSpace(parts[idx])); err == nil {
						ints[idx] = iv
					}
				}
				localAttrs = append(localAttrs, AttrCellBorderStyle{}.Set(ints[0], ints[1], ints[2], ints[3]))
			}
			if v := strings.TrimSpace(sf.Tag.Get("excel-diagonal-rgb")); v != "" {
				parts := strings.Split(v, ",")
				up, down := "", ""
				if len(parts) > 0 {
					up = strings.TrimSpace(parts[0])
				}
				if len(parts) > 1 {
					down = strings.TrimSpace(parts[1])
				}
				localAttrs = append(localAttrs, AttrCellDiagonalRGB{}.Set(up, down))
			}
			if v := strings.TrimSpace(sf.Tag.Get("excel-diagonal-style")); v != "" {
				parts := strings.Split(v, ",")
				up, down := 0, 0
				if len(parts) > 0 {
					if iv, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
						up = iv
					}
				}
				if len(parts) > 1 {
					if iv, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
						down = iv
					}
				}
				localAttrs = append(localAttrs, AttrCellDiagonalStyle{}.Set(up, down))
			}
			if v := strings.TrimSpace(sf.Tag.Get("excel-wrap-text")); v != "" {
				if b, err := strconv.ParseBool(v); err == nil {
					if b {
						localAttrs = append(localAttrs, AttrCellWrapText{}.SetTrue())
					} else {
						localAttrs = append(localAttrs, AttrCellWrapText{}.SetFalse())
					}
				}
			}

			// 处理指针字段
			if f.Kind() == reflect.Ptr {
				if f.IsNil() {
					// 空指针写空值
					merged := make([]CellAttributer, 0, len(attrs)+len(localAttrs))
					merged = append(merged, attrs...)
					merged = append(merged, localAttrs...)
					cell := APP.Cell.New(APP.CellAttr.Content.Set(""), APP.CellAttr.ContentType.Set(CellContentTypeAny)).setAttrs(merged...)
					writeCell(cell, rn+i, cn, my)
					continue
				}
				f = f.Elem()
			}

			var content any
			var ctype CellContentType = CellContentTypeAny

			switch f.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				content = int(f.Int())
				ctype = CellContentTypeInt
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				content = int(f.Uint())
				ctype = CellContentTypeInt
			case reflect.Float32, reflect.Float64:
				content = f.Float()
				ctype = CellContentTypeFloat
			case reflect.Bool:
				content = f.Bool()
				ctype = CellContentTypeBool
			case reflect.String:
				content = f.String()
				ctype = CellContentTypeAny
			case reflect.Struct:
				// special case: time.Time
				if f.Type() == reflect.TypeOf(time.Time{}) {
					content = f.Interface()
					ctype = CellContentTypeTime
				} else {
					// unsupported struct, marshal to string via fmt
					content = fmt.Sprintf("%v", f.Interface())
					ctype = CellContentTypeAny
				}
			default:
				// fallback to string representation
				content = fmt.Sprintf("%v", f.Interface())
				ctype = CellContentTypeAny
			}

			cell := APP.Cell.New(APP.CellAttr.Content.Set(content), APP.CellAttr.ContentType.Set(ctype))
			writeCell(cell, rn+i, cn, my)
		}
	}

	return my
}

// Download 下载文件
func (my *Writer) Download(w http.ResponseWriter) *Writer {
	my.lock.Lock()
	defer my.lock.Unlock()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", str.APP.Buffer.JoinString("attachment; filename=%s", url.QueryEscape(my.filename)))
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")

	my.Error = my.rawFile.Write(w)
	return my
}
