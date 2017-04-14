package api

import(
	"net/http"
	log"github.com/Sirupsen/logrus"
	"encoding/json"
	"github.com/gorilla/mux"
)

func(a *API) roles(w http.ResponseWriter, r *http.Request){

	log.Debugln(r.Method, r.URL)
	w.Header().Set("content-type", "application/json")
	roles := a.manager.Roles()

	if err := json.NewEncoder(w).Encode(&roles); err != nil {
		log.Fatalln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func(a *API) role(w http.ResponseWriter, r *http.Request){
	log.Debugln(r.Method, r.URL)
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	name := vars["name"]

	role := a.manager.Role(name)
	if role == nil{
		http.Error(w, "role does not exist", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&role); err != nil {
		log.Fatalln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
