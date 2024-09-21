package documents

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/uptrace/bun"
	"github.com/yesyash/yaskin-backend/internal/logger"
)

type fileStruct struct {
	Url string `json:"url"`
}

type uploadResponse struct {
	Success int        `json:"success"`
	File    fileStruct `json:"file"`
}

type document struct {
	bun *bun.DB
	ctx context.Context
}

// save the binary we get in request body to the public directory in the root folder
func (d *document) saveFileToDisc(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("image")
	if err != nil {
		logger.Error("error retrieving the file: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get the absolute path to the project root
	projectRoot, err := os.Getwd()
	if err != nil {
		logger.Error("error getting current working directory:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create the public folder path
	publicDir := filepath.Join(projectRoot, "public")

	if err := os.MkdirAll(publicDir, os.ModePerm); err != nil {
		logger.Error("error creating public directory: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create a new file in the public directory
	dst, err := os.Create(filepath.Join(publicDir, handler.Filename))
	if err != nil {
		logger.Error("error creating destination file: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		logger.Error("error copying file: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uploadedFileUrl := "http://localhost:4000/public/" + handler.Filename

	// Return the file URL
	fileUrl := fileStruct{Url: uploadedFileUrl}
	response := uploadResponse{Success: 1, File: fileUrl}

	jsonRes, err := json.Marshal(response)
	if err != nil {
		logger.Error("error marshalling response: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

func DocumentGroup(mux *http.ServeMux, ctx context.Context, db *bun.DB) {
	documentService := &document{db, ctx}
	mux.HandleFunc("POST /documents/upload", documentService.saveFileToDisc)
}
