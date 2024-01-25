package torm

import (
	"reflect"

	"github.com/suifengpiao14/logchan/v2"
	"github.com/suifengpiao14/sqlexec"
)

// GetSQL 生成SQL(不关联DB操作)
func GetSQL(tplIdentify string, tplName string, volume VolumeInterface) (sqls string, namedSQL string, resetedVolume VolumeInterface, err error) {
	logInfo := &LogInfoToSQL{}
	defer func() {
		logInfo.TplIdentify = tplIdentify
		logInfo.TplName = tplName
		logInfo.InputVolume = volume
		logInfo.SQL = sqls
		logInfo.Named = namedSQL
		logInfo.TPLOutVolume = resetedVolume
		logInfo.Err = err
		logchan.SendLogInfo(logInfo)
	}()
	r, err := getSQLTpl(tplIdentify)
	if err != nil {
		return "", "", nil, err
	}

	namedSQL, resetedVolume, err = execTPL(r, tplName, volume)
	if err != nil {
		return "", "", nil, err
	}

	namedData, err := getNamedData(resetedVolume)
	if err != nil {
		return "", "", nil, err
	}
	logInfo.NamedData = namedData
	sqls, err = sqlexec.ExplainSQL(namedSQL, namedData)
	if err != nil {
		return "", "", nil, err
	}
	return sqls, namedSQL, resetedVolume, nil
}

func getNamedData(data interface{}) (out map[string]interface{}, err error) {
	out = make(map[string]interface{})
	if data == nil {
		return
	}
	dataI, ok := data.(*interface{})
	if ok {
		data = *dataI
	}
	mapOut, ok := data.(map[string]interface{})
	if ok {
		out = mapOut
		return
	}
	mapOutRef, ok := data.(*map[string]interface{})
	if ok {
		out = *mapOutRef
		return
	}
	if mapOut, ok := data.(VolumeMap); ok {
		out = mapOut
		return
	}
	if mapOutRef, ok := data.(*VolumeMap); ok {
		out = *mapOutRef
		return
	}

	v := reflect.Indirect(reflect.ValueOf(data))

	if v.Kind() != reflect.Struct {
		return
	}
	vt := v.Type()
	// 提取结构体field字段
	fieldNum := v.NumField()
	for i := 0; i < fieldNum; i++ {
		fv := v.Field(i)
		ft := fv.Type()
		fname := vt.Field(i).Name
		if fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
			ft = fv.Type()
		}
		ftk := ft.Kind()
		switch ftk {
		case reflect.Int:
			out[fname] = fv.Int()
		case reflect.Int64:
			out[fname] = int64(fv.Int())
		case reflect.Float64:
			out[fname] = fv.Float()
		case reflect.String:
			out[fname] = fv.String()
		case reflect.Struct, reflect.Map:
			subOut, err := getNamedData(fv.Interface())
			if err != nil {
				return out, err
			}
			for k, v := range subOut {
				if _, ok := out[k]; !ok {
					out[k] = v
				}
			}

		default:
			out[fname] = fv.Interface()
		}
	}
	return
}
