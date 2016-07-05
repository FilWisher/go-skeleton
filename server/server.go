package main

// TODO: implement sessions

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"fmt"
	"html/template"
)

var db *sql.DB
var err error

type Item struct {
	Id int
	Title, Url string
}

var homeT = template.Must(template.ParseFiles("./site/index.html"))
var itemT = template.Must(template.ParseFiles("./site/item.html"))

func homePage(res http.ResponseWriter, req *http.Request) {

	var items []Item
	
	rows, err := db.Query("SELECT id, title FROM items LIMIT 20")
	if err != nil {
		// TODO: more than just panic
		panic(err.Error())
	}

	for rows.Next() {
		var id int
		var title string
		if err = rows.Scan(&id, &title); err != nil {
			// TODO: more than just panic
			panic(err.Error())	
		}
		items = append(items, Item{Id:id, Title:title})
	}
	
	err = homeT.Execute(res, items)
	if err != nil {
		panic(err.Error())	
	}
}

func init() {
	db, err = sql.Open("sqlite3", "./data.db")
	
	if err != nil {
		panic(err.Error())
	}
}

func create(db *sql.DB) http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.ServeFile(res, req, "./site/create.html")
			return
		} else {
		
			title := req.FormValue("title")
			url := req.FormValue("url")
			
			// TODO: check title and url aren't null
		
			// TODO: add user id once sessions are implemented
			result, err := db.Exec("INSERT INTO items(title, url) VALUES(?, ?)", title, url)
			if err != nil {
				panic(err.Error())
			}
		
			id, err := result.LastInsertId()
			if err != nil {
				http.Redirect(res, req, "/", 301)
			} else {
				http.Redirect(res, req, fmt.Sprintf("/ktem?id=%d",id), 301)
			}
		}
	}
}

func item(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		id_r := req.FormValue("id")
		// TODO: check id not ""
	
		var title string	
		var url string	
		var id int	
		err := db.QueryRow("SELECT id, title, url FROM items WHERE id=?", id_r).Scan(&id, &title, &url)
		if err != nil {
			http.NotFound(res, req)
			return
		}
		
		err = itemT.Execute(res, Item{Id: id, Title: title, Url: url})
		if err != nil {
			panic(err.Error())	
		}
	}
}

func main() {
	
	err = db.Ping()
	if err != nil {
		panic(err.Error())	
	}

	http.HandleFunc("/", homePage)
	http.HandleFunc("/login", login(db))
	http.HandleFunc("/signup", signup(db))
	http.HandleFunc("/create", create(db))
	http.HandleFunc("/item", item(db))
	http.ListenAndServe(":8080", nil)
}
