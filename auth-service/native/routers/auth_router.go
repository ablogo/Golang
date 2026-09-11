package routers

import (
	"encoding/json"
	"net/http"

	"src/models"
	"src/services"
)

type Router struct {
	UserSvc *services.UserService
}

func (a *Router) SignUp(w http.ResponseWriter, r *http.Request) {
	var model models.SignUp
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&model); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid input"))
		return
	}

	result := a.UserSvc.CreateUser(model)

	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "Invalid token", http.StatusInternalServerError)
		return
	}
}

func (a *Router) SignIn(w http.ResponseWriter, r *http.Request) {
	var model models.SignIn

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	model.UserName = r.FormValue("username")
	model.Password = r.FormValue("password")

	token := a.UserSvc.Login(model.UserName, model.Password)

	if token != "" {
		w.WriteHeader(http.StatusOK)
		token := models.Token{AccessToken: token, TokenType: "bearer"}
		json.NewEncoder(w).Encode(token)
		return
	} else {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

}

func VerifyToken(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	user_id := queryParams.Get("user_id")
	if user_id == "" {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}
