package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/qosimmax/storage-api/server/internal/handler/mock"
)

func TestTransferFile(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "dummy.txt")
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, _ = part.Write([]byte(mock.DummyData))
	_ = writer.Close()

	req, err := http.NewRequest("POST", "/file", body)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TransferFile(mock.DBMock{}, mock.FileTransferMock{}))
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("json parse eror: %v", err)
	}

	// test ReceiveFile handler
	receiveFileFile(t, response["file_id"])
}

func receiveFileFile(t *testing.T, fileID string) {
	req, err := http.NewRequest("GET", "/file", nil)
	if err != nil {
		t.Fatal(err)
	}

	params := url.Values{}
	params.Set("file_id", fileID)
	req.URL.RawQuery = params.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ReceiveFile(mock.DBMock{}, mock.FileReceiverMock{}))
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	hashGot := sha256.New()
	_, _ = io.Copy(hashGot, rr.Body)

	hashExp := sha256.New()
	hashExp.Write([]byte(mock.DummyData))

	if fmt.Sprintf("%x", hashGot.Sum(nil)) != fmt.Sprintf("%x", hashExp.Sum(nil)) {
		t.Errorf("file cheksum failed: got %x want %x",
			hashGot.Sum(nil), hashExp.Sum(nil))
	}

}
