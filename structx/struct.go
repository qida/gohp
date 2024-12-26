package structx

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func ShortContent(content interface{}, force ...string) (m map[string]interface{}) {
	m = make(map[string]interface{})
	var fields []string
	t := reflect.TypeOf(content)
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	val := reflect.ValueOf(content)
	var tt time.Time
	for _, fname := range fields {
		if val.FieldByName(fname).Interface() == "" ||
			val.FieldByName(fname).Interface() == 0 ||
			val.FieldByName(fname).Interface() == nil ||
			val.FieldByName(fname).Interface() == uint(0) ||
			val.FieldByName(fname).Interface() == int8(0) ||
			val.FieldByName(fname).Interface() == int64(0) ||
			val.FieldByName(fname).Interface() == float64(0) ||
			val.FieldByName(fname).Interface() == tt {
			if len(force) > 0 {
				for j := 0; j < len(force); j++ {
					if force[j] == fname {
						m[fname] = val.FieldByName(fname).Interface()
					}
				}
			}
		} else {
			// fmt.Printf("%s:%+v\r\n", fname, val.FieldByName(fname).Interface())
			m[fname] = val.FieldByName(fname).Interface()
		}
	}
	return
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)
	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}
	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}
	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("provided value type didn't match obj field type")
	}
	structFieldValue.Set(val)
	return nil
}

func CopyContent(master interface{}, second interface{}) (m map[string]interface{}) {
	m = make(map[string]interface{})
	var fields []string
	t := reflect.TypeOf(master)
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	val1 := reflect.ValueOf(master)
	val2 := reflect.ValueOf(second)
	var tt time.Time
	for _, fname := range fields {
		if val1.FieldByName(fname).Interface() == "" ||
			val1.FieldByName(fname).Interface() == 0 ||
			val1.FieldByName(fname).Interface() == nil ||
			val1.FieldByName(fname).Interface() == uint(0) ||
			val1.FieldByName(fname).Interface() == int8(0) ||
			val1.FieldByName(fname).Interface() == int64(0) ||
			val1.FieldByName(fname).Interface() == float64(0) ||
			val1.FieldByName(fname).Interface() == tt {
			m[fname] = val2.FieldByName(fname).Interface()
		} else {
			// fmt.Printf("%s:%+v\r\n", fname, val1.FieldByName(fname).Interface())
			m[fname] = val1.FieldByName(fname).Interface()
		}
	}
	return
}

// 对象转JSON字符串
func JsonToStr(i interface{}) string {
	byteI, _ := json.Marshal(&i)
	return string(byteI)
}

func RequestVerify(content interface{}, force ...string) (_err error) {
	var fields []string
	t := reflect.TypeOf(content)
	if t.Kind() == reflect.Ptr {
		// t = t.Elem()
		panic("字段校验传入对象不能为指针,请联系管理员")
	}
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	val := reflect.ValueOf(content)
	// var tt time.Time
	for j := 0; j < len(force); j++ {
		var found bool
		for _, fname := range fields {
			if force[j] == fname {
				found = true
				if reflect.ValueOf(val.FieldByName(fname).Interface()).IsZero() {
					_err = fmt.Errorf("%s 参数为空或0值,请联系管理员 [%+v]", fname, content)
					return
				}
				break //找到即对比下一组
			}
		}
		if !found {
			_err = errors.New("请求中检查字段 [" + force[j] + "] 不存在,请联系管理员")
		}
	}
	return
}
func RequestVerifyOr(content interface{}, force ...string) (_err error) {
	var fields []string
	t := reflect.TypeOf(content)
	if t.Kind() == reflect.Ptr {
		panic("字段校验传入对象不能为指针,请联系管理员")
	}
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	val := reflect.ValueOf(content)
	for j := 0; j < len(force); j++ {
		var found bool
		for _, fname := range fields {
			if force[j] == fname {
				found = true
				if !reflect.ValueOf(val.FieldByName(fname).Interface()).IsZero() {
					return
				}
				break //找到即对比下一组
			}
		}
		if !found {
			_err = errors.New("请求中检查字段 [" + force[j] + "] 不存在,请联系管理员")
			return
		}
	}
	_err = fmt.Errorf("请求中检查字段 %+v,全部为空,请联系管理员", force)
	return
}

// 删除字符串字段两边的空格
func RequestSpaceTrim(obj interface{}) (_err error) {
	// 使用反射获取对象的值
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr {
		_err = errors.New("对象不是指向结构体的地址")
		return
	}
	// 如果对象是指向结构体的指针，则获取其指向的结构体的值
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// 遍历结构体的所有字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		// 如果字段是字符串类型，则将其中的空格删除
		if field.Kind() == reflect.String {
			str := field.String()
			val.Field(i).SetString(strings.ReplaceAll(str, " ", ""))
		}
		// 如果字段是结构体，则递归调用该方法来处理其所有字段
		if field.Kind() == reflect.Struct {
			fieldVal := field.Addr().Interface()
			RequestSpaceTrim(fieldVal)
		}
	}
	return
}

// 通过json tag获取字段名称
func GetStructFieldNames(p interface{}, fields ...string) []string {
	var tags []string
	s := reflect.TypeOf(p)
	for i := 0; i < s.NumField(); i++ {
		tags = append(tags, s.Field(i).Tag.Get("json"))
	}
	if len(fields) == 0 {
		return tags
	}
	if strings.Contains(fields[0], "-") {
		for j := 0; j < len(fields); j++ {
			if !strings.Contains(fields[j], "-") {
				panic("不能混合使用")
			}
			for i := 0; i < len(tags); i++ {
				if fields[j] == tags[i] {
					tags = append(tags[:i], tags[(i+1):]...)
					continue
				}
			}
		}
	} else {
		for j := 0; j < len(fields); j++ {
			if strings.Contains(fields[j], "-") {
				panic("不能混合使用")
			}
		}
		return fields
	}
	return tags
}
