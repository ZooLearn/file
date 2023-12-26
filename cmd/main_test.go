package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/eventials/go-tus"
)

func TestDemo(t *testing.T) {
	f, err := os.Open("./main.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// create the tus client.
	client, err := tus.NewClient("http://localhost:8080/files", nil)
	if err != nil {
		panic(err)
	}
	// create an upload from a file.
	upload, err := tus.NewUploadFromFile(f)
	if err != nil {
		panic(err)
	}

	// create the uploader.
	uploader, err := client.CreateUpload(upload)
	if err != nil {
		panic(err)
	}

	// start the uploading process.
	if err := uploader.Upload(); err != nil {
		panic(err)
	}

	// link download
	fmt.Println(uploader.Url())
}
