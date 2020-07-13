package server


import (
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strings"
)


  
func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	sayTo := params.ByName("name")
	if ( strings.Compare(sayTo,"") == 0) {
	sayTo = "World"
}
io.WriteString(w, "Hello "+ sayTo + "!\n")
}


// Middleware without "github.com/julienschmidt/httprouter"
func StdToStdMiddleware(next http.Handler) http.Handler {

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // do stuff
        next.ServeHTTP(w, r)
    })
}

// Middleware for a standard handler returning a "github.com/julienschmidt/httprouter" Handle
func StdToJulienMiddleware(next http.Handler) httprouter.Handle {

    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        // do stuff
        next.ServeHTTP(w, r)
    }
}

// Pure "github.com/julienschmidt/httprouter" middleware
//func JulienToStdMiddleware(next httprouter.Handle) http.Handler {
//
//	return func(w http.ResponseWriter, r *http.Request) {
//		next.
//	}
//}


// Pure "github.com/julienschmidt/httprouter" middleware
func JulienToJulienMiddleware(next httprouter.Handle) httprouter.Handle {

    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        // do stuff
        next(w, r, ps)
    }
}

func JulienHandler() httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        // do stuff
    }
}

func StdHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // do stuff
    })
}

