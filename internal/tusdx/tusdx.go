package tusdx

import (
	"fmt"

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
			producer.Publish(string(event.Upload.ID))
			fmt.Printf("Upload %s finished\n", event.Upload.ID)
		}
	}()
	return h
}
