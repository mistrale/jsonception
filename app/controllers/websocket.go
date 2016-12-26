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
func (c WebSocket) RoomSocket(room_name string, ws *websocket.Conn) revel.Result {
	for {
		room := socket.Rooms[room_name]
		msg := <-room.Chan
		if err := websocket.JSON.Send(ws, msg); err != nil {
			fmt.Printf("err : %s\n", err.Error())
			// They disconnected
			return nil
		}
		if msg == "end_"+room_name {
			break
		}
	}
	return nil
}
