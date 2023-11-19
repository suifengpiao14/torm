package tormstream

import "github.com/suifengpiao14/logchan/v2"

type LogName string

func (l LogName) String() string {
	return string(l)
}

const (
	LOG_INFO_SQL LogName = "LogInfoSQL"
)

type LogInfoToSQL struct {
	SQL       string                 `json:"sql"`
	Named     string                 `json:"named"`
	NamedData map[string]interface{} `json:"namedData"`
	Data      interface{}            `json:"data"`
	Err       error                  `json:"error"`
	Level     string                 `json:"level"`
	logchan.EmptyLogInfo
}

func (l *LogInfoToSQL) GetName() logchan.LogName {
	return LOG_INFO_SQL
}
func (l *LogInfoToSQL) Error() error {
	return l.Err
}
func (l *LogInfoToSQL) GetLevel() string {
	return l.Level
}

const (
	LOG_INFO_EXEC_TEMPLATE LogName = "LogInfoExecTemplate"
)

type LogInfoExecTpl struct {
	TplName  string          `json:"tplName"`
	Volume   VolumeInterface `json:"volumne"`
	NamedSQL string          `json:"namedSql"`
	Err      error           `json:"error"`
	Level    string          `json:"level"`
	logchan.EmptyLogInfo
}

func (l *LogInfoExecTpl) GetName() logchan.LogName {
	return LOG_INFO_EXEC_TEMPLATE
}
func (l *LogInfoExecTpl) Error() error {
	return l.Err
}
func (l *LogInfoExecTpl) GetLevel() string {
	return l.Level
}
