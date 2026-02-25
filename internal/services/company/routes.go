package company

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ypezoa/bm-simplifica-back/internal/middleware"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	file "github.com/ypezoa/bm-simplifica-back/internal/services/file"
	user "github.com/ypezoa/bm-simplifica-back/internal/services/user"
	"github.com/ypezoa/bm-simplifica-back/internal/types"
	"gorm.io/gorm"
)

func CompanyRoutes(r *mux.Router) {
	subrouter := r.PathPrefix("/companies").Subrouter()
	subrouter.Use(middleware.NewJWTMiddleware(getJWTSecret()).AuthMiddleware)

	// Admin-only routes
	adminSubrouter := subrouter.PathPrefix("/admin").Subrouter()
	adminSubrouter.Use(middleware.NewJWTMiddleware(getJWTSecret()).AdminMiddleware)
	adminSubrouter.HandleFunc("", GetAllCompanies).Methods("GET")
	adminSubrouter.HandleFunc("", CreateCompanyByAdmin).Methods("POST")
	adminSubrouter.HandleFunc("/{id}", DeleteCompanyByAdmin).Methods("DELETE")

	// Client routes (own companies only)
	subrouter.HandleFunc("", GetUserCompanies).Methods("GET")
	subrouter.HandleFunc("/{id}/files", GetCompanyFiles).Methods("GET")
}

func GetCompanies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	companies, err := NewCompanyStorage().GetCompanies()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al obtener empresas: " + err.Error(),
		})
		return
	}

	var companiesData []models.CompaniesResponse

	for _, company := range companies {
		companiesData = append(companiesData, models.CompaniesResponse{
			Name:  company.Name,
			Rut:   company.Rut,
			Files: company.Files,
		})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data:    companiesData,
	})
}

func CreateCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var company models.Company

	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al leer los datos de la empresa: " + err.Error(),
		})
		return
	}

	// Se asegura de inicializar el slice de archivos como [] vacío si es nil
	if company.Files == nil {
		company.Files = []models.File{}
	}

	createdCompany, err := NewCompanyStorage().CreateCompany(company)

	createdCompanyResponse := models.Company{
		ID:     createdCompany.ID,
		Name:   createdCompany.Name,
		Rut:    createdCompany.Rut,
		Files:  createdCompany.Files,
		UserID: createdCompany.UserID,
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al crear la empresa" + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data:    createdCompanyResponse,
	})
}

func GetAllCompanies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	companies, err := NewCompanyStorage().GetCompanies()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al obtener empresas: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data:    companies,
	})
}

func CreateCompanyByAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Estructura intermedia para parsing
	var request struct {
		Name   string        `json:"name"`
		Rut    string        `json:"rut"`
		UserID string        `json:"user_id"`
		Files  []models.File `json:"files,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Formato de solicitud inválido",
		})
		return
	}

	// Debug logging
	log.Printf("Request received: Name=%s, RUT=%s, UserID=%s", request.Name, request.Rut, request.UserID)

	// Parsear UUID
	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "ID de usuario inválido",
		})
		return
	}

	// Validar que UserID no esté vacío
	if userID == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "El ID de usuario es requerido",
		})
		return
	}

	// Crear modelo Company
	company := models.Company{
		Name:   request.Name,
		Rut:    request.Rut,
		UserID: userID,
		Files:  request.Files,
	}

	// Validar datos de la compañía
	if company.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "El nombre de la compañía es requerido",
		})
		return
	}

	if company.Rut == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "El RUT de la compañía es requerido",
		})
		return
	}

	if company.Files == nil {
		company.Files = []models.File{}
	}

	createdCompany, err := NewCompanyStorage().CreateCompany(company)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(types.APIResponse{
				Success: false,
				Error:   "El usuario especificado no existe",
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.APIResponse{
				Success: false,
				Error:   "Error al crear compañía: " + err.Error(),
			})
		}
		return
	}

	// Preparar respuesta limpia (sin información sensible del usuario)
	response := map[string]interface{}{
		"id":      createdCompany.ID,
		"name":    createdCompany.Name,
		"rut":     createdCompany.Rut,
		"user_id": createdCompany.UserID,
		"files":   createdCompany.Files,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data:    response,
	})
}

func GetUserCompanies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	userStorage := user.NewUserStorage()
	user, err := userStorage.GetUserByID(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Usuario no encontrado",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data:    user.Companies,
	})
}

func GetCompanyFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	companyID, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "ID de compañía inválido",
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

	fileStorage := file.NewFileStorage()
	files, err := fileStorage.GetFiles(userID, companyID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al obtener archivos: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data:    files,
	})
}

func DeleteCompanyByAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	companyID, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "ID de compañía inválido",
		})
		return
	}

	// Obtener compañía para verificar que existe
	companyStorage := NewCompanyStorage()
	company, err := companyStorage.GetCompanyByID(companyID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Compañía no encontrada",
		})
		return
	}

	// Eliminar compañía (esto eliminará sus archivos en cascada)
	err = companyStorage.DeleteCompany(companyID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al eliminar compañía: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message": "Compañía eliminada exitosamente",
			"company_deleted": map[string]interface{}{
				"id":   company.ID,
				"name": company.Name,
				"rut":  company.Rut,
			},
		},
	})
}

func getJWTSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	return "your-secret-key-change-in-production"
}
