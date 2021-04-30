package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type serverConfig struct {
	ServerUrl  string
	ServerPort int64
	TimeFormat string
	Charset    string
	AccessLog  string
	ErrorLog   string
}

type appConfig struct {
	Secret                 string
	LogLevel               string
	DefaultPageSize        int64
	TokenPair              bool
	AccessTokenExpiration  int64
	RefreshTokenExpiration int64
	Blocklist              bool
	BlocklistPrefix        string
	ScheduleInterval       int64
	OrderLockPrefix        string
}

type wechatConfig struct {
	Code2SessionAPI string
	WechatAppID     string
	WechatSecret    string
}

type databaseConfig struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
}

type redisConfig struct {
	Addr     string
	Password string
	DB       int64
}

type Config struct {
	Server   *serverConfig
	App      *appConfig
	Wechat   *wechatConfig
	Database *databaseConfig
	Redis    *redisConfig
}

var GConfig *Config

// InitConfig 初始化配置文件
func InitConfig() {
	var (
		err      error
		filePath string
		config   Config
	)

	filePath = "./config/config.tml"
	if _, err = toml.DecodeFile(filePath, &config); err != nil {
		log.Fatalf("[ERR] 配置文件加载失败!\n %v", err)
	}
	GConfig = &config
}
