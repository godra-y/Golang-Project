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

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func NewApp(db *sql.DB) *App {
	return &App{
		Router: mux.NewRouter(),
		DB:     db,
	}
}

func (app *App) Initialize() {
	app.initializeRoutes()
}

func (app *App) initializeRoutes() {
	categoryModel := &models.CategoryModel{DB: app.DB}
	productModel := &models.ProductModel{DB: app.DB}

	categoryHandler := &handlers.CategoryHandler{CategoryModel: categoryModel}
	productHandler := &handlers.ProductHandler{ProductModel: productModel}

	// Categories
	app.Router.HandleFunc("/api/v1/categories", categoryHandler.CreateCategory).Methods("POST")
	app.Router.HandleFunc("/api/v1/categories/{id:[0-9]+}", categoryHandler.GetCategory).Methods("GET")
	app.Router.HandleFunc("/api/v1/categories/{id:[0-9]+}", categoryHandler.UpdateCategory).Methods("PUT")
	app.Router.HandleFunc("/api/v1/categories/{id:[0-9]+}", categoryHandler.DeleteCategory).Methods("DELETE")

	// Products
	app.Router.HandleFunc("/api/v1/products", productHandler.CreateProduct).Methods("POST")
	app.Router.HandleFunc("/api/v1/products/{id:[0-9]+}", productHandler.GetProduct).Methods("GET")
	app.Router.HandleFunc("/api/v1/products/{id:[0-9]+}", productHandler.UpdateProduct).Methods("PUT")
	app.Router.HandleFunc("/api/v1/products/{id:[0-9]+}", productHandler.DeleteProduct).Methods("DELETE")
}

func main() {
	var dbURI string
	flag.StringVar(&dbURI, "db", "postgres://postgres:1@localhost:5432/data?sslmode=disable", "Database URI")
	flag.Parse()

	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Fatal("Error closing the database connection:", cerr)
		}
	}()

	app := NewApp(db)
	app.Initialize()

	server := &http.Server{
		Addr:    ":8080",
		Handler: app.Router,
	}

	log.Println("Server listening on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
