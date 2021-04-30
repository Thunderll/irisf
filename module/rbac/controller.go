package rbac

import (
	"iris_project_foundation/common"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type RoleController struct {
	Ctx     iris.Context
	Service *RoleService
}

func (c *RoleController) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Register(NewRoleService)
	b.Handle("POST", "/set-permissions", "SetPermissions")
	b.Handle("POST", "/set-user", "SetUser")
	b.Handle("GET", "/user-roles/{id:int64}", "GetUserRolesBy")
}

func (c *RoleController) Post() *common.Response {
	var err error

	if err = c.Service.CreateRole(); err != nil {
		return common.Failed(err)
	}
	return common.Success(nil, 0, 0)
}

func (c *RoleController) Delete(id int64) *common.Response {
	var err error

	if err = c.Service.DeleteRole(id); err != nil {
		return common.Failed(err)
	}
	return common.Success(nil, 0, 0)
}

func (c *RoleController) SetPermissions() *common.Response {
	var err error

	if err = c.Service.SetPermissionsForRole(); err != nil {
		return common.Failed(err)
	}
	return common.Success(nil, 0, 0)
}

func (c *RoleController) SetUser() *common.Response {
	var err error

	if err = c.Service.SetRoleForUser(); err != nil {
		return common.Failed(err)
	}
	return common.Success(nil, 0, 0)
}

func (c *RoleController) Get() *common.Response {
	var (
		err   error
		roles []*Role
	)

	if roles, err = c.Service.GetRoleList(); err != nil {
		return common.Failed(err)
	}

	return common.Success(roles, 0, 0)
}

func (c *RoleController) GetBy(id int64) *common.Response {
	var (
		err    error
		result *RolePerms
	)

	if result, err = c.Service.GetRoleDetail(id); err != nil {
		return common.Failed(err)
	}

	return common.Success(result, 0, 0)
}

func (c *RoleController) GetPermissions() *common.Response {
	var (
		err   error
		perms []*Permission
	)

	if perms, err = c.Service.GetPermissions(); err != nil {
		return common.Failed(err)
	}
	return common.Success(perms, 0, 0)
}

func (c *RoleController) GetUserRolesBy(id int64) *common.Response {
	var (
		err    error
		result *UserRoles
	)

	if result, err = c.Service.GetUserRole(id); err != nil {
		return common.Failed(err)
	}
	return common.Success(result, 0, 0)
}
