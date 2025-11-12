package common

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	emptyStringMap = map[string]struct{}{
		"":      {},
		"0":     {},
		"no":    {},
		"off":   {},
		"false": {},
	}
)

type iVal interface{ Val() any }

// iString is used for type assert api for String().
type iString interface{ String() string }

// iBool is used for type assert api for Bool().
type iBool interface{ Bool() bool }

// iInt64 is used for type assert api for Int64().
type iInt64 interface{ Int64() int64 }

// iUint64 is used for type assert api for Uint64().
type iUint64 interface{ Uint64() uint64 }

// iFloat32 is used for type assert api for Float32().
type iFloat32 interface{ Float32() float32 }

// iFloat64 is used for type assert api for Float64().
type iFloat64 interface{ Float64() float64 }

// iError is used for type assert api for Error().
type iError interface{ Error() string }

// iBytes is used for type assert api for Bytes().
type iBytes interface{ Bytes() []byte }

// iInterface is used for type assert api for Interface().
type iInterface interface{ Interface() any }

// iInterfaces is used for type assert api for Interfaces().
type iInterfaces interface{ Interfaces() []any }

// iFloats is used for type assert api for Floats().
type iFloats interface{ Floats() []float64 }

// iInts is used for type assert api for Ints().
type iInts interface{ Ints() []int }

// iStrings is used for type assert api for Strings().
type iStrings interface{ Strings() []string }

// iUints is used for type assert api for Uints().
type iUints interface{ Uints() []uint }

// iMapStrAny is the interface support for converting struct parameter to map.
type iMapStrAny interface{ MapStrAny() map[string]any }

// iUnmarshalValue is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface iUnmarshalValue.
type iUnmarshalValue interface{ UnmarshalValue(any) error }

// iUnmarshalText is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface iUnmarshalText.
type iUnmarshalText interface{ UnmarshalText(text []byte) error }

// iUnmarshalText is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface iUnmarshalJSON.
type iUnmarshalJSON interface{ UnmarshalJSON(b []byte) error }

// iSet is the interface for custom value assignment.
type iSet interface{ Set(value any) (old any) }

// FloatFormat 浮点精度控制
func FloatFormat(f float64, length int) (float64, error) {
	return strconv.ParseFloat(strconv.FormatFloat(f, 'f', length, 64), 64)
}

func ToBytes(original any) []byte {
	if original == nil {
		return nil
	}

	switch value := original.(type) {
	case string:
		return []byte(value)

	case []byte:
		return value

	default:
		if f, ok := value.(iBytes); ok {
			return f.Bytes()
		}

		refType := reflect.TypeOf(original)
		refValue := reflect.ValueOf(original)

		switch refType.Kind() {
		case reflect.Map:
			bytes, err := json.Marshal(original)
			if err != nil {
				panic(err)
			}
			return bytes

		case reflect.Array, reflect.Slice:
			var (
				ok    = true
				bytes = make([]byte, refValue.Len())
			)

			for i := range bytes {
				int32Value := ToInt32(refValue.Index(i).Interface())
				if int32Value < 0 || int32Value > math.MaxUint8 {
					ok = false
					break
				}
				bytes[i] = byte(int32Value)
			}

			if ok {
				return bytes
			}
		}

		return LeEncode(original)
	}
}

func ToUint(original any) uint {
	if original == nil {
		return 0
	}
	if v, ok := original.(uint); ok {
		return v
	}

	return uint(ToInt64(original))
}

func ToUint8(original any) uint8 {
	if original == nil {
		return 0
	}
	if v, ok := original.(uint8); ok {
		return v
	}

	return uint8(ToInt64(original))
}

func ToUint16(original any) uint16 {
	if original == nil {
		return 0
	}
	if v, ok := original.(uint16); ok {
		return v
	}

	return uint16(ToInt64(original))
}

func ToUint32(original any) uint32 {
	if original == nil {
		return 0
	}
	if v, ok := original.(uint32); ok {
		return v
	}

	return uint32(ToInt64(original))
}

func ToUint64(original any) uint64 {
	if original == nil {
		return 0
	}
	switch value := original.(type) {
	case int:
		return uint64(value)
	case int8:
		return uint64(value)
	case int16:
		return uint64(value)
	case int32:
		return uint64(value)
	case int64:
		return uint64(value)
	case uint:
		return uint64(value)
	case uint8:
		return uint64(value)
	case uint16:
		return uint64(value)
	case uint32:
		return uint64(value)
	case uint64:
		return value
	case float32:
		return uint64(value)
	case float64:
		return uint64(value)
	case bool:
		if value {
			return 1
		}
		return 0
	case []byte:
		return LeDecodeToUint64(value)
	default:
		if f, ok := value.(iUint64); ok {
			return f.Uint64()
		}
		var s = ToString(value)
		if len(s) > 0 {
			if s[0] == '-' {
				s = s[1:]
			} else if s[0] == '+' {
				s = s[1:]
			}
		}
		// Hexadecimal
		if len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
			if v, e := strconv.ParseUint(s[2:], 16, 64); e == nil {
				return v
			}
		}
		// Decimal
		if v, e := strconv.ParseUint(s, 10, 64); e == nil {
			return v
		}
		// Float64
		if valueInt64 := ToFloat64(value); math.IsNaN(valueInt64) {
			return 0
		} else {
			return uint64(valueInt64)
		}
	}
}

func ToInt(original any) int {
	if original == nil {
		return 0
	}
	if v, ok := original.(int); ok {
		return v
	}

	return int(ToInt64(original))
}

func ToInt8(original any) int8 {
	if original == nil {
		return 0
	}
	if v, ok := original.(int8); ok {
		return v
	}

	return int8(ToInt64(original))
}

func ToInt16(original any) int16 {
	if original == nil {
		return 0
	}
	if v, ok := original.(int16); ok {
		return v
	}

	return int16(ToInt64(original))
}

func ToInt32(original any) int32 {
	if original == nil {
		return 0
	}
	if v, ok := original.(int32); ok {
		return v
	}

	return int32(ToInt64(original))
}

func ToInt64(original any) int64 {
	if original == nil {
		return 0
	}
	switch value := original.(type) {
	case int:
		return int64(value)
	case int8:
		return int64(value)
	case int16:
		return int64(value)
	case int32:
		return int64(value)
	case int64:
		return value
	case uint:
		return int64(value)
	case uint8:
		return int64(value)
	case uint16:
		return int64(value)
	case uint32:
		return int64(value)
	case uint64:
		return int64(value)
	case float32:
		return int64(value)
	case float64:
		return int64(value)
	case bool:
		if value {
			return 1
		}
		return 0
	case []byte:
		return LeDecodeToInt64(value)
	default:
		if f, ok := value.(iInt64); ok {
			return f.Int64()
		}
		var (
			s       = ToString(value)
			isMinus = false
		)
		if len(s) > 0 {
			if s[0] == '-' {
				isMinus = true
				s = s[1:]
			} else if s[0] == '+' {
				s = s[1:]
			}
		}
		// Hexadecimal
		if len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
			if v, e := strconv.ParseInt(s[2:], 16, 64); e == nil {
				if isMinus {
					return -v
				}
				return v
			}
		}
		// Decimal
		if v, e := strconv.ParseInt(s, 10, 64); e == nil {
			if isMinus {
				return -v
			}
			return v
		}
		// Float64
		if valueInt64 := ToFloat64(value); math.IsNaN(valueInt64) {
			return 0
		} else {
			return int64(valueInt64)
		}
	}
}

func ToDatetimeUnixMilli(original any) time.Time { return time.UnixMilli(ToInt64(original)) }

// ToFloat32 converts `any` to float32.
func ToFloat32(original any) float32 {
	if original == nil {
		return 0
	}
	switch value := original.(type) {
	case float32:
		return value
	case float64:
		return float32(value)
	case []byte:
		return LeDecodeToFloat32(value)
	default:
		if f, ok := value.(iFloat32); ok {
			return f.Float32()
		}
		v, _ := strconv.ParseFloat(ToString(original), 64)
		return float32(v)
	}
}

// ToFloat64 converts `any` to float64.
func ToFloat64(original any) float64 {
	if original == nil {
		return 0
	}
	switch value := original.(type) {
	case float32:
		return float64(value)
	case float64:
		return value
	case []byte:
		return LeDecodeToFloat64(value)
	default:
		if f, ok := value.(iFloat64); ok {
			return f.Float64()
		}
		v, _ := strconv.ParseFloat(ToString(original), 64)
		return v
	}
}

func ToString(original any) string {
	if original == nil {
		return ""
	}
	switch value := original.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *time.Time:
		if value == nil {
			return ""
		}
		return value.String()
	default:
		// Empty checks.
		// if value == nil {
		// 	return ""
		// }
		if f, ok := value.(iString); ok {
			// If the variable implements the String() interface,
			// then use that interface to perform the conversion
			return f.String()
		}
		if f, ok := value.(iError); ok {
			// If the variable implements the Error() interface,
			// then use that interface to perform the conversion
			return f.Error()
		}
		// Reflection checks.
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
		case reflect.String:
			return rv.String()
		}
		if kind == reflect.Ptr {
			return ToString(rv.Elem().Interface())
		}
		// Finally, we use json.Marshal to convert.
		if jsonContent, err := json.Marshal(value); err != nil {
			return fmt.Sprint(value)
		} else {
			return string(jsonContent)
		}
	}
}

func ToBool(original any) bool {
	if original == nil {
		return false
	}
	switch value := original.(type) {
	case bool:
		return value
	case []byte:
		if _, ok := emptyStringMap[strings.ToLower(string(value))]; ok {
			return false
		}
		return true
	case string:
		if _, ok := emptyStringMap[strings.ToLower(value)]; ok {
			return false
		}
		return true
	default:
		if f, ok := value.(iBool); ok {
			return f.Bool()
		}
		rv := reflect.ValueOf(original)
		switch rv.Kind() {
		case reflect.Ptr:
			return !rv.IsNil()
		case reflect.Map:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			return rv.Len() != 0
		case reflect.Struct:
			return true
		default:
			s := strings.ToLower(ToString(original))
			if _, ok := emptyStringMap[s]; ok {
				return false
			}
			return true
		}
	}
}

// ToStruct 任意类型转结构体
func ToStruct(v any, r any) error {
	var (
		e        error
		jsonByte []byte
	)
	jsonByte, e = json.Marshal(v)
	if e != nil {
		return e
	}

	return json.Unmarshal(jsonByte, &r)
}

func LeEncode(values ...any) []byte {
	buf := new(bytes.Buffer)
	for i := 0; i < len(values); i++ {
		if values[i] == nil {
			return buf.Bytes()
		}
		switch value := values[i].(type) {
		case int:
			buf.Write(LeEncodeInt(value))
		case int8:
			buf.Write(LeEncodeInt8(value))
		case int16:
			buf.Write(LeEncodeInt16(value))
		case int32:
			buf.Write(LeEncodeInt32(value))
		case int64:
			buf.Write(LeEncodeInt64(value))
		case uint:
			buf.Write(LeEncodeUint(value))
		case uint8:
			buf.Write(LeEncodeUint8(value))
		case uint16:
			buf.Write(LeEncodeUint16(value))
		case uint32:
			buf.Write(LeEncodeUint32(value))
		case uint64:
			buf.Write(LeEncodeUint64(value))
		case bool:
			buf.Write(LeEncodeBool(value))
		case string:
			buf.Write(LeEncodeString(value))
		case []byte:
			buf.Write(value)
		case float32:
			buf.Write(LeEncodeFloat32(value))
		case float64:
			buf.Write(LeEncodeFloat64(value))

		default:
			if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
				buf.Write(LeEncodeString(fmt.Sprintf("%v", value)))
			}
		}
	}

	return buf.Bytes()
}

func LeEncodeByLength(length int, values ...any) []byte {
	b := LeEncode(values...)
	if len(b) < length {
		b = append(b, make([]byte, length-len(b))...)
	} else if len(b) > length {
		b = b[0:length]
	}

	return b
}

func LeDecode(b []byte, values ...any) error {
	var (
		err error
		buf = bytes.NewBuffer(b)
	)
	for i := 0; i < len(values); i++ {
		if err = binary.Read(buf, binary.LittleEndian, values[i]); err != nil {
			return err
		}
	}

	return nil
}

func LeEncodeString(s string) []byte { return []byte(s) }

func LeDecodeToString(b []byte) string { return string(b) }

func LeEncodeBool(b bool) []byte {
	if b {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

func LeEncodeInt(i int) []byte {
	if i <= math.MaxInt8 {
		return LeEncodeInt8(int8(i))
	} else if i <= math.MaxInt16 {
		return LeEncodeInt16(int16(i))
	} else if i <= math.MaxInt32 {
		return LeEncodeInt32(int32(i))
	} else {
		return LeEncodeInt64(int64(i))
	}
}

func LeEncodeUint(i uint) []byte {
	if i <= math.MaxUint8 {
		return LeEncodeUint8(uint8(i))
	} else if i <= math.MaxUint16 {
		return LeEncodeUint16(uint16(i))
	} else if i <= math.MaxUint32 {
		return LeEncodeUint32(uint32(i))
	} else {
		return LeEncodeUint64(uint64(i))
	}
}

func LeEncodeInt8(i int8) []byte { return []byte{byte(i)} }

func LeEncodeUint8(i uint8) []byte { return []byte{i} }

func LeEncodeInt16(i int16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(i))

	return b
}

func LeEncodeUint16(i uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, i)

	return b
}

func LeEncodeInt32(i int32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))

	return b
}

func LeEncodeUint32(i uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)

	return b
}

func LeEncodeInt64(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))

	return b
}

func LeEncodeUint64(i uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)

	return b
}

func LeEncodeFloat32(f float32) []byte {
	bits := math.Float32bits(f)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, bits)

	return b
}

func LeEncodeFloat64(f float64) []byte {
	bits := math.Float64bits(f)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, bits)

	return b
}

func LeDecodeToInt(b []byte) int {
	if len(b) < 2 {
		return int(LeDecodeToUint8(b))
	} else if len(b) < 3 {
		return int(LeDecodeToUint16(b))
	} else if len(b) < 5 {
		return int(LeDecodeToUint32(b))
	} else {
		return int(LeDecodeToUint64(b))
	}
}

func LeDecodeToUint(b []byte) uint {
	if len(b) < 2 {
		return uint(LeDecodeToUint8(b))
	} else if len(b) < 3 {
		return uint(LeDecodeToUint16(b))
	} else if len(b) < 5 {
		return uint(LeDecodeToUint32(b))
	} else {
		return uint(LeDecodeToUint64(b))
	}
}

func LeDecodeToBool(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	if bytes.Equal(b, make([]byte, len(b))) {
		return false
	}

	return true
}

func LeDecodeToInt8(b []byte) int8 {
	if len(b) == 0 {
		panic(`empty slice given`)
	}

	return int8(b[0])
}

func LeDecodeToUint8(b []byte) uint8 {
	if len(b) == 0 {
		panic(`empty slice given`)
	}

	return b[0]
}

func LeDecodeToInt16(b []byte) int16 {
	return int16(binary.LittleEndian.Uint16(LeFillUpSize(b, 2)))
}

func LeDecodeToUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(LeFillUpSize(b, 2))
}

func LeDecodeToInt32(b []byte) int32 {
	return int32(binary.LittleEndian.Uint32(LeFillUpSize(b, 4)))
}

func LeDecodeToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(LeFillUpSize(b, 4))
}

func LeDecodeToInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(LeFillUpSize(b, 8)))
}

func LeDecodeToUint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(LeFillUpSize(b, 8))
}

func LeDecodeToFloat32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(LeFillUpSize(b, 4)))
}

func LeDecodeToFloat64(b []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(LeFillUpSize(b, 8)))
}

// LeFillUpSize fills up the bytes `b` to given length `l` using LittleEndian.
//
// Note that it creates a new bytes slice by copying the original one to avoid changing
// the original parameter bytes.
func LeFillUpSize(b []byte, l int) []byte {
	if len(b) >= l {
		return b[:l]
	}
	c := make([]byte, l)
	copy(c, b)

	return c
}
