package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func getFolder(r *http.Request) (string, error) {
	// other params
	domain := r.PostFormValue("domain")
	folder := r.PostFormValue("folder")
	path := filepath.Join(uploadPath, domain, folder)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.Mkdir(path, 0755); err != nil {
			log.Println("err creating folder: ", err)
			return "", err
		}
	}
	return path, nil
}

// ref: https://zupzup.org/go-http-file-upload-download/
func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate file size
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			log.Println(err)
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile(fileField)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		// check file type, detectcontenttype only needs the first 512 bytes
		filetype := http.DetectContentType(fileBytes)
		switch filetype {
		case "image/jpeg", "image/jpg":
		case "image/gif", "image/png":
		case "application/pdf":
			break
		default:
			renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}
		fileName := ""
		formFileName := r.PostFormValue("file_name")
		log.Println("assigned file_name", formFileName)
		if formFileName != "" {
			fileName = formFileName
		} else {
			fileEndings, err := mime.ExtensionsByType(filetype)
			if err != nil {
				renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
				return
			}
			fileName = randToken(12) + fileEndings[0]
		}
		folder, err := getFolder(r)
		if err != nil {
			renderError(w, "CANT_FIND_FOLDER", http.StatusInternalServerError)
			return
		}
		newPath := filepath.Join(folder, fileName)
		log.Printf("FileType: %s, File: %s\n", filetype, newPath)

		// write file
		newFile, err := os.Create(newPath)
		if err != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("SUCCESS"))
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	log.Println("error, ", message)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
