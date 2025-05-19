package main

import (
	"fmt"
	"log"
	"net/http"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow specific headers and methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func main() {
	http.HandleFunc("/chunk", enableCORS(chunkHandler))
	http.HandleFunc("/lookup", enableCORS(lookupHandler))
	http.HandleFunc("/create-collection", enableCORS(createCollectionHandler))
	http.HandleFunc("/delete-collection", enableCORS(deleteCollectionHandler))
	http.HandleFunc("/delete-point", enableCORS(deletePointHandler))
	fmt.Println("🧠 Chisel API running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
