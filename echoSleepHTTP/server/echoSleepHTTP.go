package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

        "github.com/julienschmidt/httprouter"

//	"github.com/gorilla/mux"

	uuid "github.com/satori/go.uuid"
)

// Define the Response Model for the echo endpoint.  Echo method, URL, request headers, query parameters, and request body 
type EchoResponseHTTP struct {
	Method      string              `json:"httpMethod"`
	Url         *url.URL            `json:"httpURL"`
	Metadata    http.Header         `json:"httpHeaders"`
	QueryParams map[string][]string `json:"httpQueryParams"`
	Body        string              `json:"httpBody"`
}

// the global server context
type echoSleepCtx struct {
	ServiceName string
}

// per request context
// Sleeper: propagate the sleep duration
// Ctx: holds any service specific context
// EchoResponseHTTP: include response in the context object
type requestCtx struct {
	Sleeper          time.Duration
	Ctx              echoSleepCtx
	EchoResponseHTTP EchoResponseHTTP
}

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

// Sleep for specified duration (request Context attribute)
func gotoSleep(requestCtx *requestCtx) {
	if requestCtx.Sleeper != 0 {
		// sleepMessage := "{ service: \"" + requestCtx.Ctx.ServiceName + "\", sleep: " + requestCtx.Sleeper.String() + "}"
                // if sleep triggered, returns sleepMesage, else ""
		//log.Debug(sleepMessage)

		time.Sleep(requestCtx.Sleeper)
	}
}

type Service interface {
}

// Dump the Request headers to log file
func dumpHeaders(r *http.Request) {
	for name, headers := range r.Header {
		for _, h := range headers {
			log.Println("%v: %v\n", name, h)
		}
	}
}

// main HTTP Request Handler
func (ctx echoSleepCtx) echoSleepHTTP(w http.ResponseWriter, r *http.Request) {

        // trace level message 
	// log.Trace("{ message: \"received echo request\" } ")

	// establish per request context
	requestCtx := &requestCtx{
		Sleeper: getSleepDuration(r),
		Ctx:     ctx,
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

	setResponseHeaders(w)

	b, err := json.Marshal(requestCtx.EchoResponseHTTP)
	if err == nil {
		w.Write(b)
	} else {
		log.Println("jsonMarshalError: ", err)
	}
}

// set standard response headers
func setResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
	var err error
	uid := uuid.Must(uuid.NewV4(), err).String()
	w.Header().Set("X-serverCorrelation", uid)
}

func RunServer2(port string) {
	http.HandleFunc("/", helloWorldHandler)
	http.ListenAndServe(":"+port, nil)
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}


// Middleware handler to emit standard request logging information
// Place in the chain AFTER security principle decoding 
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("{ path: \"" + r.RequestURI + "\"}")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}


// Run the echoSleepHTTP server process
func RunServer(port string, server *http.Server) {
	// inject the service global context values
	echoSleepCtx := &echoSleepCtx{ServiceName: "echoSleepHTTP"}

	log.Println("{ service: \"" + echoSleepCtx.ServiceName + "\", start: true,  port: " + port + "}")

	echoSleepHandler :=
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			echoSleepCtx.echoSleepHTTP(w, r)
		},
		)

	router := httprouter.New()

	router.GET("/js", router.FileServe("/"))
	router.GET("/css", router.FileServe("/"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))

	router.GET("/{rest:.*}", echoSleepHandler)
	router.PUT("/{rest:.*}", echoSleepHandler)
	router.POST("/{rest:.*}", echoSleepHandler)

	// Create the handlers that will be wrapped by the middleware.
	pushHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Push"))
	})
	pullHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pull"))
	})

	router.GET("/push", pushChain)
	router.GET("/pull", pullChain)

	//http.Handle("/", router)
	server.Handler = router.Handler()

		router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			return true
			// r.ProtoMajor == 0
		})
		// add our default logger
		router.Use(loggingMiddleware)

	fmt.Println("fire up server ")
	err := server.ListenAndServe(":80",router)
	if err == nil {
		log.Fatal(err)
	}
}
