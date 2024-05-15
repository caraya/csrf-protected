package main

import (
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

// Middleware to enable CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// Your webhook logic here
	w.Write([]byte("Webhook received"))
}

func main() {
	r := mux.NewRouter()

	// CSRF protection middleware
	csrfMiddleware := csrf.Protect(
		[]byte("32-byte-long-auth-key"), // Replace with your own 32-byte key
		csrf.Secure(false),              // Set to true in production (requires HTTPS)
		csrf.Path("/"),                  // Set the path where the CSRF token cookie is valid
	)

	r.HandleFunc("/webhook", webhookHandler).Methods("POST")

	// Setting CSRF token as a cookie for demo purposes
	r.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(csrf.Token(r)))
	}).Methods("GET")

	// Apply CORS and CSRF middleware
	handler := enableCORS(csrfMiddleware(r))

	http.Handle("/", handler)
	log.Println("Server started at :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
