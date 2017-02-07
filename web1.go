package main

import (
	"io"
	"net/http"
	"encoding/json"
	"fmt"
)

type MyHandle struct{}

type Account struct{
	User string	`json:"u"`
	Psw string	`json:"p"`

}
func main() {
	var a *Account


	s := `{"u":"lwz","p":"123"}`
	json.Unmarshal([]byte(s), &a)

	fmt.Printf("%+v",a)
	var jsonBlob = [ ] byte ( `
        { "Name" : "Platypus" , "Order" : "Monotremata" } ` )
	type Animal struct {
		Name  string
		Order string
	}
	var animal  *Animal
	err := json. Unmarshal ( jsonBlob , & animal )
	if err != nil {
		fmt. Println ( "error:" , err )
	}
	fmt. Printf ( "%+v" , animal )
}

func (*MyHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL"+r.URL.String())
}