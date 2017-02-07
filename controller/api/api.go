package api

import(
	"net/http"

	"github.com/liweizhi/containerPool/controller/manager"
	"github.com/gorilla/mux"


	log"github.com/Sirupsen/logrus"
	"github.com/liweizhi/containerPool/auth"
	"github.com/urfave/negroni"

	"github.com/liweizhi/containerPool/controller/middleware"
)

type API struct{
	listenAddr	string
	manager		manager.Manager
	cors		bool
	dockerUrl	string
}

type Credentials struct{
	Username	string
	Password	string
}
func NewAPI(listenAddr string, manager manager.Manager, cors bool, dockerUrl string) *API{
	return &API{listenAddr, manager, cors, dockerUrl}
}
func(a *API) Start(){
	mainMux := http.NewServeMux()

	//testRouter := mux.NewRouter()

	//testRouter.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request){
	//	io.WriteString(w, "hello world")
	//})
	//mainMux.Handle("/api/", testRouter)

	apiRouter := mux.NewRouter()
	apiRouter.HandleFunc("/api/accounts", a.saveAccount).Methods("POST")
	apiRouter.HandleFunc("/api/accounts", a.accounts).Methods("GET")
	apiRouter.HandleFunc("/api/accounts/{username}", a.account).Methods("GET")
	apiRouter.HandleFunc("/api/accounts/{username}", a.deleteAccount).Methods("DELETE")

	middlewareStack := negroni.New()
	authRequired := middleware.NewAuthRequired(a.manager)
	accessRequired := middleware.NewAccessRequired(a.manager)

	middlewareStack.Use(negroni.HandlerFunc(authRequired.HandlerFuncWithNext))
	middlewareStack.Use(negroni.HandlerFunc(accessRequired.HandlerFuncWithNext))
	middlewareStack.UseHandler(apiRouter)
	mainMux.Handle("/api/", apiRouter)


	accountRouter := mux.NewRouter()
	accountRouter.HandleFunc("/account/changepassword", a.changePassword).Methods("POST")
	accountMiddlewareStack := negroni.New()
	accountMiddlewareStack.Use(negroni.HandlerFunc(authRequired.HandlerFuncWithNext))
	accountMiddlewareStack.UseHandler(accountRouter)
	mainMux.Handle("/account/", accountRouter)


	loginRouter := mux.NewRouter()
	loginRouter.HandleFunc("/auth/login", a.login).Methods("POST")
	mainMux.Handle("/auth/", loginRouter)









	if _, err := a.manager.Account("admin"); err == manager.ErrAccountDoesNotExist{
		account := &auth.Account{
			Username: "admin",
			Password: "thucloud",
			Roles: []string{"admin"},
		}

		if err := a.manager.SaveAccount(account); err !=nil{
			log.Fatalln(err)
		}
	}












	http.ListenAndServe(a.listenAddr, mainMux)




}




