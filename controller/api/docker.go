package api

import (
	"net/http"

	log"github.com/Sirupsen/logrus"
	"net/url"
)

func(a* API) dockerHandler(w http.ResponseWriter, r *http.Request){
	log.Infoln(r.Method, r.URL)
	var err error
	r.URL, err = url.ParseRequestURI(a.dockerUrl)
	r.URL.Scheme = "http"
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fwd := a.fwd
	fwd.ServeHTTP(w, r)
}
