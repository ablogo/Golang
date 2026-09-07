package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"src/models"
	"src/services"
)

func GetUser(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	user_email := queryParams.Get("email")
	if user_email == "" {
		http.Error(w, "User email is required", http.StatusBadRequest)
		return
	}

	user := services.GetUserByEmail(user_email)
	if user != nil {
		json.NewEncoder(w).Encode(user)
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	users := services.GetUsers()

	if users != nil {
		json.NewEncoder(w).Encode(users)
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	user_email := queryParams.Get("email")
	if user_email == "" {
		http.Error(w, "User email is required", http.StatusBadRequest)
		return
	}

	user := services.GetUserByEmail(user_email)
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result := services.DeleteUser(user.Id)
	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusNotModified)
		return
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	var model models.User
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	result := services.UpdateUser(model)
	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusNotModified)
		return
	}

}

func ChangePassword(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	user_email := queryParams.Get("email")
	if user_email == "" {
		http.Error(w, "User email is required", http.StatusBadRequest)
		return
	}

	password := queryParams.Get("password")
	if password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	user := services.GetUserByEmail(user_email)
	if user != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result := services.ChangePassword(user, password)
	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusNotModified)
		return
	}

}

func AddAddress(w http.ResponseWriter, r *http.Request) {

	var model models.Address
	defer r.Body.Close()

	queryParams := r.URL.Query()

	user_email := queryParams.Get("email")
	if user_email == "" {
		http.Error(w, "User email is required", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	user := services.GetUserByEmail(user_email)
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	address, result := services.AddAddress(user, model)
	if result {
		json.NewEncoder(w).Encode(address)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func GetAddresses(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	user_email := queryParams.Get("email")
	if user_email == "" {
		http.Error(w, "User email is required", http.StatusBadRequest)
		return
	}

	user := services.GetUserByEmail(user_email)
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	addresses := services.GetAddressByUser(user.Id)
	if len(*addresses) > 0 {
		json.NewEncoder(w).Encode(addresses)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func DeleteAddress(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	user_email := queryParams.Get("email")
	if user_email == "" {
		http.Error(w, "User email is required", http.StatusBadRequest)
		return
	}

	address_id, err := strconv.Atoi(queryParams.Get("address_id"))
	if err != nil {
		http.Error(w, "Address id is required", http.StatusBadRequest)
		return
	}

	user := services.GetUserByEmail(user_email)
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	address := services.GetAddress(address_id)
	if address.UserId == user.Id {

		result := services.DeleteAddress(address_id)
		if result {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotModified)
		}
	} else {
		http.Error(w, "Don't belong to the user", http.StatusNotAcceptable)

	}
}

func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	var model models.Address
	defer r.Body.Close()

	queryParams := r.URL.Query()

	user_email := queryParams.Get("email")
	if user_email == "" {
		http.Error(w, "User email is required", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	user := services.GetUserByEmail(user_email)
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if model.UserId == user.Id {

		result := services.UpdateAddress(model)
		if result {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotModified)
		}
	} else {
		http.Error(w, "Don't belong to the user", http.StatusNotAcceptable)
	}
}
