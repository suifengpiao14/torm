package torm

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/sqlexec"
	"github.com/suifengpiao14/sqlexec/sqlexecparser"
	"github.com/suifengpiao14/sshmysql"
	"github.com/suifengpiao14/torm/sourceprovider"
)

type ProviderI interface {
	TypeName() string
}

type Source struct {
	Identifer string    `json:"identifer"`
	Type      string    `json:"type"`
	Config    string    `json:"config"`
	Provider  ProviderI `json:"-"`
	DDL       string    `json:"ddl"`
	SSHConfig string    `json:"sshConfig"`
}

const (
	SOURCE_TYPE_SQL_MEMORY = "SQL_MEMORY"
	SOURCE_TYPE_SQL        = "SQL"
	SOURCE_TYPE_CURL       = "CURL"
	SOURCE_TYPE_BIN        = "BIN"
	SOURCE_TYPE_REDIS      = "REDIS"
	SOURCE_TYPE_RABBITMQ   = "RABBITMQ"
)

var ERROR_SOURCE_NOT_FOUND = errors.New("not found source")

// Init 配置Provider,DDL等
func (s *Source) Init() (err error) {
	switch strings.ToUpper(s.Type) {
	case SOURCE_TYPE_SQL:
		dbConfig, err := sqlexec.JsonToDBConfig(s.Config)
		if err != nil {
			return err
		}
		sshConfig, err := sshmysql.JsonToSSHConfig(s.SSHConfig)
		if errors.Is(err, sqlexec.ERROR_EMPTY_CONFIG) {
			err = nil
		}
		if err != nil {
			return err
		}
		dbProvider := sourceprovider.NewDBProvider(*dbConfig, sshConfig)
		if s.Provider == nil {
			s.Provider = dbProvider
		}
		db := dbProvider.GetDB()
		if s.DDL == "" {
			s.DDL, err = sqlexec.GetDDL(db)
			if err != nil {
				err = errors.WithMessagef(err, "config:%s", s.Config)
				return err
			}
		}
		if s.DDL != "" { // 注册关联表结构
			err = sqlexecparser.RegisterTableByDDL(s.DDL)
			if err != nil {
				return err
			}
		}
		//todo curl , bin 提供者实现
	default:
		err = errors.Errorf("not impliment provider ;source type:%s", s.Type)
		return err
	}
	return nil
}

type Sources []Source

func (ss Sources) GetByIdentifer(identify string) (s *Source, err error) {
	for _, s := range ss {
		if strings.EqualFold(identify, s.Identifer) {
			return &s, nil
		}
	}
	err = errors.WithMessagef(ERROR_SOURCE_NOT_FOUND, "source identifer:%s", identify)
	return nil, err
}

// MakeSource 创建常规资源,方便外部统一调用
func MakeSource(identifer string, typ string, config string, sshConfig string, ddl string) (s Source, err error) {
	s = Source{
		Identifer: identifer,
		Type:      typ,
		Config:    config,
		DDL:       ddl,
		SSHConfig: sshConfig,
	}
	err = s.Init()
	if err != nil {
		return s, err
	}
	return s, nil
}
