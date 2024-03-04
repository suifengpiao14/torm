package torm

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/sqlexec"
	"github.com/suifengpiao14/sqlexec/sqlexecparser"
	"github.com/suifengpiao14/torm/sourceprovider"
)

type providerI interface {
	TypeName() string
}

type Source struct {
	Identifer string    `json:"identifer"`
	Type      string    `json:"type"`
	Config    string    `json:"config"`
	Provider  providerI `json:"-"`
	DDL       string    `json:"ddl"`
}

type Sources []Source

func (ss Sources) GetByIdentifer(identify string) (s *Source, err error) {
	for _, s := range ss {
		if strings.EqualFold(identify, s.Identifer) {
			return &s, nil
		}
	}
	err = errors.Errorf("not found source by source identifier: %s", identify)
	return nil, err
}

const (
	PROVIDER_SQL_MEMORY = "SQL_MEMORY"
	PROVIDER_SQL        = "SQL"
	PROVIDER_CURL       = "CURL"
	PROVIDER_BIN        = "BIN"
	PROVIDER_REDIS      = "REDIS"
	PROVIDER_RABBITMQ   = "RABBITMQ"
)

//MakeSource 创建常规资源,方便外部统一调用
func MakeSource(identifer string, typ string, config string, ddl string) (s Source, err error) {
	s = Source{
		Identifer: identifer,
		Type:      typ,
		Config:    config,
		DDL:       ddl,
	}
	var providerImp providerI
	switch s.Type {
	case PROVIDER_SQL:
		providerImp, err = sourceprovider.NewDBProvider(s.Config)
		if err != nil {
			return s, err
		}
		dbProvider, _ := providerImp.(*sourceprovider.DBProvider)
		db := dbProvider.GetDB()
		if s.DDL == "" {
			s.DDL, err = sqlexec.GetDDL(db)
			if err != nil {
				err = errors.WithMessagef(err, "config:%s", config)
				return s, err
			}
		}
		if s.DDL != "" { // 注册关联表结构
			err = sqlexecparser.RegisterTableByDDL(s.DDL)
			if err != nil {
				return s, err
			}
		}
		//todo curl , bin 提供者实现
	}
	s.Provider = providerImp
	return s, nil
}
