package company

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"github.com/ypezoa/bm-simplifica-back/internal/types"
)

func CompanyRoutes(r *mux.Router) {
	r.HandleFunc("/company", CreateCompany).Methods("POST")
	r.HandleFunc("/companies", GetCompanies).Methods("GET")
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
