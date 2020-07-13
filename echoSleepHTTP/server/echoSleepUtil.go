package server

// Sleep specific functions

import (
	"net/http"
	"time"
)

// given an HTTP request, retrieve sleep duration request header X-Sleep
// and save sleep value in request context
func getSleepDuration(r *http.Request) time.Duration {
	sleepTime := r.Header.Get("X-Sleep")
	if sleepTime != "" {
		sleeper, err1 := time.ParseDuration(sleepTime)
		if err1 == nil {
			return sleeper
		}
		//log.Warn("{ errmsg: \"invalid time ", err1, " \"}")
	}
	// return zero duration time object as time default
	// idea: make random within a configurable range (application level config. p90, p99, min, max, avg)
	t1 := time.Now()
	return t1.Sub(t1)
}

// Sleep for specified duration (request Context attribute contains sleep duration)
func gotoSleep(requestCtx *requestCtx) {
	if requestCtx.Sleeper != 0 {
		// sleepMessage := "{ service: \"" + requestCtx.Ctx.ServiceName + "\", sleep: " + requestCtx.Sleeper.String() + "}"
		// if sleep triggered, returns sleepMesage, else ""
		//log.Debug(sleepMessage)

		time.Sleep(requestCtx.Sleeper)
	}
}
