package server

import (
	"fmt"
	"log"
	"net/http"

	"time"
)

func RunServer2(port string) {
	http.HandleFunc("/", helloWorldHandler)
	http.ListenAndServe(":"+port, nil)
}

// Run the echoSleepHTTP server process
func RunServer(server *http.Server) {

	var router = initializeRouter()

	server.Handler = router
	
	fmt.Println("fire up server ")
	server.ReadTimeout, _ = time.ParseDuration("10s")
	server.WriteTimeout, _ = time.ParseDuration("2m")
	// server.BaseContext
	// server.ConnContext

	log.Println("about to listen for echo request VVVVVVVVVV")
	//server.Addr = "localhost:8090"
	log.Print(server.Addr)
	err := server.ListenAndServe()
	//err := http.ListenAndServe(":8090",router)
	if err == nil {
		log.Fatal(err)
	}
}
