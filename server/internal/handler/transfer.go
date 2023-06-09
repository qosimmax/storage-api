package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/qosimmax/storage-api/user"
)

const limit = 6

func TransferFile(
	db user.ServerFinder,
	ft user.FileTransfer,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		f, fh, err := r.FormFile("file")
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, true)
			return
		}
		defer f.Close()

		size := fh.Size
		filename := fh.Filename
		partSize := size / limit

		fileID, _ := generateUUID()

		servers, err := db.GetServers(ctx, limit)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, true)
			return
		}

		wg := sync.WaitGroup{}
		for i, server := range servers {
			offset := int64(i) * partSize
			size := partSize
			if i == limit-1 {
				size = partSize + size%limit
			}

			wg.Add(1)
			go func(offset, size int64, server user.ServerData) {
				_, _ = ft.TransferFile(ctx, user.FileInfo{
					ID:     fileID,
					File:   r.MultipartForm.File["file"][0],
					Offset: offset,
					Size:   size,
				}, server)
				wg.Done()

			}(offset, size, server)

		}

		wg.Wait()

		log.Println(fileID, filename, size)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	}
}

func generateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("error generating uuid: %w", err)
	}
	return id.String(), nil
}
