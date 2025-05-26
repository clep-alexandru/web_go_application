package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	fmt.Println("DB_HOST:", dbHost)
	fmt.Println("DB_PORT:", dbPort)
	fmt.Println("DB_USER:", dbUser)
	fmt.Println("DB_PASSWORD:", dbPassword)
	fmt.Println("DB_NAME:", dbName)

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("Una sau mai multe variabile de mediu pentru baza de date nu sunt setate")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Eroare la conectarea la baza de date:", err)
	}
	defer db.Close()

	http.HandleFunc("/login", PaginaLogareHandler)
	http.HandleFunc("/signup", PaginaCreareContHandler)

	//Preia fisiere statice din "static"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server pornit pe portul 8080")
	http.ListenAndServe(":8080", nil)
}

func PaginaCreareContHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		parola := r.FormValue("parola")

		hash, err := bcrypt.GenerateFromPassword([]byte(parola), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Eroare la criptarea parolei", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, string(hash))
		if err != nil {
			http.Error(w, "Eroare la creare cont: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Afisare formular HTML
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func PaginaLogareHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var parolaContCriptata string
		err := db.QueryRow("SELECT password FROM users WHERE username=$1", username).Scan(&parolaContCriptata)
		if err != nil {
			fmt.Fprintln(w, "User inexistent sau parola gresita")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(parolaContCriptata), []byte(password))
		if err != nil {
			fmt.Fprintln(w, "Parola incorecta")
			return
		}

		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
