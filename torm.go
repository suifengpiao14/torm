package torm

import (
	"strings"
	"text/template"
	"text/template/parse"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/packethandler"
	"github.com/suifengpiao14/pathtransfer"
	"github.com/tidwall/gjson"
	"golang.org/x/net/context"
)

// Torm 模板和执行器之间存在确定关系，在配置中体现, 同一个Torm 下template 内的define 共用相同资源
type Torm struct {
	Namespace        string                 `json:"namespace"`
	TplName          string                 `json:"tplName"`
	SubTemplateNames []string               `json:"subTemplateNames"`
	Source           Source                 `json:"source"`
	TplText          string                 `json:"tpl"`
	Transfers        pathtransfer.Transfers `json:"transfers"`
	PacketHandlers   packethandler.PacketHandlers
	Flow             packethandler.Flow `json:"flow"`
	template         *template.Template
}

type Torms []Torm

func (t Torm) Name() string {
	return pathtransfer.JoinPath(t.Namespace, t.TplName).String()
}

func (t Torm) GetIONamespace() (inNamespace string, outNamespace string) {
	inNamespace, outNamespace = pathtransfer.JoinPath(t.Namespace, t.TplName, pathtransfer.Transfer_Direction_input).String(), pathtransfer.JoinPath(t.TplName, pathtransfer.Transfer_Direction_output).String()
	return inNamespace, outNamespace
}

//TrimOutNamespace 去除输出的io命名空间
func (t Torm) TrimOutNamespace(input []byte) (out string, err error) {
	_, outNamespace := t.GetIONamespace()
	result := gjson.GetBytes(input, outNamespace)
	if result.Exists() {
		out = result.String()
		return out, nil
	}
	_, outTransfer := t.Transfers.SplitInOut()
	path := outTransfer.ModifySrcPath(func(path pathtransfer.Path) (newPath pathtransfer.Path) {
		return pathtransfer.Path(path.TrimIONamespace())
	}).Reverse().GjsonPath()
	result = gjson.GetBytes(input, path)
	if !result.Exists() {
		err = errors.Errorf("io/dictionary key not exists:io key:%s,dictionary key:%s;input:%s", outNamespace, path, string(input))
		return "", err
	}
	out = result.String()

	return out, nil
}

func (t Torm) GetRootTemplate() (template *template.Template) {
	return t.template
}
func (t Torm) Run(ctx context.Context, input []byte) (out []byte, err error) {
	packetHandlers, err := t.PacketHandlers.GetByName(t.Flow...)
	if err != nil {
		return nil, err
	}
	out, err = packetHandlers.Run(ctx, input)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// 解析tpl 文本，生成 Torms
func ParserTpl(source *Source, tplText string) (torms Torms, err error) {
	t := NewTemplate()
	t, err = t.Parse(tplText)
	if err != nil {
		return nil, err
	}
	prov := source.Provider
	if prov == nil {
		err = errors.Errorf("source provider requird source:%T", source)
		return nil, err
	}

	torms = make(Torms, 0)

	for _, tpl := range t.Templates() {
		tplName := tpl.Name()
		if tplName == "" {
			continue
		}
		torm := &Torm{
			TplName:          tplName,
			SubTemplateNames: make([]string, 0),
			Source:           *source,
			TplText:          tpl.Root.String(),
			template:         t, // 这里使用根模板，方便解决子模板依赖问题
		}
		torm.SubTemplateNames, err = GetSubTemplateNames(tpl, tplName)
		if err != nil {
			return nil, err
		}
		torms.AddReplace(*torm)
	}
	return torms, nil
}

func (ts *Torms) AddReplace(subTorms ...Torm) {
	if *ts == nil {
		*ts = make(Torms, 0)
	}
	for _, subTorm := range subTorms {
		exists := false
		for i, tor := range *ts {
			if strings.EqualFold(subTorm.TplName, tor.TplName) {
				(*ts)[i] = subTorm
				exists = true
			}
		}
		if !exists {
			*ts = append(*ts, subTorm)
		}

	}
}
func (ts Torms) GetByTplName(name string) (t *Torm, err error) {
	for _, t := range ts {
		if strings.EqualFold(name, t.TplName) {
			return &t, nil
		}
	}
	err = errors.Errorf("not found torm named:%s", name)
	return nil, err
}

func (ts *Torms) Transfers() (pathTransfers pathtransfer.Transfers) {
	pathTransfers = make(pathtransfer.Transfers, 0)
	for _, t := range *ts {
		pathTransfers.AddReplace(t.Transfers...)
	}

	return pathTransfers
}
func (ts *Torms) Template() (allTpl *template.Template, err error) {
	allTpl = NewTemplate()
	for _, t := range *ts {
		allTpl.AddParseTree(t.template.Name(), t.template.Tree)
	}
	return allTpl, nil
}

// GetSubTemplateNames 遍历 TemplateNode 节点
func GetSubTemplateNames(templ *template.Template, tplName string) (subTemplateNames []string, err error) {
	if templ == nil {
		err = errors.Errorf("GetSubTemplateNames:  *template.Template required")
		return nil, err
	}
	subTemplateNames = make([]string, 0)
	t := templ.Lookup(tplName)

	if t == nil {
		err = errors.Errorf("template: no template %s associated with template %s", templ.Name(), tplName)
		return nil, err
	}

	Traverse(t.Root, func(node parse.Node) {
		switch n := node.(type) {
		case *parse.TemplateNode:
			subTemplateNames = append(subTemplateNames, n.Name)
			var subNames []string // 此处需要单独声明，方便将 err 带出到外层
			subNames, err = GetSubTemplateNames(templ.Lookup(n.Name), n.Name)
			if err != nil {
				return
			}
			subTemplateNames = append(subTemplateNames, subNames...)

		}
	})
	if err != nil {
		return nil, err
	}
	return subTemplateNames, nil
}

const (
	Torm_DELIM_LEFT  = "{{"
	Torm_DELIM_RIGHT = "}}"
)

// NewTemplate 方便外部初始化模板函数
func NewTemplate() (t *template.Template) {
	return template.New("").Delims(Torm_DELIM_LEFT, Torm_DELIM_RIGHT).Funcs(TormfuncMapSQL)
}
