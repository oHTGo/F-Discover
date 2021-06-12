package logger

import (
	"time"

	"github.com/kataras/iris/v12"
)

func Handle() iris.Handler {
	return func(ctx iris.Context) {
		// all except latency to string
		var ip string
		var latency time.Duration
		var startTime, endTime time.Time
		startTime = time.Now()

		// Before Next.
		ip = ctx.RemoteAddr()

		ctx.Next()

		// no time.Since in order to format it well after
		endTime = time.Now()
		latency = endTime.Sub(startTime)

		Info("Access log", "IP: "+ip)
		Info("Access log", "Latency: "+latency.String())
	}
}
