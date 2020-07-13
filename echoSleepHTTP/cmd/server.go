package cmd

import (
	"context"
	serverLogic "github.com/buildmicroservices/TestDriveServices/echoSleepHTTP/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//   server "echoSleepHTTP"

func LogSetup() *log.Logger {
	// setup logger
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")
	return logger
}

func InterruptSetup(server *http.Server, logger *log.Logger, done chan bool) {
	// make channels for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// background gofunc to kill server on interupt
	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		//atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		// call the done channel
		close(done)
	}()

}

func RunServer(port string) {
	// wait for done channel listener
	done := make(chan bool)

	logger := LogSetup()
	server := &http.Server{Addr: "localhost:"+port, Handler: nil}
	InterruptSetup(server,logger, done)
	serverLogic.RunServer(server)

	//err != nil && err != http.ErrServerClosed {
	//	logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	//}

	//	http.Handle("/metrics", promhttp.Handler())
	//	http.ListenAndServe(":9001", nil)

	// wait here...
	<-done
	logger.Println("Server stopped")

}
