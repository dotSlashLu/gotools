package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const defaultPort = 10080

var port int

func parseFlags() {
	flag.IntVar(&port, "p", defaultPort, "binding port")
	flag.Parse()
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fileServer := http.FileServer(http.Dir("."))
	enableCors(&w)
	fileServer.ServeHTTP(w, r)
}

func main() {
	parseFlags()
	fmt.Println("listening on port", port)
	// http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))
}
