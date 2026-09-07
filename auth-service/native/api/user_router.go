package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"src/models"
	"src/services"
)

func GetUser(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	user_id, ok := strconv.Atoi(queryParams.Get("user_id"))
	if ok != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user := services.GetUser(user_id, []string{"Address"})
	if user != nil {
		jsonBytes, err := json.Marshal(user)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
		return
	} else {
		http.Error(w, "", http.StatusNotFound)
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	user_id, ok := strconv.Atoi(queryParams.Get("user_id"))
	if ok != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	result := services.DeleteUser(user_id)
	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	var model models.User
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&model); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result := services.UpdateUser(model)
	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	user_id, ok := strconv.Atoi(queryParams.Get("user_id"))
	if ok != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	password := queryParams.Get("password")
	if password != "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
	}

	user := services.GetUser(user_id, nil)
	if user == nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	result := services.ChangePassword(user, password)
	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

}

func AddPicture(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	user_id, ok := strconv.Atoi(queryParams.Get("user_id"))
	if ok != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user := services.GetUser(user_id, nil)
	if user == nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	file_form, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileBytes, _ := io.ReadAll(file_form)

	result := services.AddPicture(user, fileBytes, header.Header["Content-Type"][0], header.Filename)
	if result {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "", http.StatusBadRequest)
	}
}

func GetPicture(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	user_id, ok := strconv.Atoi(queryParams.Get("user_id"))
	if ok != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	picture := services.GetImageByUser(user_id)
	if picture != nil {
		w.Header().Set("Content-Type", *picture.ContentType)
		w.Header().Set("Content-Type", strconv.Itoa(len(*picture.Picture)))
		w.Write(*picture.Picture)
		return
	} else {
		http.Error(w, "", http.StatusNotFound)
		return
	}
}

func AddAddress(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	user_id, ok := strconv.Atoi(queryParams.Get("user_id"))
	if ok != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var model models.Address
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&model); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user := services.GetUser(user_id, nil)
	address, result := services.AddAddress(user, model)
	if result {
		jsonBytes, err := json.Marshal(address)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	} else {
		http.Error(w, "", http.StatusBadRequest)
	}
}

func GetAddresses(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	user_id, ok := strconv.Atoi(queryParams.Get("user_id"))
	if ok != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	addresses := services.GetAddressByUser(user_id)
	if len(*addresses) > 0 {
		jsonBytes, err := json.Marshal(addresses)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	} else {
		http.Error(w, "", http.StatusBadRequest)
	}
}

func DeleteAddress(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	user_id, ok := strconv.Atoi(queryParams.Get("user_id"))
	if ok != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	address_id, err := strconv.Atoi(queryParams.Get("address_id"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	address := services.GetAddress(address_id)
	if address.UserId == user_id {

		result := services.DeleteAddress(address_id)
		if result {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Invalid token", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Don't belong to the user", http.StatusNotAcceptable)
	}
}

func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	var model models.Address
	defer r.Body.Close()

	queryParams := r.URL.Query()
	user_id, ok := strconv.Atoi(queryParams.Get("user_id"))
	if ok != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&model); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if model.UserId == user_id {

		result := services.UpdateAddress(model)
		if result {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Don't belong to the user", http.StatusNotAcceptable)
	}
}
