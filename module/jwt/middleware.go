package jwt

import (
	"github.com/kataras/iris/v12/context"
	"iris_project_foundation/common"
	"iris_project_foundation/common/api_error"
	"strings"
)

// ErrorHandler 处理认证失败
func ErrorHandler(ctx *context.Context, e error) {
	ctx.JSON(common.Failed(api_error.UnauthorizedError))
}

func FromHeader(ctx *context.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// pure check: authorization header format must be Bearer {token}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "jwt" {
		return ""
	}

	return authHeaderParts[1]
}

// Serve 构造JWT中间件
func Serve() context.Handler {
	var (
		verifyMiddleware context.Handler
	)

	verifier.ErrorHandler = ErrorHandler
	//verifier.Extractors = []jwt.TokenExtractor{FromHeader}
	verifier.DisableContextUser = false

	verifyMiddleware = verifier.Verify(func() interface{} {
		return new(UserClaims)
	})
	return verifyMiddleware
}
