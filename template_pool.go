package tormstream

import (
	"sync"
	"text/template"

	"github.com/pkg/errors"
)

var sqlTemplateMap sync.Map

type SqlTplInstance struct {
	sqlTplIdentify string
	tpl            *template.Template
}

var ERROR_SQL_TEMPLATE_NOT_FOUND_DB = errors.New("sqlTplInstance.dbInstance is nil")
var ERROR_SQL_TEMPLATE_NOT_FOUND_TEMPLATE = errors.New("sqlTplInstance.dbInstance is nil")
var ERROR_DB_EXECUTOR_GETTER_REQUIRD = errors.Errorf("dbExecutorGetter required")
var ERROR_DB_EXECUTOR_REQUIRD = errors.Errorf("dbExecutor required")

func (ins *SqlTplInstance) GetTemplate() (r *template.Template) {
	return ins.tpl
}

func RegisterSQLTpl(sqlTplIdentify string, r *template.Template) {
	instance := SqlTplInstance{
		sqlTplIdentify: sqlTplIdentify,
		tpl:            r,
	}
	sqlTemplateMap.Store(sqlTplIdentify, &instance)
}

func GetSQLTpl(identify string) (sqlTplInstance *SqlTplInstance, err error) {
	val, ok := sqlTemplateMap.Load(identify)
	if !ok {
		err = errors.Errorf("not found db by identify:%s,use RegisterSQLTpl to set", identify)
		return nil, err
	}
	p, ok := val.(*SqlTplInstance)
	if !ok {
		err = errors.Errorf("required:%v,got:%v", &SqlTplInstance{}, val)
		return nil, err
	}
	return p, nil
}
