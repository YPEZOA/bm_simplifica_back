package user

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ypezoa/bm-simplifica-back/internal/middleware"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	auth "github.com/ypezoa/bm-simplifica-back/internal/services/auth"
	emailService "github.com/ypezoa/bm-simplifica-back/internal/services/email"
	"github.com/ypezoa/bm-simplifica-back/internal/types"
	"github.com/ypezoa/bm-simplifica-back/internal/validation"
	"golang.org/x/crypto/bcrypt"
)

func UserRoutes(r *mux.Router) {
	subrouter := r.PathPrefix("/users").Subrouter()
	subrouter.Use(middleware.NewJWTMiddleware(getJWTSecret()).AuthMiddleware)

	// Admin-only routes
	adminSubrouter := subrouter.PathPrefix("/admin").Subrouter()
	adminSubrouter.Use(middleware.NewJWTMiddleware(getJWTSecret()).AdminMiddleware)
	adminSubrouter.HandleFunc("", GetAllUsers).Methods("GET")
	adminSubrouter.HandleFunc("/{id}", GetUserByID).Methods("GET")
	adminSubrouter.HandleFunc("", CreateUserByAdmin).Methods("POST")
	adminSubrouter.HandleFunc("/{id}", DeleteUser).Methods("DELETE")
	adminSubrouter.HandleFunc("/{id}/change-password", ChangeUserPassword).Methods("POST")
	adminSubrouter.HandleFunc("/{id}/send-welcome-email", ResendWelcomeEmail).Methods("POST")

	// Client routes (own data only)
	subrouter.HandleFunc("/me", GetCurrentUser).Methods("GET")
}

func getJWTSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	return "your-secret-key-change-in-production"
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "ID de usuario no encontrado",
		})
		return
	}

	// Verificar que no sea el usuario actual admin
	claims, ok := r.Context().Value("userClaims").(*middleware.Claims)
	if ok && claims.UserID == userID.String() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "No puedes eliminar tu propio usuario administrador",
		})
		return
	}

	// Obtener usuario para verificar si existe
	userStorage := NewUserStorage()
	user, err := userStorage.GetUserByID(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Usuario no encontrado",
		})
		return
	}

	// Soft delete del usuario (marcar como eliminado, no eliminar físicamente)
	err = userStorage.DeleteUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al eliminar usuario: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message": "Usuario eliminado exitosamente (soft delete)",
			"user_deleted": map[string]interface{}{
				"id":         user.ID,
				"name":       user.Name,
				"email":      user.Email,
				"deleted_at": "Marcado como eliminado en base de datos",
			},
		},
	})
}

func ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "ID de usuario no encontrado",
		})
		return
	}

	// Parsear request body
	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Formato de solicitud inválido",
		})
		return
	}

	// Validar nueva contraseña
	validator := validation.UserValidator{
		Password: req.NewPassword,
	}

	if err := validator.ValidatePassword(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Contraseña inválida: " + err.Error(),
		})
		return
	}

	// Verificar que usuario existe
	userStorage := NewUserStorage()
	user, err := userStorage.GetUserByID(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Usuario no encontrado",
		})
		return
	}

	// Generar hash de la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al procesar la contraseña",
		})
		return
	}

	// Actualizar contraseña
	err = userStorage.UpdateUserPassword(userID, string(hashedPassword))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al actualizar la contraseña: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message": "Contraseña actualizada exitosamente",
			"user_updated": map[string]interface{}{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
		},
	})
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Users list from DB
	users, err := NewUserStorage().GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al obtener usuarios: " + err.Error(),
		})
		return

	}

	// Prepare response data
	var usersData []models.UserResponse
	for _, user := range users {
		usersData = append(usersData, models.UserResponse{
			ID:        user.ID,
			Role:      user.Role,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			Companies: user.Companies,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data:    usersData,
	})
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	userID, _ := uuid.Parse(vars["id"])

	user, err := NewUserStorage().GetUserByID(userID)
	userResponse := models.UserResponse{
		ID:        user.ID,
		Role:      user.Role,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Companies: user.Companies,
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Usuario no encontrado",
		})
		return

	}

	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data:    userResponse,
	})
}

func CreateUserByAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Formato de solicitud inválido",
		})
		return
	}

	validator := validation.UserValidator{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Phone:    user.Phone,
	}

	if err := validator.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	authStorage := auth.NewAuthStorage()
	createdUser, err := authStorage.CreateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "El usuario ya existe o la creación falló",
		})
		return
	}

	// Enviar email de bienvenida automáticamente
	emailSvc := emailService.NewEmailService()
	err = emailSvc.SendWelcomeEmail(createdUser.Email, createdUser.Name, user.Password)
	if err != nil {
		// Log error but don't fail request
		log.Printf("Error al enviar email de bienvenida a %s: %v", createdUser.Email, err)
	}

	userResponse := models.UserResponse{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		Role:      createdUser.Role,
		Phone:     createdUser.Phone,
		Companies: createdUser.Companies,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"user":       userResponse,
			"email_sent": err == nil,
			"message": "Usuario creado. Email de bienvenida " + func() string {
				if err == nil {
					return "enviado exitosamente"
				}
				return "falló (revisar logs)"
			}(),
		},
	})
}

func ResendWelcomeEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "ID de usuario inválido",
		})
		return
	}

	userStorage := NewUserStorage()
	user, err := userStorage.GetUserByID(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Usuario no encontrado",
		})
		return
	}

	// Leer el password desde el request (para reenviar con misma contraseña)
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Formato de solicitud inválido",
		})
		return
	}

	// Enviar email de bienvenida
	emailSvc := emailService.NewEmailService()
	err = emailSvc.SendWelcomeEmail(user.Email, user.Name, req.Password)
	if err != nil {
		log.Printf("Error al reenviar email de bienvenida a %s: %v", user.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al enviar email de bienvenida",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data:    "Email de bienvenida reenviado exitosamente",
	})
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
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

	userStorage := NewUserStorage()
	user, err := userStorage.GetUserByID(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Usuario no encontrado",
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
		Data:    userResponse,
	})
}
