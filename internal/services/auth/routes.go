package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ypezoa/bm-simplifica-back/internal/models"
	"github.com/ypezoa/bm-simplifica-back/internal/types"
)

func AuthRoutes(r *mux.Router) {
	r.HandleFunc("/sign-in", SignIn).Methods("POST")
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var creds types.AuthCredentials
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		json.NewEncoder(w).Encode(types.APIResponse{
			Success: false,
			Error:   "Credenciales incorrectas" + err.Error(),
		})
		return
	}

	user, err := NewAuthStorage().SignIn(creds.Email, creds.Password)

	fmt.Println(user)
	fmt.Println(err)
}
