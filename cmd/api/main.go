package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/godra-y/go-project/pkg/api/handlers"
	"github.com/godra-y/go-project/pkg/api/models"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8080", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgresql://postgres:1@localhost/data_go?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{config: cfg}

	categoryModel := models.CategoryModel{DB: db}
	productModel := models.ProductModel{DB: db}
	userModel := models.UserModel{DB: db}
	orderModel := models.OrderModel{DB: db}

	categoryHandler := &handlers.CategoryHandler{CategoryModel: categoryModel}
	productHandler := &handlers.ProductHandler{ProductModel: productModel}
	userHandler := &handlers.UserHandler{UserModel: userModel}
	orderHandler := &handlers.OrderHandler{OrderModel: orderModel}

	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

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

	log.Printf("Starting server on %s\n", app.config.port)
	err = http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
