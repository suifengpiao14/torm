package torm

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/funcs"
)

const (
	EOF                  = "\n"
	WINDOW_EOF           = "\r\n"
	HTTP_HEAD_BODY_DELIM = EOF + EOF
)

func execTPL(t *template.Template, tplName string, volume VolumeInterface) (namedSQL string, resetedVolume VolumeInterface, err error) {
	var b bytes.Buffer
	err = t.ExecuteTemplate(&b, tplName, volume)
	if err != nil {
		err = errors.WithStack(err)
		return "", nil, err
	}
	namedSQL = strings.ReplaceAll(b.String(), WINDOW_EOF, EOF)
	namedSQL = funcs.TrimSpaces(namedSQL)
	return namedSQL, volume, nil
}

type VolumeInterface interface {
	SetValue(key string, value interface{})
	GetValue(key string, value interface{}) (ok bool)
}

type VolumeMap map[string]interface{}

func NewVolumeMap() *VolumeMap {
	return &VolumeMap{}
}

func (v *VolumeMap) init() {
	if v == nil {
		err := errors.Errorf("*Templatemap must init")
		panic(err)
	}
	if *v == nil {
		*v = VolumeMap{} // 解决 data33 情况
	}
}

func (v *VolumeMap) SetValue(key string, value interface{}) {
	v.init()
	(*v)[key] = value

}

func (v *VolumeMap) GetValue(key string, value interface{}) (ok bool) {
	v.init()
	tmp, ok := (*v)[key]
	if !ok {
		return ok
	}
	ok = convertType(value, tmp)
	return ok
}

func (v *VolumeMap) String() (out string) {
	b, _ := json.Marshal(v)
	out = string(b)
	return out
}

func convertType(dst interface{}, src interface{}) bool {
	if src == nil || dst == nil {
		return false
	}
	rv := reflect.Indirect(reflect.ValueOf(dst))
	if !rv.CanSet() {
		err := errors.Errorf("dst :%#v reflect.CanSet() must return  true", dst)
		panic(err)
	}
	rvT := rv.Type()

	rTmp := reflect.ValueOf(src)
	if rTmp.CanConvert(rvT) {
		realValue := rTmp.Convert(rvT)
		rv.Set(realValue)
		return true
	}
	srcStr := cast.ToString(src)
	switch rvT.Kind() {
	case reflect.Int:
		srcInt, err := strconv.Atoi(srcStr)
		if err != nil {
			err = errors.WithMessagef(err, "src:%s can`t convert to int", srcStr)
			panic(err)
		}
		rv.Set(reflect.ValueOf(srcInt))
		return true
	case reflect.Int64:
		srcInt, err := strconv.ParseInt(srcStr, 10, 64)
		if err != nil {
			err = errors.WithMessagef(err, "src:%s can`t convert to int64", srcStr)
			panic(err)
		}
		rv.SetInt(int64(srcInt))
		return true
	case reflect.Float64:
		srcFloat, err := strconv.ParseFloat(srcStr, 64)
		if err != nil {
			err = errors.WithMessagef(err, "src:%s can`t convert to float64", srcStr)
			panic(err)
		}
		rv.SetFloat(srcFloat)
		return true
	case reflect.Bool:
		srcBool, err := strconv.ParseBool(srcStr)
		if err != nil {
			err = errors.WithMessagef(err, "src:%s can`t convert to bool", srcStr)
			panic(err)
		}
		rv.SetBool(srcBool)
		return true
	case reflect.String:
		rv.SetString(srcStr)
		return true

	}
	err := errors.Errorf("can not convert %v(%s) to %#v", src, rTmp.Type().String(), rvT.String())
	panic(err)
}
