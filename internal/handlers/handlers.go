package handlers

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var (
	views = jet.NewSet(
		jet.NewOSFileSystemLoader("./html"),
		jet.InDevelopmentMode(),
	)

	upgradeConnection = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

func Home(w http.ResponseWriter, r *http.Request) {

	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println("ERROR: ", err.Error())
	}

}

type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {

	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ERROR: ", err.Error())
	}

	response := WsJsonResponse{
		Action:      "do it!",
		Message:     `<em><small>Connection to the server</small></em>`,
		MessageType: "html",
	}

	t := time.Tick(5 * time.Second)

	for range t {
		err = ws.WriteJSON(response)
		if err != nil {
			panic(err)
		}
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {

	view, err := views.GetTemplate(tmpl)
	if err != nil {
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		return err
	}

	return nil
}
