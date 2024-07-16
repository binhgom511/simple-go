package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-movies-crud/models"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Response struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type Movie struct {
	Isbn  string `json:"isbn"`
	Title string `json:"title"`
	Price string `json:"price"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to postgres")
	return db
}

// GetMovies fetches all movies
// @Summary Get all movies
// @Description Retrieves a list of all movies from the database
// @Tags movies
// @Accept json
// @Produce json
// @Success 200 {array} Movie
// @Router /movies [get]
func GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := getAllMovies()

	if err != nil {
		log.Fatalf("Unable to get all movies. %v", err)
	}

	json.NewEncoder(w).Encode(movies)
}

// GetMovie retrieves a single movie by ID
// @Summary Get a movie by ID
// @Description Retrieves a movie from the database by its ID
// @Tags movies
// @Accept json
// @Produce json
// @Param id path string true "Movie ID"
// @Success 200 {object} models.Movie
// @Router /movies/{id} [get]
func GetMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	movie, err := getMovie(id)

	if err != nil {
		log.Fatalf("Unable to get movie. %v", err)
	}

	json.NewEncoder(w).Encode(movie)
}

// CreateMovies creates a new movie
// @Summary Create a new movie
// @Description Creates a new movie in the database
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body Movie true "Movie object to be created"
// @Success 200 {object} Response
// @Router /movies [post]
func CreateMovies(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)

	if err != nil {
		log.Fatal("Unabel to decode the request body. %v", err)
	}

	insertID := insertMovie(movie)

	res := Response{
		ID:      insertID,
		Message: "Movie created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

// UpdateMovie updates an existing movie by ID
// @Summary Update a movie by ID
// @Description Updates an existing movie in the database by its ID
// @Tags movies
// @Accept json
// @Produce json
// @Param id path string true "Movie ID"
// @Param movie body Movie true "Updated movie object"
// @Success 200 {object} Response
// @Router /movies/{id} [put]
func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var movie models.Movie

	err := json.NewDecoder(r.Body).Decode(&movie)

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	updatedRows := updateMovie(string(id), movie)

	msg := fmt.Sprintf("Movie updated successfully. Total rows/records affected %v", updatedRows)

	res := Response{
		ID:      string(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// DeleteMovie deletes a movie by ID
// @Summary Delete a movie by ID
// @Description Deletes a movie from the database by its ID
// @Tags movies
// @Accept json
// @Produce json
// @Param id path string true "Movie ID"
// @Success 200 {object} Response
// @Router /movies/{id} [delete]
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	deleteRows := deleteMovie(id)

	msg := fmt.Sprintf("Movie deleted successfully. Total rows/records %v", deleteRows)
	res := Response{
		ID:      string(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func insertMovie(movie models.Movie) string {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO movies(isbn, title, price) VALUES ($1, $2, $3) RETURNING id`

	err := db.QueryRow(sqlStatement, movie.Isbn, movie.Title, movie.Price).Scan(&movie.ID)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", movie.ID)
	return movie.ID
}

func getMovie(id string) (models.Movie, error) {
	db := createConnection()

	defer db.Close()

	var movie models.Movie

	sqlStatement := `SELECT * FROM movies WHERE id=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&movie.ID, &movie.Isbn, &movie.Title, &movie.Price)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No row were returned!")
		return movie, nil
	case nil:
		return movie, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return movie, err
}

func getAllMovies() ([]models.Movie, error) {
	db := createConnection()

	defer db.Close()

	var movies []models.Movie

	sqlStatement := `SELECT * FROM movies`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		err = rows.Scan(&movie.ID, &movie.Isbn, &movie.Title, &movie.Price)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		movies = append(movies, movie)
	}

	return movies, err
}

func updateMovie(id string, movie models.Movie) int64 {
	db := createConnection()

	defer db.Close()
	sqlStatement := `UPDATE movies SET isbn=$2, title=$3, price=$4 WHERE id=$1`

	res, err := db.Exec(sqlStatement, id, movie.Isbn, movie.Title, movie.Price)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/records affected. %v", rowsAffected)
	return rowsAffected
}

func deleteMovie(id string) int64 {
	db := createConnection()

	defer db.Close()
	sqlStatement := `DELETE FROM movies WHERE id=$1`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/records affected. %v", rowsAffected)
	return rowsAffected
}
