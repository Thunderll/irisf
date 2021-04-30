package throttler

import (
	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	"github.com/kataras/iris/v12"
)

func Serve(l *limiter.Limiter) iris.Handler {
	return func(ctx iris.Context) {
		httpError := tollbooth.LimitByRequest(l, ctx.ResponseWriter(), ctx.Request())
		if httpError != nil {
			ctx.ContentType(l.GetMessageContentType())
			ctx.StatusCode(httpError.StatusCode)
			ctx.WriteString(httpError.Message)
			ctx.StopExecution()
			return
		}

		ctx.Next()
	}
}
