package data

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
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
		return fmt.Errorf("No such field: %s in obj", name)
	}
	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}
	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
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

func MD5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func MD5Str(str string) string {
	return MD5Byte([]byte(str))
}

func MD5Byte(byts []byte) string {
	h := md5.New()
	h.Write(byts)
	return hex.EncodeToString(h.Sum(nil))
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
