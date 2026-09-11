package routers

import (
	"encoding/json"
	"io"
	"net/http"

	"src/models"
)

func (u *Router) GetUser(w http.ResponseWriter, r *http.Request) {

	user_id, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user := u.UserSvc.GetUser(user_id, []string{"Address"})
	if user != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}
	} else {
		http.Error(w, "", http.StatusNotFound)
		return
	}
}

func (u *Router) DeleteUser(w http.ResponseWriter, r *http.Request) {

	user_id, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	result := u.UserSvc.DeleteUser(user_id)
	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func (u *Router) UpdateUser(w http.ResponseWriter, r *http.Request) {

	user_id, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var model models.User
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&model); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	model.Id = user_id
	result := u.UserSvc.UpdateUser(model)
	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
}

func (u *Router) ChangePassword(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	user_id, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	password := queryParams.Get("password")
	if password != "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user := u.UserSvc.GetUser(user_id, nil)
	if user == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	result := u.UserSvc.ChangePassword(user, password)
	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

}

func (u *Router) AddPicture(w http.ResponseWriter, r *http.Request) {

	user_id, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user := u.UserSvc.GetUser(user_id, nil)
	if user == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	file_form, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileBytes, _ := io.ReadAll(file_form)

	result := u.UserSvc.AddPicture(user, fileBytes, header.Header["Content-Type"][0], header.Filename)
	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (u *Router) GetPicture(w http.ResponseWriter, r *http.Request) {

	/*user_id, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	picture := services.GetImageByUser(user_id)
	if picture.Picture != nil {
		w.Header().Set("Content-Type", *picture.ContentType)
		w.Header().Set("Content-Length", strconv.Itoa(len(*picture.Picture)))
		// Instruct browsers and CDNs to cache this asset
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

		// Wrap the byte array in a Reader so ServeContent can stream it
		contentReader := bytes.NewReader(*picture.Picture)
		// ServeContent automatically detects Content-Type from the filename extension
		http.ServeContent(w, r, "", time.Now(), contentReader)
		return
	} else {
		http.Error(w, "", http.StatusNotFound)
		return
	}*/
}

func (u *Router) AddAddress(w http.ResponseWriter, r *http.Request) {

	user_id, ok := r.Context().Value("userId").(int)
	if !ok {
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

	user := u.UserSvc.GetUser(user_id, nil)
	address, result := u.UserSvc.AddAddress(user, model)
	if result {
		jsonBytes, err := json.Marshal(address)
		if err != nil {
			http.Error(w, "Address wrong parse", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
		return
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func (u *Router) GetAddresses(w http.ResponseWriter, r *http.Request) {

	/*user_id, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	addresses := services.GetAddressByUser(user_id)
	if len(*addresses) > 0 {
		if err := json.NewEncoder(w).Encode(addresses); err != nil {
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}*/
}

func (u *Router) DeleteAddress(w http.ResponseWriter, r *http.Request) {

	/*queryParams := r.URL.Query()
	user_id, ok := r.Context().Value("userId").(int)
	if !ok {
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
			return
		} else {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Don't belong to the user", http.StatusNotAcceptable)
		return
	}*/
}

func (u *Router) UpdateAddress(w http.ResponseWriter, r *http.Request) {

	/*user_id, ok := r.Context().Value("userId").(int)
	if !ok {
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

	if model.UserId == user_id {

		result := services.UpdateAddress(model)
		if result {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Don't belong to the user", http.StatusNotAcceptable)
		return
	}*/
}
