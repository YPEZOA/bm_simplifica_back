package file

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"github.com/ypezoa/bm-simplifica-back/internal/types"
)

func FileRoutes(r *mux.Router) {
	r.HandleFunc("/file/{company-id}", UploadFile).Methods("POST")
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(50 << 50)

	vars := mux.Vars(r)
	companyID, err := uuid.Parse(vars["company-id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Invalid company: " + err.Error(),
		})
		return
	}

	// Archivo subido desde el form
	fileHeader, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al leer el archivo: " + err.Error(),
		})
		return
	}

	defer fileHeader.Close()

	// Generamos un nombre único y guardamos en disco
	newFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
	path := "./uploads/" + newFileName

	// Se crea archivo en disco
	dst, err := os.Create(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al subir el archivo: " + err.Error(),
		})
	}

	defer dst.Close()

	// Copiamos el contenido del archivo subido (file-Header) dentro del archivo en disco dst.
	if _, err := io.Copy(dst, fileHeader); err != nil {
		http.Error(w, "Error al copiar archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	newFile := models.File{Name: handler.Filename, Path: path, Type: handler.Header.Get("Content-Type"), CompanyID: companyID}
	storedFile, err := NewFileStorage().UploadFile(newFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al guardar el archivo: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(storedFile)
}
