package server
// HTTP Utilities
// @author: Chris Haddad

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// Standard Error Response Object 
type ErrorResponse struct {
	ErrorCode    string     `json:"errorCode"`
	ErrorMessage string     `json:"errorMessage"`
	ErrorLine    string     `json:"errorLine"`
	ErrorDetail  [][]string `json:"errorDetail"`
}

// HTTPRecorder structure to handle Response Writing and Response context
type HttpRecorder struct {
	http.ResponseWriter
	status        int
	errorResponse ErrorResponse
	// add Body ?
}

// utility function to create a new HTTP recorder that wraps ResponseWriter
func NewHttpRecorder(w http.ResponseWriter) HttpRecorder {
	return HttpRecorder{w, 0, ErrorResponse{}}
}

// Utility function to set standard response headers
func (rec *HttpRecorder) SetResponseHeaders() {
	rec.ResponseWriter.Header().Set("content-type", "application/json")
	var err error
	uid := uuid.Must(uuid.New(), err).String()
	rec.ResponseWriter.Header().Set("X-serverCorrelation", uid)
}

// Utility function to write out headers
func (rec *HttpRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

//TODO: pass Http context structure down the middleware chain
// standard HTTP middleware context - per request
type HttpCtx struct {
	ServiceName string
	TimeSpanner TimeSpanner
	Logger      log.Logger
}

//TODO: add Prometheus metrics... in other file

// Middleware handler to emit standard request logging information
// Place in the chain AFTER security principle decoding
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("{ handler: \"loggingMiddleware\" path: \"" + r.RequestURI + "\"}")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

// Span trace middleware, establish call chain and timings
func traceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("{ handler: \"traceMiddleware\", path: \"" + r.RequestURI + "\"}")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		if next != nil {
			externalId := r.Header.Get("externalId")
			timeSpan := NewTimeSpanner("echo", externalId)
			span, _ := timeSpan.addTimeSpan("domainCall")

			//TODO: pass the time spanner across sub-calls.

			span.StartTimer()
			next.ServeHTTP(w, r)
			span.StopTimer()

			b, err := json.Marshal(timeSpan)
				if err == nil {
					log.Println("timespanner:"+string(b))
				}

				// TODO: emit the internal call timer to telemetry service
				// TODO: add openTelemetry
		}
	})
}

func pushHandle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("in the push handler ")
	rec := NewHttpRecorder(w)
	rec.SetResponseHeaders()
	rec.ResponseWriter.WriteHeader(201)
	// TODO: log URI parameters
	w.Write([]byte("{ method: pushHandler, time: 100, uri=\"/push\"}"))
}

func pullHandle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("in the pull handler ")

	rec := NewHttpRecorder(w)
	rec.SetResponseHeaders()
	rec.ResponseWriter.WriteHeader(201)
	// TODO: log URI parameters
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
