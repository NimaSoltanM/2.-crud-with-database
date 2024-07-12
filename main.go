package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Album struct {
	gorm.Model
	ID     string
	Title  string
	Artist string
	Price  float64
}

var album Album
var albums []Album

func main() {
	mux := http.NewServeMux()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&Album{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	mux.HandleFunc("GET /albums", func(w http.ResponseWriter, r *http.Request) {
		result := db.Find(&albums)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response := struct {
			Data []Album `json:"data"`
		}{
			Data: albums,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	mux.HandleFunc("GET /albums/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Extract the id from the URL
		id := r.PathValue("id")

		var album Album
		result := db.First(&album, "id = ?", id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Album not found", http.StatusNotFound)
			} else {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := struct {
			Data Album `json:"data"`
		}{
			Data: album,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	mux.HandleFunc("POST /albums", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&album)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := db.Create(&album)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(album)
	})

	http.ListenAndServe(":8080", mux)
}
