package main

import (
	"net/http"

	"github.com/godra-y/go-project/pkg/api/handlers"
	"github.com/godra-y/go-project/pkg/api/models"
	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	v1 := r.PathPrefix("/api/v1").Subrouter()

	categoryModel := models.CategoryModel{DB: app.db}
	productModel := models.ProductModel{DB: app.db}
	userModel := models.UserModel{DB: app.db}
	orderModel := models.OrderModel{DB: app.db}

	categoryHandler := &handlers.CategoryHandler{CategoryModel: categoryModel}
	productHandler := &handlers.ProductHandler{ProductModel: productModel}
	userHandler := &handlers.UserHandler{UserModel: userModel}
	orderHandler := &handlers.OrderHandler{OrderModel: orderModel}

	// Category routes
	v1.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
	v1.HandleFunc("/categories/{id}", categoryHandler.GetCategory).Methods("GET")
	v1.HandleFunc("/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT")
	v1.HandleFunc("/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")

	// Product routes
	v1.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	v1.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	v1.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	v1.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	// User routes
	v1.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	v1.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	v1.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	v1.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Order routes
	v1.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	v1.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET")
	v1.HandleFunc("/orders/{id}", orderHandler.UpdateOrder).Methods("PUT")
	v1.HandleFunc("/orders/{id}", orderHandler.DeleteOrder).Methods("DELETE")

	return r
}
