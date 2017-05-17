package api

import "net/http"

func(a *API)info(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "application/json")

}
