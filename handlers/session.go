package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Session struct {
	l *log.Logger
}

func NewSession(l *log.Logger) *Session {
	return &Session{l}
}

type Mystruct struct {
	Email  string  `json:"email"`
	Points float64 `json:"points"`
	Solved []int   `json:"solved"`
}

func (q *Mystruct) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(q)
}

func (q *Session) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		q.session(rw, r)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (q *Session) session(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS"))
	c, _ := r.Cookie("accessToken")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	var email string
	var points float64
	id := c.String()[12:]
	err = db.QueryRow("SELECT email FROM sessions where id=$1", id).Scan(&email)
	err = db.QueryRow("SELECT points from users where email=$1", email).Scan(&points)
	res, err := db.Query("SELECT qid from solved where email=$1", email)
	var a []int
	for res.Next() {
		var ele int
		res.Scan(&ele)
		a = append(a, ele)
	}
	switch {
	case err == sql.ErrNoRows:
		http.Error(rw, "Failed to login", http.StatusForbidden)
	case err != nil:
		log.Fatal(err)
	default:
		ele := &Mystruct{
			Email:  email,
			Points: points,
			Solved: a,
		}
		ele.ToJSON(rw)

	}
}
