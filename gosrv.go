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

func main() {
	parseFlags()
	fmt.Println("listening on port", port)
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))
}
