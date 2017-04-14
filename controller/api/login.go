package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	log.Debugln(r.Method, r.URL)
	var creds *Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginSuccessful, err := a.manager.Authenticate(creds.Username, creds.Password)


	if !loginSuccessful {
		log.Warnf("invalid login for %s from %s", creds.Username, r.RemoteAddr)
		http.Error(w, "invalid username/password", http.StatusForbidden)
		return
	}
	if err != nil {
		log.Errorf("error during login for %s from %s: %s", creds.Username, r.RemoteAddr, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}




	// return token
	token, err := a.manager.NewAuthToken(creds.Username, r.UserAgent())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := a.manager.Store().Get(r, a.manager.StoreKey())
	session.Values["username"] = creds.Username
	session.Save(r, w)

	if err := json.NewEncoder(w).Encode(token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}



}

func (a *API) changePassword(w http.ResponseWriter, r *http.Request) {

	var creds *Credentials
	var username string

	seesion, _ := a.manager.Store().Get(r, a.manager.StoreKey())
	username = seesion.Values["username"].(string)
	if username == "" {
		http.Error(w, "unauthroized", http.StatusUnauthorized)
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
