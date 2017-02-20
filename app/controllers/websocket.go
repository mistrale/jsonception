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
func (c WebSocket) ListenExecutionRun(room_name string, ws *websocket.Conn) revel.Result {
	fmt.Printf("room name : %s\n", room_name)
	room, ok := socket.Rooms[room_name]
	if !ok {
		if err := websocket.JSON.Send(ws, "Connection to websocket invalide : room not found"); err != nil {
			fmt.Printf("err : %s\n", err.Error())
			// They disconnected
			return nil
		}
	}
	for {

		response := <-room.Chan
		//fmt.Printf("SENDING DATA :%s\n", response)
		if err := websocket.JSON.Send(ws, response); err != nil {
			go func() {
				for {
					msg := <-room.Chan
					if msg["response"] == "end_"+room_name {
						break
					}
				}
			}()
			fmt.Printf("err : %s\n", err.Error())
			// They disconnected
			return nil
		}
		if response["response"] == "end_"+room_name {
			break
		}
	}
	return nil
}

//
// for {
// 	fmt.Println("Boucle pour listen")
// 	room.Mux.Lock()
// 	if len(room.Tmp) > 0 {
// 		response := room.Tmp[0]
// 		room.Tmp = append(room.Tmp[:0], room.Tmp[1:]...)
// 		room.Mux.Unlock()
//
// 		fmt.Printf("on  a send en socket : %s avec size tmp : %d\n", response, len(room.Tmp))
// 		// msg := <-room.Chan
// 		// if msg["response"] == "end_"+room_name {
// 		// 	break
// 		// }
// 		if err := websocket.JSON.Send(ws, response); err != nil {
// 			fmt.Printf("err : %s\n", err.Error())
// 			return nil
// 		}
// 		if room.IsClosed == true {
// 			fmt.Println("Boucle pour listeN CA FA FINI")
//
// 			break
// 		}
// 	}
// 	fmt.Println("Boucle pour listen - END")
//
// }
