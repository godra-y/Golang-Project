package main

import (
	"database/sql"
	"encoding/json"
	"goproject/model"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	_`github.com/lib/pq`
)

var (
	// Initialize a session store
	// Replace "your-secret-key" with a strong key
	store = sessions.NewCookieStore([]byte("blablabla"))
)

func main() {
	connStr := "user=postgres dbname=goproject password='postgrass' sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	userModel := &model.UserModel{DB: db}
	productModel := &model.ProductModel{DB: db}
	categoryModel := &model.CategoryModel{DB: db}
	orderModel := &model.OrderModel{DB: db}

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := userModel.Authenticate(credentials.Username, credentials.Password)
		if err != nil {
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}

		// Create a session and save the user ID in the session
		session, _ := store.Get(r, "session-name")
		session.Values["user_id"] = user.ID
		session.Save(r, w)

		json.NewEncoder(w).Encode(map[string]string{"status": "logged in"})
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			ids, ok := r.URL.Query()["id"]

			if !ok || len(ids[0]) < 1 {
				http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(ids[0])
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}

			user, err := userModel.GetUserByID(id)
			if err != nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode(user)

		case "POST":
			var user model.User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			id, err := userModel.CreateUser(&user)
			if err != nil {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]int{"id": id})

		case "PUT":
			ids, ok := r.URL.Query()["id"]

			if !ok || len(ids[0]) < 1 {
				http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(ids[0])
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}

			var user model.User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			user.ID = id

			if err := userModel.UpdateUser(&user); err != nil {
				http.Error(w, "Failed to update user", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(map[string]string{"status": "success"})

		case "DELETE":
			ids, ok := r.URL.Query()["id"]

			if !ok || len(ids[0]) < 1 {
				http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(ids[0])
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}

			if err := userModel.DeleteUser(id); err != nil {
				http.Error(w, "Failed to delete user", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		}
	})

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			ids, ok := r.URL.Query()["id"]

			if !ok || len(ids[0]) < 1 {
				http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(ids[0])
			if err != nil {
				http.Error(w, "Invalid product ID", http.StatusBadRequest)
				return
			}

			product, err := productModel.Get(id)
			if err != nil {
				http.Error(w, "Product not found", http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode(product)

		case "POST":
			var product model.Product
			if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			id, err := productModel.Insert(&product)
			if err != nil {
				http.Error(w, "Failed to create product", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]int{"id": id})

		case "PUT":
			ids, ok := r.URL.Query()["id"]

			if !ok || len(ids[0]) < 1 {
				http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(ids[0])
			if err != nil {
				http.Error(w, "Invalid product ID", http.StatusBadRequest)
				return
			}

			var product model.Product
			if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			product.ID = id

			if err := productModel.Update(&product); err != nil {
				http.Error(w, "Failed to update product", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(map[string]string{"status": "success"})

		case "DELETE":
			ids, ok := r.URL.Query()["id"]

			if !ok || len(ids[0]) < 1 {
				http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(ids[0])
			if err != nil {
				http.Error(w, "Invalid product ID", http.StatusBadRequest)
				return
			}

			if err := productModel.Delete(id); err != nil {
				http.Error(w, "Failed to delete product", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		}
	})

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			ids, ok := r.URL.Query()["id"]

			if !ok || len(ids[0]) < 1 {
				http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(ids[0])
			if err != nil {
				http.Error(w, "Invalid category ID", http.StatusBadRequest)
				return
			}

			category, err := categoryModel.Get(id)
			if err != nil {
				http.Error(w, "Category not found", http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode(category)

		case "POST":
			var category model.Category
			if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			id, err := categoryModel.Insert(&category)
			if err != nil {
				http.Error(w, "Failed to create category", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]int{"id": id})

		case "PUT":
			ids, ok := r.URL.Query()["id"]

			if !ok || len(ids[0]) < 1 {
				http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(ids[0])
			if err != nil {
				http.Error(w, "Invalid category ID", http.StatusBadRequest)
				return
			}

			var category model.Category
			if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			category.ID = id

			if err := categoryModel.Update(&category); err != nil {
				http.Error(w, "Failed to update category", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(map[string]string{"status": "success"})

		case "DELETE":
			ids, ok := r.URL.Query()["id"]

			if !ok || len(ids[0]) < 1 {
				http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(ids[0])
			if err != nil {
				http.Error(w, "Invalid category ID", http.StatusBadRequest)
				return
			}

			if err := categoryModel.Delete(id); err != nil {
				http.Error(w, "Failed to delete category", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		}
	})

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			ids, ok := r.URL.Query()["id"]

			if !ok || len(ids[0]) < 1 {
				http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(ids[0])
			if err != nil {
				http.Error(w, "Invalid order ID", http.StatusBadRequest)
				return
			}

			order, err := orderModel.GetOrder(id)
			if err != nil {
				http.Error(w, "Order not found", http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode(order)

		case "POST":
			session, _ := store.Get(r, "session-name")
			userID, ok := session.Values["user_id"]
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var order model.Order
			if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Set the user ID from the session
			order.UserID = userID.(int)

			// Create the order with the user ID from the session
			id, err := orderModel.CreateOrder(&order)
			if err != nil {
				http.Error(w, "Failed to create order", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]int{"id": id})

		case "PUT":
			ids, ok := r.URL.Query()["id"]

			if !ok || len(ids[0]) < 1 {
				http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(ids[0])
			if err != nil {
				http.Error(w, "Invalid order ID", http.StatusBadRequest)
				return
			}

			var order model.Order
			if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			order.ID = id

			if err := orderModel.UpdateOrder(&order); err != nil {
				http.Error(w, "Failed to update order", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(map[string]string{"status": "success"})

		case "DELETE":
			ids, ok := r.URL.Query()["id"]

			if !ok || len(ids[0]) < 1 {
				http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(ids[0])
			if err != nil {
				http.Error(w, "Invalid order ID", http.StatusBadRequest)
				return
			}

			if err := orderModel.DeleteOrder(id); err != nil {
				http.Error(w, "Failed to delete order", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		}
	})

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
