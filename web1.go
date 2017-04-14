package main

import (
	"io"
	"net/http"
	//"encoding/json"
	//"fmt"

	"fmt"
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
	s:= "haha"
	a := string("haha")
	fmt.Println(s == a)
}

func (*MyHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL"+r.URL.String())

}