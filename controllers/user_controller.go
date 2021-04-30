package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"iris_project_foundation/common"
	"iris_project_foundation/services"
)

type UserController struct {
	Ctx     iris.Context
	Service services.IUserService
}

func (c *UserController) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Register(services.NewUserService)
}

// Post 增
func (c *UserController) Post() *common.Response {
	var (
		err error
	)

	if err = c.Service.CreateUser(); err != nil {
		return common.Failed(err)
	}
	return common.Success(nil, 0, 0)
}

// Get 查, 列表
func (c *UserController) Get() *common.Response {
	var (
		err  error
		resp *common.Response
	)

	if resp, err = c.Service.GetUserList(); err != nil {
		return common.Failed(err)
	}

	return resp
}

// GetBy 查, 详情
func (c *UserController) GetBy(id int64) *common.Response {
	var (
		err  error
		resp *common.Response
	)
	if resp, err = c.Service.GetUser(id); err != nil {
		return common.Failed(err)
	}
	return resp
}
