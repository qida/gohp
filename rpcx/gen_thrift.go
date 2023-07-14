package rpcx

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/qida/gohp/logx"
)

//e.g. 	AutoGetThrift(db.BaseImage{},db.TbBaseLabel{},db.TbBaseLabelType{},db.TbBaseSchemeAi{},db.TbRecLabelImage{}}

type GenThrift struct {
	PrefixDelete string
}

func (g *GenThrift) AutoGetThrift(values ...interface{}) error {
	for _, value := range values {
		var txtThrift string
		nameStruct := strings.Replace(strings.SplitAfterN(reflect.TypeOf(value).String(), ".", 2)[1], g.PrefixDelete, "", 1)
		txtThrift += fmt.Sprintf("namespace go %s \r\n", ToSnakeCase(nameStruct))
		txtThrift += fmt.Sprintf("struct %s {", nameStruct)
		t := reflect.TypeOf(value)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		for i := 0; i < t.NumField(); i++ {
			fieldName := t.Field(i).Name
			fieldType := t.Field(i).Type.String()
			switch fieldType {
			case "string",
				"time.Time":
				fieldType = "string"
			case "int",
				"uint",
				"uint64",
				"int64":
				fieldType = "i64 "
			case "int32",
				"uint32":
				fieldType = "i32"
			case "int16",
				"uint16":
				fieldType = "i16"
			case "int8",
				"byte",
				"uint8":
				fieldType = "i8"
			case "float64",
				"float32":
				fieldType = "double"
			case "bool":
				fieldType = "bool"
			default:
				// fieldType = "string"
			}
			tagSetting := ParseTagSetting(t.Field(i).Tag.Get("gorm"), ";")
			txtThrift += fmt.Sprintf("%4d: %8s %20s , //%s", i+1, fieldType, fieldName, tagSetting["COMMENT"])
		}
		txtThrift += fmt.Sprintln("}")
		logx.Info(txtThrift)
		CreateFile(fmt.Sprintf("%s.thrift", ToSnakeCase(nameStruct)), txtThrift)
	}

	return nil
}

func ToSnakeCase(str string) string {
	var result string
	for i, v := range str {
		if i > 0 && v >= 'A' && v <= 'Z' {
			result += "_"
		}
		result += strings.ToLower(string(v))
	}
	return result
}

func CreateFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}

func ParseTagSetting(str string, sep string) map[string]string {
	settings := map[string]string{}
	names := strings.Split(str, sep)

	for i := 0; i < len(names); i++ {
		j := i
		if len(names[j]) > 0 {
			for {
				if names[j][len(names[j])-1] == '\\' {
					i++
					names[j] = names[j][0:len(names[j])-1] + sep + names[i]
					names[i] = ""
				} else {
					break
				}
			}
		}

		values := strings.Split(names[j], ":")
		k := strings.TrimSpace(strings.ToUpper(values[0]))

		if len(values) >= 2 {
			settings[k] = strings.Join(values[1:], ":")
		} else if k != "" {
			settings[k] = k
		}
	}

	return settings
}
