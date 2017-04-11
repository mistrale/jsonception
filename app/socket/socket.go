package socket

import (
	"fmt"
	"sync"

	"github.com/mistrale/jsonception/app/dispatcher"
)

var (
	// archive initially, and then new messages as they come in.
	Rooms map[string]*Room = make(map[string]*Room)
)

type Room struct {
	Name     string
	Tmp      []dispatcher.Event
	Chan     chan dispatcher.Event
	Mux      sync.Mutex
	IsClosed bool
}

func CreateRoom(name string) *Room {
	fmt.Printf("creatiing room : %s\n", name)
	room := &Room{Name: name, Chan: make(chan dispatcher.Event), Tmp: make([]dispatcher.Event, 0), IsClosed: false}
	Rooms[name] = room
	fmt.Println("on demarrrre la boucle")

	go func() {
		for {
			msg := <-room.Chan
			if msg.Body == "end_"+room.Name {
				room.Mux.Lock()
				fmt.Printf("[room][%s] : Closing room\n", room.Name)
				room.IsClosed = true
				room.Mux.Unlock()
				break
			}
			//fmt.Printf("[room][%s] : on recoit %s\n", name, msg)
			room.Mux.Lock()
			room.Tmp = append(room.Tmp, msg)
			room.Mux.Unlock()

		}
	}()
	return room
}

func DeleteRoom(name string) {
	delete(Rooms, name)
}
