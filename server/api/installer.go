package api

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"server/config"
	"server/controllers"
	"server/handler"
	"server/persistence"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Installer struct {
	router     *mux.Router
	db         *gorm.DB
	httpServer *http.Server
}

func NewApp(database *persistence.Database) *Installer {
	app := &Installer{
		db: database.Instance,
	}
	app.initialize()
	return app
}

func (a *Installer) initialize() {

	a.router = mux.NewRouter()
	a.setRouters()
}

func (a *Installer) setRouters() {
	a.Get("/status", a.Status)
	a.Get("/step", handler.BasicAuth(a.GetAllSteps))
	a.Get("/step/{id}", handler.BasicAuth(a.GetStep))
	a.Get("/execution", handler.BasicAuth(a.GetAllExecutions))
	a.Get("/execution/{id}", handler.BasicAuth(a.GetExecution))
	a.Post("/execution/{id}/restart", handler.BasicAuth(a.Restart))
	a.Post("/deleteContainer", handler.BasicAuth(a.DeleteContainer))
	a.Post("/terminate", handler.BasicAuth(a.Terminate))
}

func (a *Installer) Status(w http.ResponseWriter, r *http.Request) {
	handler.RespondOk(w)
}

func (a *Installer) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, f).Methods("GET")
}

func (a *Installer) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, f).Methods("POST")
}

func (a *Installer) GetAllSteps(w http.ResponseWriter, r *http.Request) {
	handler.GetAllSteps(a.db, w, r)
}

func (a *Installer) GetStep(w http.ResponseWriter, r *http.Request) {
	handler.GetStep(a.db, w, r)
}

func (a *Installer) GetAllExecutions(w http.ResponseWriter, r *http.Request) {
	handler.GetAllExecutions(a.db, w, r)
}

func (a *Installer) GetExecution(w http.ResponseWriter, r *http.Request) {
	handler.GetExecution(a.db, w, r)
}

func (a *Installer) Restart(w http.ResponseWriter, r *http.Request) {
	handler.Restart(a.db, w, r)
}

func (a *Installer) DeleteContainer(w http.ResponseWriter, r *http.Request) {
	handler.DeleteContainer(a.db, w, r)
}

func (a *Installer) Terminate(w http.ResponseWriter, r *http.Request) {
	handler.Terminate(a.db, w, r)
}

func (a *Installer) Run() {
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
