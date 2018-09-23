// blog assets manager - client
// bamgr [-d <domain>] [<folder>] <file>
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	exitOk = iota
	exitParam
	exitHTTP
)

var flags struct {
	server       string
	domain       string
	authUser     string
	authPass     string
	preserveName bool
}

func init() {
	flag.StringVar(&flags.server, "s", "", "server address")
	flag.StringVar(&flags.domain, "d", "", "availability domain")
	flag.StringVar(&flags.authUser, "u", "", "basic auth user")
	flag.StringVar(&flags.authPass, "p", "", "basic auth pass")
	flag.BoolVar(&flags.preserveName, "n", false, "preserve file name")
	flag.Parse()

	if flags.server == "" {
		flag.PrintDefaults()
		os.Exit(exitParam)
	}
	if !strings.HasPrefix(flags.server, "http") {
		flags.server = "http://" + flags.server
	}
}

func parseArgs() (folder, file string) {
	argv := flag.Args()
	argc := len(flag.Args())
	if argc < 1 || argc > 2 {
		flag.PrintDefaults()
		os.Exit(exitParam)
	}
	if argc == 2 {
		folder = argv[0]
		file = argv[1]
	} else {
		file = argv[0]
	}
	return
}

func main() {
	folder, file := parseArgs()
	params := make(map[string]string)
	params["domain"] = flags.domain
	params["folder"] = folder
	if flags.preserveName {
		fmt.Println("should preserve name")
		params["file_name"] = filepath.Base(file)
	}
	req, err := newfileUploadRequest(flags.server, params, "file", file)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(exitHTTP)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(exitHTTP)
	} else {
		var bodyContent []byte
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		resp.Body.Read(bodyContent)
		resp.Body.Close()
		fmt.Println(bodyContent)
	}
}

// ref https://gist.github.com/mattetti/5914158/f4d1393d83ebedc682a3c8e7bdc6b49670083b84
// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string,
	params map[string]string, paramName,
	path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	fmt.Println("file type", writer.FormDataContentType())
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if flags.authUser != "" {
		req.SetBasicAuth(flags.authUser, flags.authPass)
	}
	return req, err
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	// rewind file
	out.Seek(0, 0)

	return contentType, nil
}
