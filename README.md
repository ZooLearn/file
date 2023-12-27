Open Source


# prepare:
- install [golang](https://go.dev/)
- install [pre-commit](https://pre-commit.com/)
- install [rabbitmq](https://www.rabbitmq.com/tutorials/tutorial-one-go.html)
- install [ffmpeg](https://ffmpeg.org/download.html)
  
# upload file example
```
package main

import (
	"fmt"
	"os"

	"github.com/eventials/go-tus"
)

func main() {
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
```