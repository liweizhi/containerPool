package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"strings"

)

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var creds *Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginSuccessful, err := a.manager.Authenticate(creds.Username, creds.Password)
	if err != nil {
		log.Errorf("error during login for %s from %s: %s", creds.Username, r.RemoteAddr, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !loginSuccessful {
		log.Warnf("invalid login for %s from %s", creds.Username, r.RemoteAddr)
		http.Error(w, "invalid username/password", http.StatusForbidden)
		return
	}




	// return token
	token, err := a.manager.NewAuthToken(creds.Username, r.UserAgent())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	if err := json.NewEncoder(w).Encode(token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *API) changePassword(w http.ResponseWriter, r *http.Request) {

	var creds *Credentials
	var username string
	authHeader := r.Header.Get("X-Access-Token")
	parts := strings.Split(authHeader, ":")
	if len(parts) == 2 {
		username = parts[0]
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.manager.ChangePassword(username, creds.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
