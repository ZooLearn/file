package tusdx

import (
	"fmt"

	"github.com/ZooLearn/file/internal/log"
	"github.com/ZooLearn/file/internal/rabbitx"
	"github.com/tus/tusd/pkg/filestore"
	tusdHandler "github.com/tus/tusd/pkg/handler"
)

func TusdHandler(producer *rabbitx.Producer) *tusdHandler.UnroutedHandler {
	store := filestore.FileStore{
		Path: "./uploads",
	}

	composer := tusdHandler.NewStoreComposer()
	store.UseIn(composer)

	h, err := tusdHandler.NewUnroutedHandler(tusdHandler.Config{
		BasePath:              "/files/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})

	if err != nil {
		panic(fmt.Errorf("Unable to create handler: %s", err))
	}

	go func() {
		for {
			event := <-h.CompleteUploads
			err := producer.Publish(string(event.Upload.ID))
			if err != nil {
				log.Errorf("publish to rabbit: %s", err)
			}
			fmt.Printf("Upload %s finished\n", event.Upload.ID)
		}
	}()
	return h
}
