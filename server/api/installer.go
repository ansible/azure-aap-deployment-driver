package api

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"server/config"
	"server/controllers"
	"server/engine"
	"server/handler"
	"server/model"
	"server/persistence"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Installer struct {
	router     *mux.Router
	db         *gorm.DB
	engine     *engine.Engine
	httpServer *http.Server
	loginMgr   handler.LoginManager
}

// func NewApp(database *persistence.Database) *Installer {
func NewApp(database *persistence.Database, engine *engine.Engine, loginManager handler.LoginManager) *Installer {
	app := &Installer{
		db:       database.Instance,
		engine:   engine,
		loginMgr: loginManager,
	}
	app.initialize()
	return app
}

func (a *Installer) initialize() {

	a.router = mux.NewRouter()
	a.configureSessionHelper()
	a.setRouters()
}

func (a *Installer) GetRouter() *mux.Router {
	return a.router
}

func (a *Installer) configureSessionHelper() {
	sessionConfig := model.SessionConfig{}
	a.db.Find(&sessionConfig)

	// session auth key is not setup yet, generate one and store it
	if sessionConfig.SessionAuthKey == nil {
		b, err := handler.GenerateSessionAuthKey()
		if err != nil {
			log.Fatalf("Could not generate session key. %v", err)
		}
		sessionConfig.SessionAuthKey = b
		a.db.Save(&sessionConfig)
	}

	handler.ConfigureSessionHelper(handler.SessionHelperConfiguration{
		AuthKey:      sessionConfig.SessionAuthKey,
		CookieDomain: config.GetEnvironment().SESSION_COOKIE_DOMAIN,
		CookieName:   config.GetEnvironment().SESSION_COOKIE_NAME,
		CookiePath:   config.GetEnvironment().SESSION_COOKIE_PATH,
		Secure:       config.GetEnvironment().SESSION_COOKIE_SECURE,
		MaxAge:       config.GetEnvironment().SESSION_COOKIE_MAX_AGE,
	})
}

func (a *Installer) setRouters() {
	a.Get("/status", a.WrapHandlerWithDB(handler.Status))
	a.Get("/authtype", handler.AuthType)
	a.Post("/login", a.WrapHandlerWithDB(handler.CredentialsHandler{}.GetLoginHandler())) // GetLoginHandler() must be executed to do its initialization
	a.Post("/logout", handler.EnsureAuthenticated(handler.Logout))
	a.Get("/step", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.GetAllSteps)))
	a.Get("/step/{id}", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.GetStep)))
	a.Get("/execution", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.GetAllExecutions)))
	a.Get("/execution/{id}", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.GetExecution)))
	a.Post("/execution/{id}/restart", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.Restart)))
	a.Post("/cancelAllSteps", handler.EnsureAuthenticated(a.WrapHandlerWithDBAndEngine(handler.CancelAllSteps)))
	a.Post("/deleteContainer", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.DeleteContainer)))
	a.Post("/terminate", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.Terminate)))
	a.Get("/azmarketplaceentitlementscount", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.GetNumOfAzureMarketplaceEntitlements)))
	if config.IsSsoEnabled() {
		ssoClient, ok := a.loginMgr.(*handler.SsoHandler)
		if ok {
			log.Trace("Configuring SSO API endpoints.")
			// Could have failed, in which case loginMgr is a credentials login handler and these aren't needed
			a.Get(handler.REDIRECT_PATH, a.WrapHandlerWithDB(ssoClient.SsoRedirect))
			// GET handling for /login due to nginx forwarding here
			a.Get("/login", a.WrapHandlerWithDB(a.loginMgr.GetLoginHandler()))
		}
	}
}

func (a *Installer) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, f).Methods("GET")
}

func (a *Installer) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, f).Methods("POST")
}

// WrapHandlerWithDB returns an HTTP HandlerFunc for invoking handlers that need DB as first argument
func (a *Installer) WrapHandlerWithDB(fn handler.HandleFuncWithDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(a.db, w, r)
	}
}

func (a *Installer) WrapHandlerWithDBAndEngine(fn handler.HandleFuncWithDBAndEngine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(a.db, a.engine, w, r)
	}
}

func (a *Installer) Run() {
	log.Println("Starting httpServer...")
	a.httpServer = &http.Server{}
	a.httpServer.Addr = fmt.Sprintf("%s:%s", config.Args.Host, config.Args.Port)
	a.httpServer.Handler = a.router
	err := controllers.AddCancelHandler("API Server", a.stopServer)
	if err != nil {
		log.Fatalf("Error while adding cancel handler to API server: %v", err)
	}
	runServer(a)
}

func runServer(a *Installer) {
	if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Errorf("HTTP server shut down: %v", err)
	}
}

func (a *Installer) stopServer() {
	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		log.Errorf("HTTP server error while shutting down: %v", err)
	}
}
