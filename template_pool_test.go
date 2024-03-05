package torm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/torm"
)

type ServerGetByIdEntity struct {
	Id     int //主键
	IDList []int
	torm.VolumeMap
}

func (t *ServerGetByIdEntity) TplName() string {
	return "ServerGetById"
}

func (t *ServerGetByIdEntity) Torm() string {
	return "{{define \"ServerGetById\"}}\nselect * from `server`  where `id`=:Id  and `id` in ({{in . .IDList}}) and `deleted_at` is null;\n{{end}}\n\n\n"
}

func (t *ServerGetByIdEntity) Namespace() string {
	return "curdservice"
}

type ServerGetByNameEntity struct {
	Name        string //项目标识
	ServiceName string //服务名称
	torm.VolumeMap
}

func (t *ServerGetByNameEntity) TplName() string {
	return "serverGetByName"
}

func (t *ServerGetByNameEntity) Torm() string {
	return "{{define \"serverGetByName\"}}\nselect * from `server` where `name`=:Name and `service_name`=:ServiceName and `deleted_at` is null;\n{{end}}\n"
}

func (t *ServerGetByNameEntity) Namespace() string {
	return "curdservice"
}

type ServicePaginateEntity struct {
	Limit                      int //
	Offset                     int //
	ServicePaginateWhereEntity     //
	torm.VolumeMap
}

func (t *ServicePaginateEntity) TplName() string {
	return "ServicePaginate"
}

func (t *ServicePaginateEntity) Torm() string {
	return "{{define \"ServicePaginate\"}}\nselect * from `service`  where 1=1 {{template \"ServicePaginateWhere\" .}}   and `deleted_at` is null {{template \"ServicePaginateOrder\"}}  limit :Offset,:Limit ;\n{{end}}\n\n\n\n"
}

func (t *ServicePaginateEntity) Namespace() string {
	return "curdservice"
}

type ServicePaginateWhereEntity struct {
	Name string //项目标识
	torm.VolumeMap
}

func (t *ServicePaginateWhereEntity) TplName() string {
	return "ServicePaginateWhere"
}

func (t *ServicePaginateWhereEntity) Torm() string {
	return "{{define \"ServicePaginateWhere\"}}\n{{noEmpty \"and `name` like \\\"%%%s%%\\\"\" .Name }}\n{{end}}\n\n"
}

func (t *ServicePaginateWhereEntity) Namespace() string {
	return "curdservice"
}

type ServicePaginateOrderEntity struct {
	torm.VolumeMap
}

func (t *ServicePaginateOrderEntity) TplName() string {
	return "ServicePaginateOrder"
}

func (t *ServicePaginateOrderEntity) TplType() string {
	return "text"
}

func (t *ServicePaginateOrderEntity) Torm() string {
	return "{{define \"ServicePaginateOrder\"}}\n    order by `id` desc \n{{end}}\n"
}

func (t *ServicePaginateOrderEntity) Namespace() string {
	return "curdservice"
}

func TestRegisterTorm(t *testing.T) {
	byId := &ServerGetByIdEntity{
		Id:     1,
		IDList: []int{1},
	}
	err := torm.RegisterTorm(byId)
	require.NoError(t, err)
	byName := new(ServerGetByNameEntity)
	err = torm.RegisterTorm(byName)
	require.NoError(t, err)
	err = torm.RegisterTorm(new(ServicePaginateEntity))
	require.NoError(t, err)

	err = torm.RegisterTorm(new(ServicePaginateWhereEntity))
	require.NoError(t, err)
	err = torm.RegisterTorm(new(ServicePaginateOrderEntity))
	require.NoError(t, err)

	sql, _, _, err := torm.GetSQL(byId.Namespace(), byId.TplName(), byId)
	require.NoError(t, err)
	assert.Equal(t, "select * from `server` where `id`=1 and `id` in (1) and `deleted_at` is null;", sql)

	byName = &ServerGetByNameEntity{
		Name: "a",
	}
	sql, _, _, err = torm.GetSQL(byName.Namespace(), byName.TplName(), byName)
	require.NoError(t, err)
	require.Equal(t, "select * from `server` where `name`='a' and `service_name`='' and `deleted_at` is null;", sql)

	page := &ServicePaginateEntity{
		Limit: 10,
	}
	sql, _, _, err = torm.GetSQL(page.Namespace(), page.TplName(), page)
	require.NoError(t, err)
	require.Equal(t, "select * from `service` where 1=1 and `deleted_at` is null order by `id` desc limit 0,10 ;", sql)

}
