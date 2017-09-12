package main

import (
	"flag"
	"net/http"
	"fmt"
	"log"
	"github.com/suncat2000/bid-service/requester"
)

func main() {
	listenAddr := flag.String("addr", ":8080", "http listen address")
	flag.Parse()

	requestHandler := requester.NewRequester()

	http.HandleFunc("/winner", func(writer http.ResponseWriter, request *http.Request) {
		requestHandler.Handle(writer, request)
	})

	fmt.Printf("Listen on %s\n", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}