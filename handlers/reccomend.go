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

type Recommend struct {
	l *log.Logger
}

func NewRecommend(l *log.Logger) *Recommend {
	return &Recommend{l}
}

type mystruct_list struct {
	List []mystruct `json:"list"`
}
type mystruct struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func (q *mystruct_list) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(q)
}

func (q *Recommend) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		q.recommend(rw, r)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (q *Recommend) recommend(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	c, _ := r.Cookie("accessToken")
	id := c.String()[12:]

	rw.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS"))
	var points int
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	if err := db.Ping(); err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	defer db.Close()
	db.QueryRow("select points from users where email in (select email from sessions where id = $1);", id).Scan(&points)
	res, err := db.Query("select q.id,q.name as distance from question q order by abs(q.level-$1) limit 5;", points)
	li := &mystruct_list{}
	for res.Next() {
		var ele mystruct
		res.Scan(&ele.Id, &ele.Name)
		li.List = append(li.List, ele)
	}
	switch {
	case err == sql.ErrNoRows:
		http.Error(rw, "Failed to login", http.StatusForbidden)
	case err != nil:
		log.Fatal(err)
	default:
		li.ToJSON(rw)
	}
}
