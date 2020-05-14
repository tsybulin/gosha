package ws

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// LocalDir ...
const LocalDir = "/local/"

// WebService ...
type WebService interface {
	Start()
}

type webservice struct {
	eventBus   evt.Bus
	listenAddr string
	authKey    string
	sessionID  int32
}

// NewWebService ...
func NewWebService(eventBus evt.Bus, la, authKey string) WebService {
	return &webservice{
		eventBus:   eventBus,
		listenAddr: la,
		authKey:    authKey,
		sessionID:  0,
	}
}

func (wbs *webservice) Start() {
	go func() {
		wbs.start()
	}()
}

func (wbs *webservice) wsHandler(w http.ResponseWriter, r *http.Request) {
	wbs.sessionID++
	sessionID := wbs.sessionID
	wbs.eventBus.Publish(logger.Topic, logger.LevelWarn, "WebService.wsHandler new", sessionID)
	wss := newWsocketService(wbs.eventBus, wbs.authKey, sessionID)
	wss.wsHandler(w, r)
	wbs.eventBus.Publish(logger.Topic, logger.LevelWarn, "WebService.wsHandler done", sessionID)
}

func (wbs *webservice) start() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", wbs.apiHome)

	router.PathPrefix(LocalDir).Handler(http.StripPrefix(LocalDir, http.FileServer(http.Dir("./config/www/"))))

	router.HandleFunc("/api/domains", wbs.apiDomains).Methods("GET")

	router.HandleFunc("/api/services", wbs.apiServices).Methods("GET")

	router.HandleFunc("/api/service/{service}/components", wbs.apiComponents).Methods("GET")

	router.HandleFunc("/api/service/{service}/execute/{method}/{entity}", wbs.apiExecute).Methods("POST")

	router.HandleFunc("/api/ws", wbs.wsHandler)

	http.ListenAndServe(wbs.listenAddr, router)
}
