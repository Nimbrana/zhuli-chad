// Package signaling establish the connection between clients of WebRTC protocol
// This project is based on signaling-go project from sfRTC-framework of Karina Romero (@KarinaRomero)
// The original github repository is https://github.com/KarinaRomero/signaling-go
package signaling

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type clientWs struct {
	ws   *websocket.Conn
	Name string
}

var clients = make(map[*websocket.Conn]bool)
var clientsWs = []clientWs{}

var itemsMap map[string]interface{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

// HandleConnections manage the incomming connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	for {
		var message interface{}

		messageType, p, err := ws.ReadMessage()

		if err != nil {
			log.Printf("error: %v", err)
			fmt.Println(messageType)
			delete(clients, ws)
			break
		}

		json.Unmarshal(p, &message)

		itemsMap = message.(map[string]interface{})
		//fmt.Println(itemsMap)

		switch itemsMap["type"] {

		case "login":
			if err := login(ws); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		case "offer":
			if err := processOffer(ws); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		case "answer":

			if err := processAnswer(ws); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		case "candidate":
			if err := processCandidate(ws); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		case "leave":
			if err := leave(ws); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		default:
			if err := unknownCommand(ws); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func login(ws *websocket.Conn) error {

	var js []byte
	var err error

	var lm = loginMessage{}

	if findClient(itemsMap["name"].(string), nil) == -1 {
		var client = clientWs{
			ws:   ws,
			Name: itemsMap["name"].(string),
		}
		clientsWs = append(clientsWs, client)

		lm.Type = itemsMap["type"].(string)
		lm.Succes = true

		js, err = json.Marshal(lm)
	} else {
		lm.Type = itemsMap["type"].(string)
		lm.Succes = false
		js, err = json.Marshal(lm)
	}

	if err != nil {
		return err
	}

	ws.WriteMessage(1, js)

	return nil
}

func processOffer(ws *websocket.Conn) error {
	fmt.Println("Send offer to", itemsMap["name"].(string))

	var Otherindex = findClient(itemsMap["name"].(string), nil)
	var myIndex = findClient("", ws)

	if myIndex != -1 && Otherindex != -1 {

		var off map[string]interface{}
		off = itemsMap["offer"].(map[string]interface{})

		var offerStructure = offer{
			Type: "offer",
			SDP:  off["sdp"].(string),
		}

		var om = offerMessage{
			"offer",
			offerStructure,
			clientsWs[myIndex].Name,
		}

		js, err := json.Marshal(om)

		if err != nil {
			return err
		}
		clientsWs[Otherindex].ws.WriteMessage(1, js)
	}

	return nil
}

func processAnswer(ws *websocket.Conn) error {
	fmt.Println("answer to: ", itemsMap["name"].(string))

	var Otherindex = findClient(itemsMap["name"].(string), nil)

	if Otherindex != -1 {

		var aws map[string]interface{}
		aws = itemsMap["answer"].(map[string]interface{})

		//fmt.Println(aws["sdp"].(string))

		var answerStructure = answer{
			Type: "answer",
			SDP:  aws["sdp"].(string),
		}

		var am = answerMessage{
			"answer",
			answerStructure,
		}

		js, err := json.Marshal(am)

		if err != nil {
			return err
		}
		clientsWs[Otherindex].ws.WriteMessage(1, js)
	}

	return nil
}

func processCandidate(ws *websocket.Conn) error {
	fmt.Println("Sending candidate to ", itemsMap["name"].(string))

	var Otherindex = findClient(itemsMap["name"].(string), nil)

	if Otherindex != -1 {
		var cnddt map[string]interface{}
		cnddt = itemsMap["candidate"].(map[string]interface{})

		var candidateToSend = candidate{
			Candidate:     cnddt["candidate"].(string),
			SdpMid:        cnddt["sdpMid"].(string),
			SdpMLineIndex: cnddt["sdpMLineIndex"].(float64),
		}

		var cm = candidateMessage{
			"candidate",
			candidateToSend,
		}

		js, err := json.Marshal(cm)

		if err != nil {
			return err
		}
		clientsWs[Otherindex].ws.WriteMessage(1, js)
	}

	return nil
}

func leave(ws *websocket.Conn) error {
	//fmt.Println("leave: ", itemsMap["name"].(string))

	var Otherindex = findClient(itemsMap["name"].(string), nil)

	if Otherindex != -1 {

		lm := leaveMessage{"leave"}

		js, err := json.Marshal(lm)
		if err != nil {
			return err
		}
		clientsWs[Otherindex].ws.WriteMessage(1, js)
	}

	return nil
}

func unknownCommand(ws *websocket.Conn) error {
	//fmt.Println("default")

	um := unrecognizedMessage{
		"error",
		"Unrecognized command: " + itemsMap["candidate"].(string),
	}

	js, err := json.Marshal(um)
	if err != nil {
		return err
	}
	ws.WriteMessage(1, js)

	return nil
}

func findClient(c string, ws *websocket.Conn) int {
	if c == "" {
		for i := range clientsWs {
			if clientsWs[i].ws == ws {
				//found
				return i
			}
		}
	} else {
		for i := range clientsWs {
			if clientsWs[i].Name == c {
				//found
				return i
			}
		}
	}
	//not found
	return -1
}
