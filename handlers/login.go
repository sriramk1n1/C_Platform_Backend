package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Login struct {
	l *log.Logger
}

func NewLogin(l *log.Logger) *Login {
	return &Login{l}
}

func (q *Login) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS"))
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodPost {
		q.login(rw, r)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (q *Login) login(rw http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	if err := db.Ping(); err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	defer db.Close()
	var res string
	err = db.QueryRow("SELECT email FROM users where email=$1 and password=$2", email, password).Scan(&res)
	switch {
	case err == sql.ErrNoRows:
		http.Error(rw, "Invalid Credentials", http.StatusForbidden)
	case err != nil:
		log.Fatal(err)
	default:
		id := uuid.New().String()
		db.Exec("INSERT INTO sessions(id,email) values($1,$2)", id, email)
		cookie := http.Cookie{}
		cookie.Name = "accessToken"
		cookie.Value = id
		cookie.Path = "/"
		cookie.Expires = time.Now().Add(365 * 24 * time.Hour)
		// cookie.Secure = true
		cookie.Domain = "localhost"
		// cookie.SameSite = http.SameSiteNoneMode

		http.SetCookie(rw, &cookie)
		rw.Write([]byte("success"))
	}
}
