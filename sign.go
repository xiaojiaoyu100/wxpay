package wxpay

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func checkSign(stream []byte, key string) (err error) {
    defer func() {
        if err != nil {
            notifyAsync(string(stream), err)
        }
    }()
    reqMap := make(map[string]string)
    err = xml.Unmarshal(stream, (*Map)(&reqMap))
    if err != nil {
        return
    }

    if v, ok := reqMap["sign"]; ok {
        md5Sign := sign(reqMap, key)
        if v != md5Sign {
            err = signNotMatchErr
            return
        }
    }
    return
}

func sign(req map[string]string, key string) string {
	// #1.对参数按照key=value的格式，并按照参数名ASCII字典序排序生成字符串：
	sortedKeys := make([]string, 0)
	for k := range req {
		if k == "sign" {
			continue
		}
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	var signStrings string
	for _, k := range sortedKeys {
		value := req[k]
		if value != "" {
			signStrings += k + "=" + value + "&"
		}
	}

	// #2.连接商户key：
	signStrings = signStrings + "key=" + key

	// #3.生成sign并转成大写：
	hash := md5.New()
	hash.Write([]byte(signStrings))
	upperSign := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))

	// #4.校验结果：
	globalLogger.printf("待签名字符串: %s", signStrings)
	globalLogger.printf("签名: %s", upperSign)
	return upperSign
}

func signStruct(v interface{}, key string) string {
	req := convert(v)
	return sign(req, key)
}

func convert(str interface{}) map[string]string {
	m := make(map[string]string)
	val := reflect.ValueOf(str).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		tag := typeField.Tag.Get("xml")

		list := strings.Split(tag, ",")
		if len(list) == 0 {
			continue
		}

		key := list[0]

		if key == "xml" ||
			key == "" ||
			key == "-" {
			continue
		}

		value := fmt.Sprintf("%v", valueField)

		if strings.Contains(tag, "omitempty") {
			switch valueField.Kind() {
			case reflect.String:
				if value == "" {
					continue
				}
			case reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64,
				reflect.Uint,
				reflect.Uint8,
				reflect.Uint16,
				reflect.Uint32,
				reflect.Uint64,
				reflect.Float32,
				reflect.Float64:
				if value == "0" {
					continue
				}

			case reflect.Ptr:
				if valueField.IsNil() {
					continue
				}
			}
		}

		m[key] = value
	}

	return m
}
