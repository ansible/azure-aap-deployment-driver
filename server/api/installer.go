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
}

// func NewApp(database *persistence.Database) *Installer {
func NewApp(database *persistence.Database, engine *engine.Engine) *Installer {
	app := &Installer{
		db:     database.Instance,
		engine: engine,
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
	a.Post("/login", handler.GetLoginHandler("admin", config.GetEnvironment().PASSWORD))
	a.Post("/logout", handler.EnsureAuthenticated(handler.Logout))
	a.Get("/step", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.GetAllSteps)))
	a.Get("/step/{id}", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.GetStep)))
	a.Get("/execution", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.GetAllExecutions)))
	a.Get("/execution/{id}", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.GetExecution)))
	a.Post("/execution/{id}/restart", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.Restart)))
	a.Post("/cancelAllSteps", handler.EnsureAuthenticated(a.WrapHandlerWithDBAndEngine(handler.CancelAllSteps)))
	a.Post("/deleteContainer", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.DeleteContainer)))
	a.Post("/terminate", handler.EnsureAuthenticated(a.WrapHandlerWithDB(handler.Terminate)))
	a.Post("/eventhook", a.WrapHandlerWithDBAndEngine(handler.EventHook))
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
	controllers.AddCancelHandler("API Server", a.stopServer)
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
