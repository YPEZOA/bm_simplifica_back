package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ypezoa/bm-simplifica-back/internal/middleware"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	emailService "github.com/ypezoa/bm-simplifica-back/internal/services/email"
	"github.com/ypezoa/bm-simplifica-back/internal/types"
	"github.com/ypezoa/bm-simplifica-back/internal/validation"
)

func AuthRoutes(r *mux.Router) {
	r.HandleFunc("/sign-in", SignIn).Methods("POST")
	r.HandleFunc("/contact", ContactRequest).Methods("POST")
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	log.Printf("SignIn endpoint called - Method: %s, Path: %s", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	var creds types.AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Formato de solicitud inválido",
		})
		return
	}

	log.Printf("Received credentials - Email: '%s', Password: '%s'", creds.Email, creds.Password)

	validator := validation.LoginValidator{
		Email:    creds.Email,
		Password: creds.Password,
	}

	if err := validator.Validate(); err != nil {
		log.Printf("Validation error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	log.Printf("Validation passed")

	authStorage := NewAuthStorage()
	user, err := authStorage.SignIn(creds.Email, creds.Password)
	if err != nil {
		log.Printf("SignIn error: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Credenciales incorrectas",
		})
		return
	}

	log.Printf("SignIn successful for user: %s", user.Email)

	jwtMiddleware := middleware.NewJWTMiddleware(getJWTSecret())
	token, err := jwtMiddleware.GenerateToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al generar token",
		})
		return
	}

	userResponse := models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Phone:     user.Phone,
		Companies: user.Companies,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"user":  userResponse,
			"token": token,
		},
	})
}

func ContactRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var contact validation.ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Formato de solicitud inválido",
		})
		return
	}

	if err := contact.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	emailSvc := emailService.NewEmailService()
	err := emailSvc.SendNewUserNotification(contact.Name, contact.Email, contact.Phone, contact.Message)
	// For development: log the error but don't fail the request
	if err != nil {
		log.Printf("Error sending email notification: %v", err)
		log.Printf("Contact form submission: Name=%s, Email=%s, Phone=%s",
			contact.Name, contact.Email, contact.Phone)

		// Still return success to user, but note that email failed
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: true,
			Data:    "Su solicitud ha sido recibida. Hemos registrado tu información y te contactaremos pronto.",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data:    "Su solicitud ha sido enviada exitosamente. Nos contactaremos pronto.",
	})
}

func getJWTSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	return "your-secret-key-change-in-production"
}
