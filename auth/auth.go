
package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

var (
	ErrUnauthorized  = errors.New("unauthorized")
	ErrNoUserInToken = errors.New("no user sent in token")
)

type (
	Account struct {
		ID        string       `json:"id,omitempty" gorethink:"id,omitempty"`
		Username  string       `json:"username,omitempty" gorethink:"username"`
		Password  string       `json:"password,omitempty" gorethink:"password"`
		Tokens    []*AuthToken `json:"-" gorethink:"tokens"`
		Roles     []string     `json:"roles,omitempty" gorethink:"roles"`
	}

	AuthToken struct {
		Token     string `json:"auth_token,omitempty" gorethink:"auth_token"`
		UserAgent string `json:"user_agent,omitempty" gorethink:"user_agent"`
	}

	AccessToken struct {
		Token    string
		Username string
	}

	ServiceKey struct {
		Key         string `json:"key,omitempty" gorethink:"key"`
		Description string `json:"description,omitempty" gorethink:"description"`
	}

	Authenticator interface {
		Authenticate(username, password, hash string) (bool, error)
		GenerateToken() (string, error)
		IsUpdateSupported() bool
		Name() string
	}
)

func Hash(data string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	return string(h[:]), err
}

func GenerateToken() (string, error) {
	return Hash(time.Now().String())
}

// GetAccessToken returns an AccessToken from the access header
func GetAccessToken(authToken string) (*AccessToken, error) {
	parts := strings.Split(authToken, ":")

	if len(parts) != 2 {
		return nil, ErrNoUserInToken

	}

	return &AccessToken{
		Username: parts[0],
		Token:    parts[1],
	}, nil

}

//ContainerPool Authentication
type CPAuthenticator struct {
		name string
}


func NewAuthenticator(name string) Authenticator {
	return &CPAuthenticator{
		name : name,
	}
}

func (a CPAuthenticator) IsUpdateSupported() bool {
	return true
}

func (a CPAuthenticator) Name() string {
	return a.name
}

func (a CPAuthenticator) Authenticate(username, password, hash string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err == nil {
		return true, nil
	}
	return false, nil
}

func (a CPAuthenticator) GenerateToken() (string, error) {
	return GenerateToken()
}