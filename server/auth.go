package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func signup(db *sql.DB) http.HandlerFunc {
	return func (res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.ServeFile(res, req, "./site/signup.html")
			return
		}
		
		username := req.FormValue("username")
		password := req.FormValue("password")
	
		var user string
		err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)
		
		switch {
		case err == sql.ErrNoRows:
			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			
			_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashPassword)
			if err != nil {
				panic(err.Error())	
			}
			
			res.Write([]byte("User created"))
			return
		case err != nil:
			panic(err.Error())
		default:
			http.Redirect(res, req, "/", 301)
		}
	}
}

func login(db *sql.DB) http.HandlerFunc {
	return func (res http.ResponseWriter, req *http.Request) {

		if req.Method != "POST" {
			http.ServeFile(res, req, "./site/login.html")
			return
		}
		
		username := req.FormValue("username")
		password := req.FormValue("password")
		
		var databaseUsername string
		var databasePassword string
		
		err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)
		if err != nil {
			http.Redirect(res, req, "/login", 301)
			return
		}
		
		err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
		
		if err != nil {
			http.Redirect(res, req, "/login", 301)
			return
		}
		
		res.Write([]byte("Hello " + databaseUsername))
	}
}
