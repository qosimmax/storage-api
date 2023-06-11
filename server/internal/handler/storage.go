package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/qosimmax/storage-api/user"
)

const limit = 6

func TransferFile(
	db user.ServerFileHandler,
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

		totalSize := fh.Size
		filename := fh.Filename
		partSize := fh.Size / limit
		remainSize := fh.Size % limit

		fileID, _ := generateUUID()

		servers, err := db.FindAvailableServers(ctx, limit)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, true)
			return
		}

		wg := sync.WaitGroup{}
		offset := int64(0)
		for i, server := range servers {
			size := partSize
			if remainSize > 0 {
				size += 1
				remainSize -= 1
			}

			wg.Add(1)
			go func(offset, size int64, order int, server user.ServerData) {
				_, _ = ft.SendFile(ctx, user.SrcFileInfo{
					ID:     fileID,
					File:   r.MultipartForm.File["file"][0],
					Offset: offset,
					Size:   size,
				}, server)

				_ = db.AddPartitionFileInfo(ctx, user.PartitionFileInfo{
					FileID:   fileID,
					ServerID: server.ID,
					PartSize: size,
					Order:    order,
				})
				wg.Done()

			}(offset, size, i, server)

			offset += size

		}

		wg.Wait()

		log.Println(fileID, filename, totalSize)
		err = db.AddFileInfo(ctx, user.FileInfo{
			ID:   fileID,
			Name: filename,
			Size: totalSize,
		})
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, true)
			return
		}

		response, _ := json.Marshal(map[string]string{"file_id": fileID})
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)

	}
}

func ReceiveFile(
	db user.FileInfoHandler,
	fr user.FileReceiver,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		fileID := r.URL.Query().Get("file_id")
		fileInfo, err := db.GetFileInfo(ctx, fileID)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, true)
			return
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name))
		w.Header().Set("Content-Type", "application/octet-stream")

		servers, err := db.FindFileLocationServers(ctx, fileID)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, true)
			return
		}

		for _, server := range servers {
			n, err := fr.ReceiveFile(ctx, fileID, server, w)
			log.Println(n, err)
		}

	}
}

func generateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("error generating uuid: %w", err)
	}
	return id.String(), nil
}