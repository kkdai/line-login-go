package main

import (
	"html/template"
	"log"
	"net/http"
)

func browse(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login.tmpl"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("Template err:", err)
	}
}
