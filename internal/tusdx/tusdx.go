package tusdx

import (
	"encoding/json"
	"fmt"

	"github.com/ZooLearn/file/internal/log"
	"github.com/ZooLearn/file/internal/rabbitx"
	"github.com/tus/tusd/pkg/filestore"
	tusdHandler "github.com/tus/tusd/pkg/handler"
)

type Data struct {
	ID string
}

func TusdMediaHandler(producer *rabbitx.Producer) *tusdHandler.UnroutedHandler {
	store := filestore.FileStore{
		Path: "./uploads",
	}

	composer := tusdHandler.NewStoreComposer()
	store.UseIn(composer)

	h, err := tusdHandler.NewUnroutedHandler(tusdHandler.Config{
		BasePath:              "/media/files/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})

	if err != nil {
		panic(fmt.Errorf("unable to create handler: %s", err))
	}

	go func() {
		for {
			event := <-h.CompleteUploads
			input, err := json.Marshal(Data{
				ID: event.Upload.ID,
			})
			if err != nil {
				log.Errorf("publish to rabbit: %s", err)
				continue
			}
			fmt.Println(":=====>", string(input))
			err = producer.Publish(input)
			if err != nil {
				log.Errorf("publish to rabbit: %s", err)
			}
			fmt.Printf("Upload %s finished\n", event.Upload.ID)
		}
	}()
	return h
}
