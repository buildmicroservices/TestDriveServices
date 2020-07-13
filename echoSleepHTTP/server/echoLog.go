package server

// HTTP logging specific utilities


import (

	"net/http"
	"log"
)

// Utility function to dump the Request headers to log stream
func dumpHeaders(r *http.Request) {
	for name, headers := range r.Header {
		for _, h := range headers {
			log.Println("%v: %v\n", name, h)
		}
	}
}

