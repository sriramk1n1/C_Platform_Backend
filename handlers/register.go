package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Register struct {
	l *log.Logger
}

func NewRegister(l *log.Logger) *Register {
	return &Register{l}
}

func (q *Register) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS"))

	if r.Method == http.MethodPost {
		q.register(rw, r)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (q *Register) register(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS"))

	email := r.FormValue("email")
	password := r.FormValue("password")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	_, err = db.Exec("INSERT INTO users (email,password) values ($1,$2)", email, password)
	if err != nil {
		rw.Write([]byte("User already Exists"))
	}
}
