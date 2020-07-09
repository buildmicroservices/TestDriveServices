package server

import (
	"net/http"
	"testing"
	"time"
)

func TestEchoSleepHTTP(t *testing.T) {
	//  want := "Hello, world."
	//  if got := Hello(); got != want {
	//      t.Errorf("Hello() = %q, want %q", got, want)
	//    }
	var duration time.Duration

	r := http.Request{}
	duration = getSleepDuration(&r)
	if duration.String() != "0s" {
		t.Errorf("no header check failed: %s", duration.String())
	}
	r.Header = http.Header{}
	r.Header.Set("X-Sleep", "2s")
	duration = getSleepDuration(&r)
	if duration.String() != "2s" {
		t.Errorf("header present check failed: %s", duration.String())
	}
	var requestCtx = requestCtx{Sleeper: duration}
	gotoSleep(&requestCtx)

	r.Header.Del("X-Sleep")
	r.Header.Set("X-Sleep", "5")
	duration = getSleepDuration(&r)
	if duration.String() != "0s" {
		t.Errorf("invalid sleep arg check failed: %s", duration.String())
	}

}
