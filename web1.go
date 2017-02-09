package main

import (
	"io"
	"net/http"
	//"encoding/json"
	//"fmt"
	"log"
)

type MyHandle struct{}

type Account struct{
	User string	`json:"u"`
	Psw string	`json:"p"`

}
func redirect(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "www.sina.com", 301)
}

func main() {

	http.HandleFunc("/", redirect)
	err := http.ListenAndServe(":9010", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (*MyHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL"+r.URL.String())

}