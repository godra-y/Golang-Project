package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	v1 := r.PathPrefix("/api/v1").Subrouter()

	//Category routes
	v1.HandleFunc("/categories", app.getCategoriesList).Methods("GET")
	v1.HandleFunc("/categories", app.createCategoryHandler).Methods("POST")
	v1.HandleFunc("/categories/{id}", app.getCategoryHandler).Methods("GET")
	v1.HandleFunc("/categories/{id}", app.updateCategoryHandler).Methods("PUT")
	v1.HandleFunc("/categories/{id}", app.deleteCategoryHandler).Methods("DELETE")

	//Product routes
	v1.HandleFunc("/products", app.getProductsList).Methods("GET")
	v1.HandleFunc("/products", app.createProductHandler).Methods("POST")
	v1.HandleFunc("/products/{id}", app.getProductHandler).Methods("GET")
	v1.HandleFunc("/products/{id}", app.updateProductHandler).Methods("PUT")
	v1.HandleFunc("/products/{id}", app.deleteProductHandler).Methods("DELETE")
	v1.HandleFunc("/categories/{id}/products", app.getProductsByCategoryHandler).Methods("GET")

	//User Routes
	v1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	v1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	v1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	return app.authenticate(r)
}
