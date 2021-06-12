package errors

import (
	"errors"
	"f-discover/interfaces"
	"f-discover/logger"
	"fmt"
	"net/http/httputil"
	"runtime"
	"strings"

	"github.com/kataras/iris/v12/context"
)

// New returns a new recover middleware,
// it recovers from panics and logs
// the panic message to the application's logger "Warn" level.
func Handle() context.Handler {
	return func(ctx *context.Context) {
		defer func() {
			if err := recover(); err != nil {
				if ctx.IsStopped() { // handled by other middleware.
					return
				}

				logger.Error("errors.Handler", fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.HandlerName()), nil)
				{
					rawReq, _ := httputil.DumpRequest(ctx.Request(), true)
					for _, requestData := range strings.Split(string(rawReq), "\n") {
						logger.Error("errors.Handler", "Request", errors.New(requestData))
					}
				}
				for i := 1; ; i++ {
					_, file, line, got := runtime.Caller(i)
					if !got {
						break
					}
					logger.Error("errors.Handler", "Callers", errors.New(fmt.Sprintf("%s:%d", file, line)))
				}

				ctx.StopWithJSON(500, interfaces.IFail{
					Message: "Something's wrong",
				})
			}
		}()

		ctx.Next()
	}
}
