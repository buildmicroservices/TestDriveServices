package cmd

import (
	"net/http"
	serverLogic "github.com/buildmicroservices/TestDriveServices/echoSleepHTTP/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"context"
	"time"
)

//   server "echoSleepHTTP"

func RunServer(port string) {

	// setup logger
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")

	// make channels for interrupt
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	server := &http.Server{Addr: ":"+port, Handler: nil}

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
		close(done)
	}()

	fmt.Println("Port ", port)
	serverLogic.RunServer(port,server)


	//err != nil && err != http.ErrServerClosed {
	//	logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	//}

	<-done
	logger.Println("Server stopped")

	//	http.Handle("/metrics", promhttp.Handler())
	//	http.ListenAndServe(":9001", nil)

}
