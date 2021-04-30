package rotate_logger

import (
	"github.com/kataras/iris/v12"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"iris_project_foundation/config"
	"log"
	"time"
)

func InitRotateLogger(app *iris.Application) {
	var (
		accessLogger *rotatelogs.RotateLogs
		errorLogger  *rotatelogs.RotateLogs
	)
	// 配置日志分割
	accessLogger = BuildAccessLogger()
	errorLogger = BuildErrorLogger()
	app.Logger().SetLevelOutput("info", accessLogger)
	app.Logger().SetLevelOutput("warn", errorLogger)
	app.Logger().SetLevelOutput("error", errorLogger)
}

func BuildAccessLogger() *rotatelogs.RotateLogs {
	var (
		err    error
		logger *rotatelogs.RotateLogs
	)
	if logger, err = rotatelogs.New(
		config.GConfig.Server.AccessLog+"/access_log.%Y%m%d",
		rotatelogs.WithLinkName(config.GConfig.Server.AccessLog+"/access_log"),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationCount(30),
	); err != nil {
		log.Fatalf("[ERR] 初始化访问日志失败!\n%v", err)
	}
	return logger
}

func BuildErrorLogger() *rotatelogs.RotateLogs {
	var (
		err    error
		logger *rotatelogs.RotateLogs
	)
	if logger, err = rotatelogs.New(
		config.GConfig.Server.ErrorLog+"/error_log.%Y%m%d",
		rotatelogs.WithLinkName(config.GConfig.Server.ErrorLog+"/error_log"),
		rotatelogs.WithRotationTime(720*time.Hour),
		rotatelogs.WithRotationCount(13),
	); err != nil {
		log.Fatalf("[ERR] 初始化访问日志失败!\n%v", err)
	}
	return logger
}
