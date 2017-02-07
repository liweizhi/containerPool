package middleware

import(
	"net/http"
	"github.com/liweizhi/containerPool/controller/manager"
	"strings"
	"fmt"
	log"github.com/Sirupsen/logrus"
)
type AuthRequired struct {
	deniedHostHandler http.HandlerFunc
	manager           manager.Manager

}
func defaultDeniedHostHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}
func NewAuthRequired(m manager.Manager) *AuthRequired {
	return &AuthRequired{
		deniedHostHandler: http.HandlerFunc(defaultDeniedHostHandler),
		manager:           m,

	}
}

func (a *AuthRequired)handleRequest(w http.ResponseWriter, r *http.Request) error{
	valid := false
	// service key takes priority
	serviceKey := r.Header.Get("X-Service-Key")
	if serviceKey != "" {
		if err := a.manager.VerifyServiceKey(serviceKey); err == nil {
			valid = true
		}
	} else { // check for authHeader
		authHeader := r.Header.Get("X-Access-Token")
		parts := strings.Split(authHeader, ":")
		if len(parts) == 2 {
			// validate
			user := parts[0]
			token := parts[1]
			if err := a.manager.VerifyAuthToken(user, token); err == nil {
				valid = true

			}
		}
	}

	if !valid {
		a.deniedHostHandler(w, r)
		return fmt.Errorf("unauthorized %s", r.RemoteAddr)
	}

	return nil
}
func (a *AuthRequired) HandlerFuncWithNext(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := a.handleRequest(w, r)
	log.Infoln("authrequire")
	if err != nil {
		log.Warnf("unauthorized request for %s from %s", r.URL.Path, r.RemoteAddr)
		return
	}

	if next != nil {
		next(w, r)
	}
}