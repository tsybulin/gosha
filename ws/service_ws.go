package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
	"github.com/tsybulin/gosha/svc"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	livePeriod = 5 * time.Second
)

type syncbool struct {
	mux sync.Mutex
	v   bool
}

func (sb *syncbool) get() bool {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	return sb.v
}

func (sb *syncbool) set(b bool) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	sb.v = b
}

type wsocketService struct {
	eventBus  evt.Bus
	authKey   string
	auth      bool
	upgrader  websocket.Upgrader
	statechan chan evt.Message
	sb        *syncbool
	mux       sync.Mutex
	sessionID int32
}

func (wss *wsocketService) writeJSON(ws *websocket.Conn, v interface{}) error {
	wss.mux.Lock()
	defer wss.mux.Unlock()
	return ws.WriteJSON(v)
}

func (wss *wsocketService) writeMessage(ws *websocket.Conn, messageType int, data []byte) error {
	wss.mux.Lock()
	defer wss.mux.Unlock()
	return ws.WriteMessage(messageType, data)
}

func (wss *wsocketService) stateChageHandler(event evt.Message) {
	if wss.sb.get() {
		wss.statechan <- event
	} else {
		wss.eventBus.Publish(logger.Topic, logger.LevelWarn, "WSocketService.eventBus statechan closed ", wss.sessionID)
	}
}

func (wss *wsocketService) reader(ws *websocket.Conn) {
	defer func() {
		wss.eventBus.Unsubscribe(evt.StateChangedTopic, fmt.Sprintf("WsocketService.stateChangeHandler-%v", wss.sessionID))

		wss.eventBus.Publish(logger.Topic, logger.LevelDebug, "WSocketService.reader close statechan", wss.sessionID)
		wss.sb.set(false)
		close(wss.statechan)

		ws.Close()
		wss.eventBus.Publish(logger.Topic, logger.LevelInfo, "WSocketService.reader Close", wss.sessionID)
	}()

	ws.SetReadLimit(8192)
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		mt, bs, err := ws.ReadMessage()
		if err != nil {
			wss.eventBus.Publish(logger.Topic, logger.LevelDebug, "WSocketService.reader read error", wss.sessionID, err.Error())
			break
		}

		if mt != websocket.TextMessage {
			continue
		}

		var message evt.Message
		if json.Unmarshal(bs, &message) == nil {

			if message.Type == "auth" {
				am := evt.Message{
					Type: "auth_failed",
				}

				if message.Token == wss.authKey {
					wss.auth = true
					am.Type = "auth_ok"
					am.ID = message.ID
				}

				ws.SetWriteDeadline(time.Now().Add(writeWait))
				wss.writeJSON(ws, &am)

				continue
			}

			if message.Type == "ping" {
				pm := evt.Message{
					ID:   message.ID,
					Type: "pong",
				}
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := wss.writeJSON(ws, &pm); err != nil {
					wss.eventBus.Publish(logger.Topic, logger.LevelDebug, "WSocketService.reader write error", wss.sessionID, err.Error())
					return
				}

				continue
			}

			if !wss.auth {
				am := evt.Message{
					Type: "auth_required",
				}

				ws.SetWriteDeadline(time.Now().Add(writeWait))
				wss.writeJSON(ws, &am)

				continue
			}

			if message.Type == "get_states" {
				states := svc.NewRegistry(wss.eventBus).States()
				states.ID = message.ID

				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := wss.writeJSON(ws, &states); err != nil {
					wss.eventBus.Publish(logger.Topic, logger.LevelDebug, "WSocketService.reader write error", wss.sessionID, err.Error())
					return
				}
			}
		}

		var command Command
		if json.Unmarshal(bs, &command) != nil {
			continue
		}

		if command.Type == "call_service" {
			svc.NewRegistry(wss.eventBus).Execute("service."+command.Domain, command.Command, command.ServiceData.EntityID)

			rm := Result{
				ID:      message.ID,
				Type:    "result",
				Success: "true",
			}

			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := wss.writeJSON(ws, &rm); err != nil {
				wss.eventBus.Publish(logger.Topic, logger.LevelDebug, "WSocketService.reader write error", wss.sessionID, err.Error())
			}
		}
	}
}

func (wss *wsocketService) writer(ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		pingTicker.Stop()
		ws.Close()
		wss.eventBus.Publish(logger.Topic, logger.LevelInfo, "WSocketService.writer Close", wss.sessionID)
	}()

	var wsmid int = 0

	for {
		select {
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := wss.writeMessage(ws, websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case sm := <-wss.statechan:
			if !wss.auth {
				continue
			}

			wsmid++
			if wsmid > 1000 {
				wsmid = 1
			}

			sm.ID = wsmid

			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := wss.writeJSON(ws, &sm); err != nil {
				wss.eventBus.Publish(logger.Topic, logger.LevelDebug, "WSocketService.writer write error", wss.sessionID, err.Error())
				return
			}
		}
	}
}

func (wss *wsocketService) wsHandler(w http.ResponseWriter, r *http.Request) {
	wss.eventBus.Publish(logger.Topic, logger.LevelInfo, "WSocketService.wsHandler start ", wss.sessionID)
	defer wss.eventBus.Publish(logger.Topic, logger.LevelInfo, "WSocketService.wsHandler stop ", wss.sessionID)

	ws, err := wss.upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			wss.eventBus.Publish(logger.Topic, logger.LevelWarn, "WSocketService.wsHandler error", wss.sessionID, err.Error())
		}
		return
	}

	wss.eventBus.Subscribe(evt.StateChangedTopic, fmt.Sprintf("WsocketService.stateChangeHandler-%v", wss.sessionID), wss.stateChageHandler)

	go wss.writer(ws)
	wss.reader(ws)
}

func newWsocketService(eventBus evt.Bus, authKey string, sessionID int32) *wsocketService {
	wss := &wsocketService{
		eventBus:  eventBus,
		authKey:   authKey,
		auth:      false,
		statechan: make(chan evt.Message),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		sb:        &syncbool{v: true},
		sessionID: sessionID,
	}

	return wss
}
