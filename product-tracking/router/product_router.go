package router

import (
	"encoding/json"
	"net/http"
	"product-tracking/model"
	"product-tracking/service"
	"strconv"
)

type Router struct {
	ProductSvc *service.ProductService
}

func (u *Router) CreateProduct(w http.ResponseWriter, r *http.Request) {

	/*user_id, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}*/

	var model model.Product
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&model); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result := u.ProductSvc.Create(model)
	if result != nil {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
}

func (u *Router) GetProduct(w http.ResponseWriter, r *http.Request) {

	/*user_id, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}*/

	queryParams := r.URL.Query()
	product_id, err := strconv.Atoi(queryParams.Get("id"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	history, err := strconv.ParseBool(queryParams.Get("history"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	product := u.ProductSvc.GetById(product_id, history)
	if product != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(product); err != nil {
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}
	} else {
		http.Error(w, "", http.StatusNotFound)
		return
	}
}

func (u *Router) GetProducts(w http.ResponseWriter, r *http.Request) {

	/*user_id, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}*/

	queryParams := r.URL.Query()
	history, err := strconv.ParseBool(queryParams.Get("history"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	products := u.ProductSvc.GetProducts(history)
	if products != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(products); err != nil {
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}
	} else {
		http.Error(w, "", http.StatusNotFound)
		return
	}
}

func (u *Router) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	/*user_id, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}*/

	queryParams := r.URL.Query()
	product_id, err := strconv.Atoi(queryParams.Get("id"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result := u.ProductSvc.DeleteProduct(product_id)
	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func (u *Router) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	/*user_id, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}*/

	var model model.Product
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&model); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result := u.ProductSvc.UpdateProduct(model)
	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
}

func (u *Router) UpdateProductPrice(w http.ResponseWriter, r *http.Request) {

	/*user_id, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}*/

	queryParams := r.URL.Query()
	product_id, err := strconv.Atoi(queryParams.Get("id"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	price, err := strconv.Atoi(queryParams.Get("price"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result := u.ProductSvc.UpdateProductPrice(product_id, price)
	if result {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
}
