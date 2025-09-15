package user

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"github.com/ypezoa/bm-simplifica-back/internal/types"
)

func UserRoutes(r *mux.Router) {
	r.HandleFunc("/users", GetAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", GetUserByID).Methods("GET")
	r.HandleFunc("/users", CreateUser).Methods("POST")
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := NewUserStorage().GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al obtener usuarios: " + err.Error(),
		})
		return

	}

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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al leer los datos del usuario: " + err.Error(),
		})
		return
	}

	createdUser, err := NewUserStorage().CreateUser(user)
	createdUserResponse := models.UserResponse{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		Phone:     createdUser.Phone,
		Role:      createdUser.Role,
		Companies: createdUser.Companies,
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Error al crear usuario: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: true,
		Data:    createdUserResponse,
	})
}
