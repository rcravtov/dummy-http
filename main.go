package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
)

var hub *Hub

func main() {

	var port int
	flag.IntVar(&port, "port", 8080, "HTTP Server port number")
	flag.Parse()

	hub = NewHub()

	http.HandleFunc("GET /log", HandleLog)
	http.HandleFunc("/handle", HandleRequest)
	http.HandleFunc("GET /ws", HandleWS)

	fmt.Printf("Starting server on port: %d\n", port)
	fmt.Println("Expecting connections on /handel")
	fmt.Println("Check incoming requests on /log")
	http.ListenAndServe(":"+strconv.Itoa(port), nil)

}
