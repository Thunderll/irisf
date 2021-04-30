package main

import (
	"fmt"
	"iris_project_foundation/config"
	"iris_project_foundation/middleware/cors"
	"iris_project_foundation/models"
	"iris_project_foundation/module/jwt"
	"iris_project_foundation/module/rbac"
	"iris_project_foundation/routes"
	"log"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
)

func newApp() *iris.Application {
	var (
		app *iris.Application
	)

	app = iris.New()

	// 加载配置文件
	config.InitConfig()

	// 日志配置
	//rotate_logger.InitRotateLogger(app)

	app.Logger().SetLevel(config.GConfig.App.LogLevel)
	app.Logger().Debug("Log level set to debug")
	app.UseRouter(logger.New())

	app.UseRouter(requestid.New())
	app.Logger().Debugf("Using <UUID4> to identify requests")
	app.UseRouter(recover.New())

	// 初始化跨域设置
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	app.UseRouter(crs)

	// 初始化数据库
	models.InitDatabase()

	// 初始化JWT
	jwt.InitJWT()

	// 初始化权限管理
	rbac.InitRBAC()

	// 路由注册
	app.PartyFunc("/web-api", routes.CommonRouteSetup)

	app.Configure(
		iris.WithOptimizations,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithTimeFormat(config.GConfig.Server.TimeFormat),
		iris.WithCharset(config.GConfig.Server.Charset),
	)
	return app
}

func main() {
	var (
		app *iris.Application
		err error
	)

	app = newApp()
	err = app.Listen(fmt.Sprintf("%s:%d", config.GConfig.Server.ServerUrl, config.GConfig.Server.ServerPort))
	if err != nil {
		log.Fatal(err)
	}
}
