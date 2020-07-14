package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"

	//	"github.com/gorilla/mux"
)

// Define the Response Model for the echo endpoint.  Echo method, URL, request headers, query parameters, and request body
type EchoResponseHTTP struct {
	Method      string              `json:"httpMethod"`
	Url         *url.URL            `json:"httpURL"`
	Metadata    http.Header         `json:"httpHeaders"`
	QueryParams map[string][]string `json:"httpQueryParams"`
	Body        string              `json:"httpBody"`
}

type Service interface {
}


// the global server context
type echoSleepCtx struct {
	ServiceName     string
	// Note, handler chain could be a httprouter
	EchoSleepHandleChain http.Handler
}

// per request context
// Sleeper: propagate the sleep duration
// Ctx: holds any service specific context
// TimeSpan: trace and timer telemetry
// EchoResponseHTTP: include response in the context object
type requestCtx struct {
	Sleeper          time.Duration
	Ctx              echoSleepCtx
	TimeSpan TimeSpanner
	EchoResponseHTTP EchoResponseHTTP
}



// standard interface to handle request
func (ctx echoSleepCtx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx.echoSleepHTTP(w, r)
}

// main HTTP Request Handler   format http.Handler
func (ctx echoSleepCtx) echoSleepHTTP(w http.ResponseWriter, r *http.Request) {

	// trace level message
	// log.Trace("{ message: \"received echo request\" } ")
	externalId := r.Header.Get("externalId")
	// establish per request context
	requestCtx := &requestCtx{
		Sleeper: getSleepDuration(r),
		Ctx:     ctx,
		TimeSpan: NewTimeSpanner(ctx.ServiceName+"-echo",externalId),
		EchoResponseHTTP: EchoResponseHTTP{
			Method:      r.Method,
			Url:         r.URL,
			QueryParams: r.URL.Query(),
			Metadata:    r.Header,
			// Body: nil,
		},
	}

	// if header set, sleep before echo
	gotoSleep(requestCtx)

	if r.Method == "PUT" || r.Method == "POST" {
		// if PUT or POST, ensure we drain the input stream
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err == nil {
			// TODO: protect from DOS and check content-length, cap body[:?] slice
			// attach the body to the response object
			requestCtx.EchoResponseHTTP.Body = string(body[:])
			// log.Debug("body", body)
		} else {
			log.Println("{ errorType: \"bodyReadError\", error: \"" + err.Error() + "\"}")
		}
	}

	// end of the line handler.... can start writing out
	rec := NewHttpRecorder(w)
	rec.SetResponseHeaders()
	rec.ResponseWriter.WriteHeader(223)

	b, err := json.Marshal(requestCtx.EchoResponseHTTP)
	if err == nil {
		w.Write(b)
	} else {
		log.Println("jsonMarshalError: ", err)
	}
}

func (ctx echoSleepCtx) echo1() http.Handler {
	log.Println("{ service: \"" + ctx.ServiceName + "\", start: true, message:\"registering echo handler\"}")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx.echoSleepHTTP(w, r)
	},
	)
}

func initializeRouter() *httprouter.Router {

	// inject the service global context values
	echoSleepCtx := &echoSleepCtx{ServiceName: "echoSleepHTTP"}

	echoSleepCtx.EchoSleepHandleChain = echoSleepCtx.echo1()

	router := httprouter.New()

	// SERVE STATIC FILES...
	//http.Handle("/static/", http.StripPrefix("/static/", fs))
	router.ServeFiles("/static/*filepath", http.Dir("/var/www/public/"))

	// EStablish the echoSleep Handler off the root URI basepath
	var handle httprouter.Handle

	//	router.Handle(route.Method, route.Path, handle)
	// inject the logger BEFORE the echoSleepHTTP handler
	handle = loggingMiddleware(echoSleepCtx.echo1())
	router.GET("/echo1", handle)

	//router.Handle("GET", "/{rest:.*}", handle)
	//router.GET("/{rest:.*}", echoSleepCtx.echoSleepHTTP)
	//router.PUT("/{rest:.*}", handle)

	router.GET("/echo", StdToJulienMiddleware(echoSleepCtx.EchoSleepHandleChain))

	router.GET("/push", pushHandle)
	router.GET("/pull", pullHandle)

	//http.Handle("/", router)
	//server.Handler = router.Handler()
	//router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
	//	return true
	// r.ProtoMajor == 0
	//})
	// add our default logger
	//router.Use(loggingMiddleware)

	// establish default not available handler for not allowed use case
	router.HandleMethodNotAllowed = true
	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	//	router.MethodNotAllowed = func() ()

	return router
}
