package controllers

import (
	"fmt"

	"github.com/mistrale/jsonception/app/socket"

	"github.com/revel/revel"
	"golang.org/x/net/websocket"
)

// WebSocket controller
type WebSocket struct {
	*revel.Controller
}

// RoomSocket : WebSocket listen to room
func (c WebSocket) ListenScriptRun(room_name string, ws *websocket.Conn) revel.Result {
	fmt.Printf("room name : %s\n", room_name)
	room, ok := socket.Rooms[room_name]
	if !ok {
		if err := websocket.JSON.Send(ws, "Connection to websocket invalide : room not found"); err != nil {
			fmt.Printf("err : %s\n", err.Error())
			// They disconnected
			return nil
		}
	}
	i := 0
	for {
		//fmt.Println("Boucle pour listen")
		room.Mux.Lock()
		for ; i < len(room.Tmp); i++ {
			if err := websocket.JSON.Send(ws, room.Tmp[i]); err != nil {
				fmt.Printf("err : %s\n", err.Error())
			}
		}
		if room.IsClosed == true {
			fmt.Println("Boucle pour listeN CA FA FINI")
			return nil
		}
		room.Mux.Unlock()
	}
	return nil
}
