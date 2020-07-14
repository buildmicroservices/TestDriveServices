package server


import (
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/google/uuid"

)

// Standard Error Response Object 
type ErrorResponse struct {
	ErrorCode	string	`json:"errorCode"`
	ErrorMessage	string	`json:"errorMessage"`
	ErrorLine	string	`json:"errorLine"`
	ErrorDetail	[][]string	`json:"errorDetail"`
  }


type HttpRecorder struct {
	http.ResponseWriter
	status int
	errorResponse ErrorResponse
}

func NewHttpRecorder(w http.ResponseWriter) HttpRecorder {
   return HttpRecorder { w , 0, ErrorResponse{} }
}

// set standard response headers
func (rec *HttpRecorder )SetResponseHeaders() {
	rec.ResponseWriter.Header().Set("content-type", "application/json")
	var err error
	uid := uuid.Must(uuid.New(), err).String()
	rec.ResponseWriter.Header().Set("X-serverCorrelation", uid)
}

func (rec *HttpRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

//TODO: test log middleware

//TODO: test timer middleware

//TODO: add Prometheus metrics... in other file

// Middleware handler to emit standard request logging information
// Place in the chain AFTER security principle decoding
func loggingMiddleware(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Do stuff here
		log.Println("{ path: \"" + r.RequestURI + "\"}")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		if next != nil {
			next.ServeHTTP(w, r)
		}
	}
}

func pushHandle (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("in the push handler ")
	rec := NewHttpRecorder(w)
	rec.SetResponseHeaders()
	rec.ResponseWriter.WriteHeader(201)
	w.Write([]byte("{ method: pushHandler, time: 100, uri=\"/push\"}"))
}

func pullHandle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("in the pull handler ")

	rec := NewHttpRecorder(w)
	rec.SetResponseHeaders()
	rec.ResponseWriter.WriteHeader(201)
	w.Write([]byte("{ method: pullHandler, time: 100, uri=\"/pull\"}"))
}

/// STANDARD API RESULT FRAMEWORK

// Writes the response as a standard JSON response with StatusOK
//func writeOKResponse(w http.ResponseWriter, m interface{}) {
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(http.StatusOK)
//	if err := json.NewEncoder(w).Encode(&JsonResponse{Data: m}); err != nil {
//		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
//	}
//}

// Writes the error response as a Standard API JSON response with a response code
//func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(errorCode)
//	json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: errorCode, Title: errorMsg}})
//}

/// EXAMPLE ROUTE SETUP

/*
Define all the routes here.
A new Route entry passed to the routes slice will be automatically
translated to a handler with the NewRouter() function
*/
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

/*
type Routes []Route

func AllRoutes() Routes {
	routes := Routes{
		Route{"Index", "GET", "/", echoSleepHTTPHandle },
		Route{"BookIndex", "GET", "/books", echoSleepHTTPHandle },
		Route{"Bookshow", "GET", "/books/:isdn", echoSleepHTTPHandle },
		Route{"Bookshow", "POST", "/books", echoSleepHTTPHandle },
	}
	return routes
}
*/