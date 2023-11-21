package tormstream

import (
	"sync"
	"text/template"

	"github.com/pkg/errors"
)

type TormI interface {
	TplName() string
	Identity() string
	Torm() string
}

var sqlTemplateMap sync.Map

func RegisterTorm(torm TormI) (err error) {
	tmp, err := template.New("").Funcs(TormfuncMapSQL).Parse(torm.Torm())
	if err != nil {
		return err
	}
	identity := torm.Identity()
	val, ok := sqlTemplateMap.Load(identity)
	var r *template.Template
	if ok {
		r = val.(*template.Template)
	}
	if r == nil {
		r = template.New("").Funcs(TormfuncMapSQL)
	}
	tpls := tmp.Templates()
	for _, tpl := range tpls {
		name := tpl.Name()
		if name == "" {
			continue
		}
		r.AddParseTree(name, tpl.Tree)
	}
	sqlTemplateMap.Store(torm.Identity(), r)
	return nil
}

func getSQLTpl(identify string) (r *template.Template, err error) {
	val, ok := sqlTemplateMap.Load(identify)
	if !ok {
		err = errors.Errorf("not found sqlTpl by identify:%s,use RegisterSQLTpl to set", identify)
		return nil, err
	}
	p, ok := val.(*template.Template)
	if !ok {
		err = errors.Errorf("required:%v,got:%v", template.New(""), val)
		return nil, err
	}
	return p, nil
}
