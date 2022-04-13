package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/rs/cors"
)

// func IndexHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "Hello world")
// }

func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	n := new(Note)

	notes, err := n.GetAll()

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j, err := json.Marshal(notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")

	w.Write(j)
}

func CreateNotesHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = note.Create()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

func UpdateNotesHandler(w http.ResponseWriter, r *http.Request) {
	var note Note

	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = note.Update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteNotesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}
	var note Note

	err = note.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func NotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetNotesHandler(w, r)
	case "POST":
		CreateNotesHandler(w, r)
	case "PUT":
		UpdateNotesHandler(w, r)
	case "DELETE":
		DeleteNotesHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/notes", NotesHandler)

	migrate := flag.Bool("migrate", false, "Run migrations")
	flag.Parse()
	if *migrate {
		if err := MakeMigrations(); err != nil {
			log.Fatal(err)
		}
	}
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	})
	core := c.Handler(router)
	http.ListenAndServe(":3000", core)
}
