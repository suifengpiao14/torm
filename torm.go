package torm

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/packethandler"
	"github.com/suifengpiao14/pathtransfer"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Torm 模板和执行器之间存在确定关系，在配置中体现, 同一个Torm 下template 内的define 共用相同资源
type Torm struct {
	Name           string                 `json:"name"`
	Source         Source                 `json:"source"`
	TplText        string                 `json:"tpl"`
	Transfers      pathtransfer.Transfers `json:"transfers"`
	PacketHandlers packethandler.PacketHandlers
	Flow           packethandler.Flow `json:"flow"`
	template       *template.Template
}
type Torms []Torm

// NamespaceInput 从标准库转换过来的命名空间
func (t Torm) NamespaceInput() (namespaceInput string) {
	return fmt.Sprintf("%s.input", t.Name)
}

// NamespaceInput 返回到标准库时，增加的命名空间
func (t Torm) NamespaceOutput() (namespaceOutput string) {
	return fmt.Sprintf("%s.output", t.Name)
}

// FormatInput 从标准输入中获取数据
func (t Torm) FormatInput(data []byte) (input string) {
	inTransfers, _ := t.Transfers.SplitInOut(t.Name)
	gjsonPath := inTransfers.Reverse().String()
	input = gjson.GetBytes(data, gjsonPath).String()
	input = gjson.Get(input, t.NamespaceInput()).String()
	return input
}

// FormatOutput 格式化输出标准数据
func (t Torm) FormatOutput(data []byte) (output []byte) {
	_, outTransfers := t.Transfers.SplitInOut(t.Name)
	gjsonPath := outTransfers.String()
	output = []byte(gjson.GetBytes(data, gjsonPath).String())
	output, err := sjson.SetRawBytes(output, t.NamespaceOutput(), output)
	if err != nil {
		panic(err)
	}
	return output
}

//解析tpl 文本，生成 Torms
func ParserTpl(source *Source, tplText string, pathtransferLine pathtransfer.TransferLine, flow packethandler.Flow, packetHandlers packethandler.PacketHandlers) (torms Torms, err error) {
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

	transfers := pathtransferLine.Transfer()
	torms = make(Torms, 0)

	for _, tpl := range t.Templates() {
		tplName := tpl.Name()
		if tplName == "" {
			continue
		}
		torm := &Torm{
			Name:           tplName,
			Source:         *source,
			TplText:        tpl.Root.String(),
			Transfers:      transfers,
			PacketHandlers: packetHandlers,
			Flow:           flow,
			template:       tpl,
		}
		torms.Add(*torm)
	}
	return torms, nil
}

func (ts *Torms) Add(subTorms ...Torm) {
	if *ts == nil {
		*ts = make(Torms, 0)
	}
	*ts = append(*ts, subTorms...)
}
func (ts Torms) GetByName(name string) (t *Torm, err error) {
	for _, t := range ts {
		if strings.EqualFold(name, t.Name) {
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
