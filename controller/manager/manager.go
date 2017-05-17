package manager


import (

	"github.com/liweizhi/containerPool/auth"
	"github.com/samalba/dockerclient"
	"github.com/gorilla/sessions"
	r"gopkg.in/dancannon/gorethink.v2"
	"errors"

	log"github.com/Sirupsen/logrus"

	"fmt"
	"github.com/gorilla/securecookie"
)
const (
	tblNameConfig      = "config"
	tblNameEvents      = "events"
	tblNameAccounts    = "accounts"
	tblNameRoles       = "roles"
	tblNameServiceKeys = "service_keys"
	tblNameExtensions  = "extensions"
	tblNameWebhookKeys = "webhook_keys"
	tblNameRegistries  = "registries"
	tblNameConsole     = "console"
	storeKey           = "thucloud"
)

var (
	ErrCannotPingRegistry         = errors.New("Cannot ping registry")
	ErrLoginFailure               = errors.New("invalid username or password")
	ErrAccountExists              = errors.New("account already exists")
	ErrAccountDoesNotExist        = errors.New("account does not exist")
	ErrRoleDoesNotExist           = errors.New("role does not exist")
	ErrNodeDoesNotExist           = errors.New("node does not exist")
	ErrServiceKeyDoesNotExist     = errors.New("service key does not exist")
	ErrInvalidAuthToken           = errors.New("invalid auth token")
	ErrExtensionDoesNotExist      = errors.New("extension does not exist")
	ErrWebhookKeyDoesNotExist     = errors.New("webhook key does not exist")
	ErrRegistryDoesNotExist       = errors.New("registry does not exist")
	ErrConsoleSessionDoesNotExist = errors.New("console session does not exist")
)
type DefaultManager struct {
	storeKey         string
	database         string
	authKey          string
	session          *r.Session
	authenticator    auth.Authenticator
	store            *sessions.CookieStore
	client           *dockerclient.DockerClient

}

type ScaleResult struct {
	Scaled []string
	Errors []string
}

type Manager interface {
	Accounts() ([]*auth.Account, error)
	Account(username string) (*auth.Account, error)
	SaveAccount(account *auth.Account) error
	DeleteAccount(account *auth.Account) error

	Authenticate(username, password string) (bool, error)
	GetAuthenticator() auth.Authenticator

	NewAuthToken(username string, userAgent string) (*auth.AuthToken, error)
	VerifyAuthToken(username, token string) error

	VerifyServiceKey(key string) error
	SaveServiceKey(key *auth.ServiceKey) error
	RemoveServiceKey(key string) error
	ServiceKey(key string) (*auth.ServiceKey, error)
	ServiceKeys() ([]*auth.ServiceKey, error)
	NewServiceKey(description string) (*auth.ServiceKey, error)

	ChangePassword(username, password string) error

	Store() *sessions.CookieStore
	StoreKey() string

	Roles()  []*auth.ACL
	Role(name string)  *auth.ACL
	//NodeInfo()

}

func NewManager(addr string, database string, authKey string, client *dockerclient.DockerClient,  authenticator auth.Authenticator) (Manager, error) {
	log.Debug("setting up rethinkdb session")
	session, err := r.Connect(r.ConnectOpts{
		Address:  addr,
		Database: database,
		AuthKey:  authKey,
	})
	if err != nil {
		return nil, err
	}
	log.Info("checking database")

	r.DBCreate(database).Run(session)
	store := sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

	m := &DefaultManager{
		database:         database,
		authKey:          authKey,
		session:          session,
		authenticator:    authenticator,
		store:            store,
		client:           client,
		storeKey:         storeKey,
	}

	m.initdb()
	return m, nil
}

func (m DefaultManager) initdb() {
	// create tables if needed
	tables := []string{tblNameConfig, tblNameEvents, tblNameAccounts, tblNameRoles, tblNameConsole, tblNameServiceKeys, tblNameRegistries, tblNameExtensions, tblNameWebhookKeys}
	for _, tbl := range tables {
		_, err := r.Table(tbl).Run(m.session)
		if err != nil {
			if _, err := r.DB(m.database).TableCreate(tbl).Run(m.session); err != nil {
				log.Fatalf("error creating table: %s", err)
			}
		}
	}
}

func (m DefaultManager) Accounts() ([]*auth.Account, error) {
	res, err := r.Table(tblNameAccounts).OrderBy(r.Asc("username")).Run(m.session)
	if err != nil {
		return nil, err
	}
	accounts := []*auth.Account{}
	if err := res.All(&accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (m DefaultManager) Account(username string) (*auth.Account, error) {
	res, err := r.Table(tblNameAccounts).Filter(map[string]string{"username": username}).Run(m.session)
	if err != nil {
		return nil, err

	}
	if res.IsNil() {
		return nil, ErrAccountDoesNotExist
	}
	var account *auth.Account
	if err := res.One(&account); err != nil {
		return nil, err
	}
	return account, nil
}

func (m DefaultManager) SaveAccount(account *auth.Account) error {
	var (
		hash      string
		eventType string
	)
	if account.Password != "" {
		h, err := auth.Hash(account.Password)
		if err != nil {
			return err
		}

		hash = h
	}
	// check if exists; if so, update
	acct, err := m.Account(account.Username)
	if err != nil && err != ErrAccountDoesNotExist {
		return err
	}

	// update
	if acct != nil {
		updates := map[string]interface{}{
			"roles":      account.Roles,
		}
		if account.Password != "" {
			updates["password"] = hash
		}

		if _, err := r.Table(tblNameAccounts).Filter(map[string]string{"username": account.Username}).Update(updates).RunWrite(m.session); err != nil {
			return err
		}

		eventType = "update-account"
	} else {
		account.Password = hash
		if _, err := r.Table(tblNameAccounts).Insert(account).RunWrite(m.session); err != nil {
			return err
		}

		eventType = "add-account"
	}

	log.Debugln(eventType)
	//m.logEvent(eventType, fmt.Sprintf("username=%s", account.Username), []string{"security"})

	return nil
}

func (m DefaultManager) DeleteAccount(account *auth.Account) error {
	log.Infoln("delete account", account.ID)
	res, err := r.Table(tblNameAccounts).Filter(map[string]string{"id": account.ID}).Delete().Run(m.session)
	if err != nil {
		return err
	}

	if res.IsNil() {
		return ErrAccountDoesNotExist
	}

	//m.logEvent("delete-account", fmt.Sprintf("username=%s", account.Username), []string{"security"})

	return nil
}

func (m DefaultManager) GetAuthenticator() auth.Authenticator {
	return m.authenticator
}

func (m DefaultManager) Authenticate(username, password string) (bool, error) {
	// only get the account to get the hashed password if using the builtin auth
	passwordHash := ""
	if m.authenticator.Name() == "builtin" {
		acct, err := m.Account(username)
		if err != nil {
			log.Error(err)
			return false, ErrLoginFailure
		}

		passwordHash = acct.Password
	}

	a, err := m.authenticator.Authenticate(username, password, passwordHash)
	if !a || err != nil {
		log.Error(ErrLoginFailure)
		return false, ErrLoginFailure
	}

	return true, nil
}


func (m DefaultManager) VerifyAuthToken(username, token string) error {
	acct, err := m.Account(username)
	if err != nil {
		return err
	}
	found := false
	for _, t := range acct.Tokens {
		if token == t.Token {
			found = true
			break
		}
	}
	if !found {
		return ErrInvalidAuthToken
	}
	return nil
}

func (m DefaultManager) NewAuthToken(username string, userAgent string) (*auth.AuthToken, error) {
	tk, err := m.authenticator.GenerateToken()
	if err != nil {
		return nil, err
	}
	acct, err := m.Account(username)
	if err != nil {
		return nil, err
	}
	token := &auth.AuthToken{}
	tokens := acct.Tokens
	found := false
	for _, t := range tokens {
		if t.UserAgent == userAgent {
			found = true
			t.Token = tk
			token = t
			break
		}
	}
	if !found {
		token = &auth.AuthToken{
			UserAgent: userAgent,
			Token:     tk,
		}
		tokens = append(tokens, token)
	}
	// delete token

	if _, err := r.Table(tblNameAccounts).Filter(map[string]string{"username": username}).Filter(r.Row.Field("user_agent").Eq(userAgent)).Delete().Run(m.session); err != nil {
		return nil, err
	}
	// add
	if _, err := r.Table(tblNameAccounts).Filter(map[string]string{"username": username}).Update(map[string]interface{}{"tokens": tokens}).RunWrite(m.session); err != nil {
		return nil, err
	}
	return token, nil
}

func (m DefaultManager) VerifyServiceKey(key string) error {
	if _, err := m.ServiceKey(key); err != nil {
		return err
	}
	return nil
}
func (m DefaultManager) ServiceKeys() ([]*auth.ServiceKey, error) {
	res, err := r.Table(tblNameServiceKeys).Run(m.session)
	if err != nil {
		return nil, err
	}
	keys := []*auth.ServiceKey{}
	if err := res.All(&keys); err != nil {
		return nil, err
	}
	return keys, nil
}

func (m DefaultManager) ServiceKey(key string) (*auth.ServiceKey, error) {
	res, err := r.Table(tblNameServiceKeys).Filter(map[string]string{"key": key}).Run(m.session)
	if err != nil {
		return nil, err

	}
	if res.IsNil() {
		return nil, ErrServiceKeyDoesNotExist
	}
	var k *auth.ServiceKey
	if err := res.One(&k); err != nil {
		return nil, err
	}
	return k, nil
}
func (m DefaultManager) NewServiceKey(description string) (*auth.ServiceKey, error) {
	k, err := m.authenticator.GenerateToken()
	if err != nil {
		return nil, err
	}
	key := &auth.ServiceKey{
		Key:         k[24:],
		Description: description,
	}
	if err := m.SaveServiceKey(key); err != nil {
		return nil, err
	}
	return key, nil
}
func (m DefaultManager) SaveServiceKey(key *auth.ServiceKey) error {
	if _, err := r.Table(tblNameServiceKeys).Insert(key).RunWrite(m.session); err != nil {
		return err
	}

	//m.logEvent("add-service-key", fmt.Sprintf("description=%s", key.Description), []string{"security"})

	return nil
}

func (m DefaultManager) RemoveServiceKey(key string) error {
	if _, err := r.Table(tblNameServiceKeys).Filter(map[string]string{"key": key}).Delete().RunWrite(m.session); err != nil {
		return err
	}

	//m.logEvent("delete-service-key", fmt.Sprintf("key=%s", key), []string{"security"})

	return nil
}


func (m DefaultManager) ChangePassword(username, password string) error {
	if !m.authenticator.IsUpdateSupported() {
		return fmt.Errorf("not supported for authenticator: %s", m.authenticator.Name())
	}

	hash, err := auth.Hash(password)
	if err != nil {
		return err
	}

	if _, err := r.Table(tblNameAccounts).Filter(map[string]string{"username": username}).Update(map[string]string{"password": hash}).Run(m.session); err != nil {
		return err
	}

	//m.logEvent("change-password", username, []string{"security"})

	return nil
}

func (m DefaultManager) Store() *sessions.CookieStore{
	return m.store
}

func (m DefaultManager) StoreKey() string{
	return storeKey
}

func (m DefaultManager) Roles() []*auth.ACL{
	return auth.DefaultACLs()
}

func (m DefaultManager) Role(name string) *auth.ACL{
	//var acl *auth.ACL
	acls := auth.DefaultACLs()
	for _, acl := range acls{
		if acl.RoleName == name{
			return acl
		}
	}

	return nil
}