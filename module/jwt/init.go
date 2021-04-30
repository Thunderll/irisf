package jwt

import (
	"iris_project_foundation/config"
	"log"
	"time"

	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/kataras/iris/v12/middleware/jwt/blocklist/redis"
)

var (
	signer   *jwt.Signer
	verifier *jwt.Verifier
)

func InitJWT() {
	var (
		err       error
		blocklist *redis.Blocklist
	)
	signer = jwt.NewSigner(jwt.HS256, config.GConfig.Wechat.WechatSecret,
		time.Duration(config.GConfig.App.AccessTokenExpiration)*time.Minute)

	verifier = jwt.NewVerifier(jwt.HS256, config.GConfig.Wechat.WechatSecret)

	// 是否使用jwt黑名单
	if config.GConfig.App.Blocklist {
		if config.GConfig.Redis.Addr != "" {
			// 如果配置了redis, 则使用redis存储黑名单
			blocklist = redis.NewBlocklist()

			blocklist.ClientOptions.Addr = config.GConfig.Redis.Addr
			blocklist.ClientOptions.Password = config.GConfig.Redis.Password
			blocklist.ClientOptions.DB = int(config.GConfig.Redis.DB + 1)
			blocklist.Prefix = config.GConfig.App.BlocklistPrefix

			if err = blocklist.Connect(); err != nil {
				log.Fatalf("[ERR] JWT黑名单初始化失败!\n%v", err)
			}
			log.Println("[INFO] JWT黑名单初始化成功.")
			verifier.Blocklist = blocklist
		} else {
			// 使用默认黑名单存储,即存储在内存中
			verifier.WithDefaultBlocklist()
		}
	}

}
