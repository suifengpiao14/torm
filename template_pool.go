package torm

import (
	"sync"
	"text/template"

	"github.com/pkg/errors"
)

type TormI interface {
	TplName() string
	Namespace() string
	Torm() string
}

var sqlTemplateMap sync.Map

const (
	Torm_DELIM_LEFT  = "{{"
	Torm_DELIM_RIGHT = "}}"
)

//NewTemplate 方便外部初始化模板函数
func NewTemplate() (t *template.Template) {
	return template.New("").Delims(Torm_DELIM_LEFT, Torm_DELIM_RIGHT).Funcs(TormfuncMapSQL)
}

func RegisterTorm(torm TormI) (err error) {
	tmp, err := NewTemplate().Parse(torm.Torm())
	if err != nil {
		return err
	}
	identity := torm.Namespace()
	val, ok := sqlTemplateMap.Load(identity)
	var r *template.Template
	if ok {
		r = val.(*template.Template)
	}
	if r == nil {
		r = NewTemplate()
	}
	tpls := tmp.Templates()
	for _, tpl := range tpls {
		name := tpl.Name()
		if name == "" {
			continue
		}
		r.AddParseTree(name, tpl.Tree)
	}
	sqlTemplateMap.Store(torm.Namespace(), r)
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
