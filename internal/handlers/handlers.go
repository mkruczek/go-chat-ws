package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

const (
	actionUserList   = "user_list"
	actionPublishMeg = "publish_message"
)

var (
	wsChan  = make(chan WsPayload)
	clients = make(map[WebSocketConnection]string)

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

type WebSocketConnection struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"messageType"`
	ConnectedUsers []string `json:"connectedUsers"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	UserName string              `json:"userName"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
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

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		panic(err)
	}

	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(fmt.Sprintf("recovery... : %v", r))
		}
	}()

	var body WsPayload

	for {
		err := conn.ReadJSON(&body)
		if err != nil {
			log.Println(fmt.Sprintf("error for reading body from WS : %s", err.Error()))
		} else {
			body.Conn = *conn
			wsChan <- body
		}
	}
}

func ListenForWsChannel() {

	var response WsJsonResponse

	for {
		e := <-wsChan

		switch e.Action {
		case "username":
			clients[e.Conn] = e.UserName
			response.ConnectedUsers = getUserList()
			response.Action = actionUserList
			broadcastToAll(response)
		case "left":
			delete(clients, e.Conn)
			response.ConnectedUsers = getUserList()
			response.Action = actionUserList
			broadcastToAll(response)
		case "message":
			response.Action = actionPublishMeg
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.UserName, e.Message)
			broadcastToAll(response)
		default:
			log.Println(fmt.Sprintf("error, unknown action: %s", e.Action))
		}
	}
}

func getUserList() []string {
	var result []string

	for _, v := range clients {
		if v != "" {
			result = append(result, v)
		}
	}

	sort.Strings(result)

	return result
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println(fmt.Sprintf("error during sending WsJsonResponse: %s", err.Error()))
			_ = client.Close()
			delete(clients, client)
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
