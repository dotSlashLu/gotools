package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

const (
	maxUploadSize = 2 * 1024 * 1024 // 2 mb
	uploadPath    = "/var/www/img"
	fileField     = "file"
)

var flags struct {
	Port int64
}

func init() {
	flag.Int64Var(&flags.Port, "p", 8080, "Listening port")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", uploadFileHandler())

	port := strconv.Itoa(int(flags.Port))
	log.Print("Server started on localhost:" + port +
		", use / for uploading files")
	log.Fatal(http.ListenAndServe("localhost:"+port, nil))
}
