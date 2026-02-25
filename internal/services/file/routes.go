package file

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ypezoa/bm-simplifica-back/internal/middleware"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"github.com/ypezoa/bm-simplifica-back/internal/types"
)

func FileRoutes(r *mux.Router) {
	subrouter := r.PathPrefix("/files").Subrouter()
	subrouter.Use(middleware.NewJWTMiddleware(getJWTSecret()).AuthMiddleware)

	// Admin-only routes
	adminSubrouter := subrouter.PathPrefix("/admin").Subrouter()
	adminSubrouter.Use(middleware.NewJWTMiddleware(getJWTSecret()).AdminMiddleware)
	adminSubrouter.HandleFunc("/upload/{company-id}", UploadFileByAdmin).Methods("POST")
	adminSubrouter.HandleFunc("/delete", DeleteFilesByAdmin).Methods("POST")

	// Client routes (download files from own companies)
	subrouter.HandleFunc("/{id}/download", DownloadFile).Methods("GET")
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(50 << 50)

	vars := mux.Vars(r)
	companyID, err := uuid.Parse(vars["company-id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Compañía inválida: " + err.Error(),
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

	// Validate file
	if err := validateFile(fileHeader, handler.Filename, handler.Header.Get("Content-Type"), 50<<20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Validación de archivo falló: " + err.Error(),
		})
		return
	}

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
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data:    storedFile,
	})
}

func UploadFileByAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseMultipartForm(50 << 50)

	vars := mux.Vars(r)
	companyID, err := uuid.Parse(vars["company-id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "ID de compañía inválido: " + err.Error(),
		})
		return
	}

	fileHeader, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al leer archivo: " + err.Error(),
		})
		return
	}
	defer fileHeader.Close()

	if err := validateFile(fileHeader, handler.Filename, handler.Header.Get("Content-Type"), 50<<20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Validación de archivo falló: " + err.Error(),
		})
		return
	}

	newFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
	path := "./uploads/" + newFileName

	dst, err := os.Create(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al crear archivo: " + err.Error(),
		})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, fileHeader); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al copiar archivo: " + err.Error(),
		})
		return
	}

	newFile := models.File{
		Name:      handler.Filename,
		Path:      path,
		Type:      handler.Header.Get("Content-Type"),
		CompanyID: companyID,
	}

	storedFile, err := NewFileStorage().UploadFile(newFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al guardar archivo: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data:    storedFile,
	})
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "ID de archivo inválido",
		})
		return
	}

	claims, ok := r.Context().Value("userClaims").(*middleware.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Claims de usuario no encontrados",
		})
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "ID de usuario inválido en token",
		})
		return
	}

	fileStorage := NewFileStorage()
	files, err := fileStorage.GetFiles(userID, uuid.Nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al obtener archivo: " + err.Error(),
		})
		return
	}

	var targetFile *models.File
	for _, f := range files {
		if f.ID == fileID {
			targetFile = &f
			break
		}
	}

	if targetFile == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Archivo no encontrado o acceso denegado",
		})
		return
	}

	http.ServeFile(w, r, targetFile.Path)
}

func DeleteFilesByAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		FileIDs []string `json:"file_ids" binding:"required,min=1"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Formato de solicitud inválido",
		})
		return
	}

	if len(req.FileIDs) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Debe proporcionar al menos un ID de archivo",
		})
		return
	}

	// Convertir string UUIDs a UUID objects
	fileUUIDs := make([]uuid.UUID, len(req.FileIDs))
	for i, idStr := range req.FileIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(types.APIResponse{
				Success: false,
				Error:   "UUID de archivo inválido: " + idStr,
			})
			return
		}
		fileUUIDs[i] = id
	}

	// Eliminar archivos
	fileStorage := NewFileStorage()
	deletedFiles, err := fileStorage.DeleteMultipleFiles(fileUUIDs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al eliminar archivos: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message":       "Archivos eliminados exitosamente",
			"deleted_count": len(deletedFiles),
			"deleted_files": deletedFiles,
		},
	})
}

func getJWTSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	return "your-secret-key-change-in-production"
}

func validateFile(fileHeader io.Reader, filename, contentType string, maxSize int64) error {
	if maxSize > 0 {
		fileHeader.(*os.File).Seek(0, 0)
		buf := make([]byte, maxSize+1)
		n, _ := fileHeader.Read(buf)
		if int64(n) > maxSize {
			return fmt.Errorf("el tamaño del archivo excede el máximo permitido")
		}
		fileHeader.(*os.File).Seek(0, 0)
	}

	allowedTypes := []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"image/jpeg",
		"image/png",
		"text/plain",
	}

	isAllowed := false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return fmt.Errorf("tipo de archivo no permitido: %s", contentType)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return fmt.Errorf("el archivo debe tener una extensión válida")
	}

	dangerousExtensions := []string{".exe", ".bat", ".cmd", ".scr", ".pif", ".com", ".js", ".vbs", ".jar", ".app", ".deb", ".pkg", ".dmg"}
	for _, dangerousExt := range dangerousExtensions {
		if ext == dangerousExt {
			return fmt.Errorf("extensión de archivo no permitida por razones de seguridad")
		}
	}

	return nil
}
