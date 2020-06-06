package ws

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
	"github.com/tsybulin/gosha/svc"
)

func (wbs *webservice) apiHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/local/", http.StatusMovedPermanently)
	wbs.eventBus.Publish(logger.Topic, logger.LevelInfo, "WebService.handle redirect", r.URL.String())
}

func (wbs *webservice) apiDomains(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(cmp.DomainNames)
	wbs.eventBus.Publish(logger.Topic, logger.LevelInfo, "WebService.handle", r.URL.String())
}

func (wbs *webservice) apiServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	services := struct {
		Services []svc.Description `json:"services,omitempty"`
	}{
		Services: svc.NewRegistry(wbs.eventBus).Services(),
	}
	json.NewEncoder(w).Encode(services)
	wbs.eventBus.Publish(logger.Topic, logger.LevelInfo, "WebService.handle", r.URL.String())
}

func (wbs *webservice) apiComponents(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["service"]
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	srv := svc.NewRegistry(wbs.eventBus).GetService(id)
	components := struct {
		Components []string `json:"components,omitempty"`
	}{
		Components: make([]string, 0),
	}

	if srv != nil {
		components.Components = srv.Components()
	}

	json.NewEncoder(w).Encode(components)
	wbs.eventBus.Publish(logger.Topic, logger.LevelInfo, "WebService.handle", r.URL.String())
}

func (wbs *webservice) apiExecute(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	svc.NewRegistry(nil).Execute(v["service"], v["method"], v["entity"])
	w.Write([]byte("ok"))
	wbs.eventBus.Publish(logger.Topic, logger.LevelInfo, "WebService.handle", r.URL.String())
}

func (wbs *webservice) apiReloadAutomations(w http.ResponseWriter, r *http.Request) {
	wbs.eventBus.Publish(evt.CfgReloadAutomationsTopic, true)
	w.Write([]byte("ok"))
	wbs.eventBus.Publish(logger.Topic, logger.LevelInfo, "WebService.handle", r.URL.String())
}
