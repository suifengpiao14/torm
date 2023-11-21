package tormstream

import "github.com/suifengpiao14/logchan/v2"

type _LogName string

func (l _LogName) String() string {
	return string(l)
}

const (
	LOG_INFO_SQL _LogName = "LogInfoSQL"
)

type LogInfoToSQL struct {
	TplIdentify string           `json:"tplIdentify"`
	TplName     string           `json:"tplName"`
	InputVolume VolumeInterface `json:"inputVolume"`

	SQL          string                 `json:"sql"`
	Named        string                 `json:"named"`
	NamedData    map[string]interface{} `json:"namedData"`
	TPLOutVolume VolumeInterface       `json:"tplOutVolume"`
	Err          error                  `json:"error"`
	Level        string                 `json:"level"`
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
	LOG_INFO_EXEC_TEMPLATE _LogName = "LogInfoExecTemplate"
)
