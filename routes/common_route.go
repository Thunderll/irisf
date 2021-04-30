package routes

import (
	"iris_project_foundation/controllers"
	"iris_project_foundation/module/jwt"
	"iris_project_foundation/module/rbac"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func CommonRouteSetup(r iris.Party) {
	var (
		v1 iris.Party
	)
	v1 = r.Party("/v1")

	// token生成接口
	v1.Post("/web-obtain-token", jwt.WebObtainToken)
	mvc.Configure(v1, noAuthWebAPI)

	// jwt认证中间件
	v1.Use(jwt.Serve())

	// token注销接口
	v1.Get("/logout", jwt.Logout)
	mvc.Configure(v1, authWebAPI)
}

// authWebAPI 注册需要认证的接口和中间件
func authWebAPI(app *mvc.Application) {
	// 添加casbin权限管理中间件
	app.Router.Use(rbac.Serve())

	mvc.New(app.Router.Party("/web-user")).Handle(new(controllers.UserController))
}

// noAuthWebAPI 注册不需要认证的接口和中间件
func noAuthWebAPI(app *mvc.Application) {
	mvc.New(app.Router.Party("/role")).Handle(new(rbac.RoleController))
}
