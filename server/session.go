package main

import (
	"fmt"
	"net/http"
)

func checkSession(h http.HandlerFunc) (http.HandlerFunc) {
	return func(res http.ResponseWriter, req *http.Request) {
	
		c, err := req.Cookie("test")

		h(res, req)
	}
}
