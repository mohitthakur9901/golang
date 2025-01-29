package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)

}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	for _, movie := range movies {
		if movie.ID == param["id"] {
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movie Movie
	_ = json.NewDecoder(r.Body).Decode((&movie))
	movie.ID = strconv.Itoa((rand.Intn(100000000)))
	for _, item := range movies {
		if item.Title == movie.Title {
			json.NewEncoder(w).Encode("Movie Already Exisit")
			return
		} else {
			movies = append(movies, movie)

		}
	}
	json.NewEncoder(w).Encode(movie)

}


func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// Decode the request body into a struct
	var updatedMovie Movie
	err := json.NewDecoder(r.Body).Decode(&updatedMovie)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find and update the movie
	for i, movie := range movies {
		if movie.ID == params["id"] {
			movies[i].Title = updatedMovie.Title
			movies[i].Isbn = updatedMovie.Isbn
			movies[i].Director = updatedMovie.Director

			json.NewEncoder(w).Encode(movies[i])
			return
		}
	}

	// If the movie is not found
	http.Error(w, "Movie not found", http.StatusNotFound)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	found := false

	for index, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			found = true
			break
		}
	}
	if !found {
		json.NewEncoder(w).Encode("NOT FOUND")
	}
	json.NewEncoder(w).Encode(movies)
}

func main() {

	r := mux.NewRouter()

	movies = append(movies, Movie{
		ID:    "1",
		Isbn:  "123444",
		Title: "Hello",
		Director: &Director{
			FirstName: "Mohit",
			LastName:  "Thakur",
		},
	})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting Port Server Is 4000")
	log.Fatal(http.ListenAndServe(":3000", r))

}
