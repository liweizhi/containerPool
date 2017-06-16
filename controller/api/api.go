package api

import(
	"net/http"

	"github.com/liweizhi/containerPool/controller/manager"
	"github.com/gorilla/mux"


	log"github.com/Sirupsen/logrus"
	"github.com/liweizhi/containerPool/auth"
	"github.com/urfave/negroni"
	"github.com/vulcand/oxy/forward"
	"github.com/liweizhi/containerPool/controller/middleware"
	"io"
)

type API struct{
	listenAddr	string
	manager		manager.Manager
	cors		bool
	dockerUrl	string
	fwd		*forward.Forwarder
}

type Credentials struct{
	Username	string
	Password	string
}
func NewAPI(listenAddr string, manager manager.Manager, cors bool, dockerUrl string) *API{
	return &API{
		listenAddr: listenAddr,
		manager: manager,
		cors: cors,
		dockerUrl: dockerUrl,

	}
}
func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world")
}
func(a *API) Start(){
	mainMux := http.NewServeMux()

	a.fwd, _ = forward.New()
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
	apiRouter.HandleFunc("/api/roles", a.roles).Methods("GET")
	apiRouter.HandleFunc("/api/roles/{name}", a.role).Methods("GET")
	apiRouter.HandleFunc("/api/node", a.nodeInfo).Methods("GET")
	mainMux.Handle("/", http.FileServer(http.Dir("myStatic")))
	//mainMux.HandleFunc("/", hello)
	authRequired := middleware.NewAuthRequired(a.manager)
	accessRequired := middleware.NewAccessRequired(a.manager)

	middlewareStack := negroni.New()
	middlewareStack.Use(negroni.HandlerFunc(authRequired.HandlerFuncWithNext))
	middlewareStack.Use(negroni.HandlerFunc(accessRequired.HandlerFuncWithNext))
	middlewareStack.UseHandler(apiRouter)
	mainMux.Handle("/api/", middlewareStack)


	accountRouter := mux.NewRouter()
	accountRouter.HandleFunc("/account/changepassword", a.changePassword).Methods("POST")
	accountMiddlewareStack := negroni.New()
	accountMiddlewareStack.Use(negroni.HandlerFunc(authRequired.HandlerFuncWithNext))
	accountMiddlewareStack.UseHandler(accountRouter)
	mainMux.Handle("/account/", accountMiddlewareStack)


	loginRouter := mux.NewRouter()
	loginRouter.HandleFunc("/auth/login", a.login).Methods("POST")
	mainMux.Handle("/auth/", loginRouter)

	dockerRouter := mux.NewRouter()

	urlMap := map[string]map[string]http.HandlerFunc{
		"GET": {
			"/containers/json":			a.dockerHandler,  //list containers
			"/containers/{id}/json":		a.dockerHandler, //Inspect a container
			"/containers/{id}/top":			a.dockerHandler, //List processes running inside a container
			"/containers/{id}/logs":		a.dockerHandler, //Get container logs
			"/containers/{id}/stats":		a.dockerHandler, //Get container stats based on resource usage
			"/images/json":				a.dockerHandler,
			"/images/{name}/json":			a.dockerHandler,
			"/networks":				a.dockerHandler, //List networks
			"/networks/{id}":			a.dockerHandler,
			"/swarm":				a.dockerHandler,
			"/nodes":				a.dockerHandler,
			"/nodes/{id}":				a.dockerHandler,
			"/info":				a.dockerHandler,
			"/events":				a.dockerHandler,
			"/version":				a.dockerHandler,
			"/_ping":				a.dockerHandler,

		},

		"POST": {
			"/containers/create":			a.dockerHandler, //Create a container
			"/containers/{id}/start":		a.dockerHandler, //Start a container
			"/containers/{id}/stop":		a.dockerHandler, //Stop a container
			"/containers/{id}/restart":		a.dockerHandler, //Restart a container
			"/containers/{id}/kill":		a.dockerHandler, //Kill a container
			"/containers/{id}/update":		a.dockerHandler,
			"/containers/{id}/rename":		a.dockerHandler,
			"/containers/{id}/pause":		a.dockerHandler,
			"/build":				a.dockerHandler, //Build an image
			"/images/create":			a.dockerHandler, //Create an image by either pulling it from a registry or importing it.
			"/images/{name}/push":			a.dockerHandler,
			"/commit":				a.dockerHandler, //Create a new image from a container
			"/images/load":				a.dockerHandler, //Load a set of images and tags into a repository.
			"/swarm/init":				a.dockerHandler,
			"/swarm/join":				a.dockerHandler,
			"/auth":				a.dockerHandler,



		},

		"DELETE": {
			"/containers/{id}":			a.dockerHandler,
			"/images/{name}":			a.dockerHandler,
			"/nodes/{id}":				a.dockerHandler,
		},

	}

	for method, routes := range urlMap{
		for route, handler := range routes{
			dockerRouter.HandleFunc(route, handler).Methods(method)
		}
	}

	dockerMiddlewareStack := negroni.New()
	dockerMiddlewareStack.Use(negroni.HandlerFunc(authRequired.HandlerFuncWithNext))
	dockerMiddlewareStack.Use(negroni.HandlerFunc(accessRequired.HandlerFuncWithNext))
	dockerMiddlewareStack.UseHandler(dockerRouter)
	mainMux.Handle("/containers/", dockerMiddlewareStack)
	mainMux.Handle("/images/", dockerMiddlewareStack)
	mainMux.Handle("/networks", dockerMiddlewareStack)
	mainMux.Handle("/networks/", dockerMiddlewareStack)
	mainMux.Handle("/nodes", dockerMiddlewareStack)
	mainMux.Handle("/build", dockerMiddlewareStack)
	mainMux.Handle("/_ping", dockerMiddlewareStack)
	mainMux.Handle("/version", dockerMiddlewareStack)
	mainMux.Handle("/info", dockerMiddlewareStack)
	mainMux.Handle("/auth", dockerMiddlewareStack)










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




