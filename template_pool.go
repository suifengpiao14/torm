package tormstream

import (
	"sync"
	"text/template"

	"github.com/pkg/errors"
)

var sqlTemplateMap sync.Map

type _SqlTplInstance struct {
	sqlTplIdentify string
	tpl            *template.Template
}

func RegisterSQLTpl(sqlTplIdentify string, r *template.Template) {
	instance := _SqlTplInstance{
		sqlTplIdentify: sqlTplIdentify,
		tpl:            r,
	}
	sqlTemplateMap.Store(sqlTplIdentify, &instance)
}

func getSQLTpl(identify string) (sqlTplInstance *_SqlTplInstance, err error) {
	val, ok := sqlTemplateMap.Load(identify)
	if !ok {
		err = errors.Errorf("not found sqlTpl by identify:%s,use RegisterSQLTpl to set", identify)
		return nil, err
	}
	p, ok := val.(*_SqlTplInstance)
	if !ok {
		err = errors.Errorf("required:%v,got:%v", &_SqlTplInstance{}, val)
		return nil, err
	}
	return p, nil
}
