package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/cors"
)

var hub *Hub

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "HTTP Server port number")
	flag.Parse()

	hub = NewHub()

	mux := http.NewServeMux()

	mux.HandleFunc("/", HandleLog)
	mux.HandleFunc("/handle", HandleRequest)
	mux.HandleFunc("/ws", HandleWS)

	fmt.Printf("Starting server on port: %d\n", port)
	fmt.Println("Expecting connections on /handle")
	fmt.Println("Check incoming requests on /")

	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":"+strconv.Itoa(port), handler)
}
