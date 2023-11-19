package tormstream

import (
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/suifengpiao14/funcs"
	"github.com/suifengpiao14/logchan/v2"
	gormLogger "gorm.io/gorm/logger"
)

// GetSQL 生成SQL(不关联DB操作)
func GetSQL(sqlTplIdentify string, tplName string, volume VolumeInterface) (sqls string, namedSQL string, resetedVolume VolumeInterface, err error) {
	sqlTplInstance, err := GetSQLTpl(sqlTplIdentify)
	if err != nil {
		return "", "", nil, err
	}
	namedSQL, resetedVolume, err = ExecTPL(sqlTplInstance.tpl, tplName, volume)
	if err != nil {
		return "", "", nil, err
	}
	sqls, err = ToSQL(namedSQL, resetedVolume)
	if err != nil {
		return "", "", nil, err
	}
	return sqls, namedSQL, resetedVolume, nil
}

// ToSQL 将字符串、数据整合为sql
func ToSQL(namedSql string, data interface{}) (sql string, err error) {
	namedSql = funcs.StandardizeSpaces(funcs.TrimSpaces(namedSql)) // 格式化sql语句
	logInfo := &LogInfoToSQL{
		Named: namedSql,
		Data:  data,
		Err:   err,
	}

	defer func() {
		logInfo.SQL = sql
		logInfo.Err = err
		logchan.SendLogInfo(logInfo)
	}()
	namedData, err := getNamedData(data)
	if err != nil {
		return "", err
	}
	logInfo.NamedData = namedData
	statment, arguments, err := sqlx.Named(namedSql, namedData)
	if err != nil {
		err = errors.WithStack(err)
		return "", err
	}
	sql = gormLogger.ExplainSQL(statment, nil, `'`, arguments...)
	return sql, nil
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
