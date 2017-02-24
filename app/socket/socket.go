package socket

import (
	"fmt"
	"sync"
)

var (
	// archive initially, and then new messages as they come in.
	Rooms map[string]*Room = make(map[string]*Room)
)

type Room struct {
	Name     string
	Tmp      []map[string]interface{}
	Chan     chan map[string]interface{}
	Mux      sync.Mutex
	IsClosed bool
}

func CreateRoom(name string) *Room {
	fmt.Printf("creatiing room : %s\n", name)
	room := &Room{Name: name, Chan: make(chan map[string]interface{}), Tmp: make([]map[string]interface{}, 0), IsClosed: false}
	Rooms[name] = room
	fmt.Println("on demarrrre la boucle")

	go func() {
		for {
			msg := <-room.Chan
			//fmt.Printf("[room][%s] : on recoit %s\n", name, msg)
			room.Mux.Lock()
			room.Tmp = append(room.Tmp, msg)
			room.Mux.Unlock()
			if msg["response"] == "end_"+room.Name {
				room.Mux.Lock()
				fmt.Printf("[room][%s] : Closing room\n", room.Name)
				room.IsClosed = true
				room.Mux.Unlock()
				break
			}
		}
	}()
	return room
}

func DeleteRoom(name string) {
	delete(Rooms, name)
}
