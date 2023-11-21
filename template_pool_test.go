package tormstream_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/tormstream"
)

type ServerGetByIdEntity struct {
	Id     int //主键
	IDList []int
	tormstream.VolumeMap
}

func (t *ServerGetByIdEntity) TplName() string {
	return "ServerGetById"
}

func (t *ServerGetByIdEntity) Torm() string {
	return "{{define \"ServerGetById\"}}\nselect * from `server`  where `id`=:Id  and `id` in ({{in . .IDList}}) and `deleted_at` is null;\n{{end}}\n\n\n"
}

func (t *ServerGetByIdEntity) Identity() string {
	return "curdservice"
}

type ServerGetByNameEntity struct {
	Name        string //项目标识
	ServiceName string //服务名称
	tormstream.VolumeMap
}

func (t *ServerGetByNameEntity) TplName() string {
	return "serverGetByName"
}

func (t *ServerGetByNameEntity) Torm() string {
	return "{{define \"serverGetByName\"}}\nselect * from `server` where `name`=:Name and `service_name`=:ServiceName and `deleted_at` is null;\n{{end}}\n"
}

func (t *ServerGetByNameEntity) Identity() string {
	return "curdservice"
}

func TestRegisterTorm(t *testing.T) {
	byId := &ServerGetByIdEntity{
		Id:     1,
		IDList: []int{1},
	}
	err := tormstream.RegisterTorm(byId)
	require.NoError(t, err)
	byName := new(ServerGetByNameEntity)
	err = tormstream.RegisterTorm(byName)
	require.NoError(t, err)
	sql, _, _, err := tormstream.GetSQL(byId.Identity(), byId.TplName(), byId)
	require.NoError(t, err)
	assert.Equal(t, "select * from `server` where `id`=1 and `id` in (1) and `deleted_at` is null;", sql)
	byName = &ServerGetByNameEntity{
		Name: "a",
	}

	sql, _, _, err = tormstream.GetSQL(byName.Identity(), byName.TplName(), byName)
	require.NoError(t, err)
	fmt.Println(sql)
}
