package socket

import "fmt"

var (
	// archive initially, and then new messages as they come in.
	Rooms map[string]*Room = make(map[string]*Room)
)

type Room struct {
	Name string
	Chan chan map[string]interface{}
}

func CreateRoom(name string) *Room {
	fmt.Printf("creatiing room : %s\n", name)
	Rooms[name] = &Room{Name: name, Chan: make(chan map[string]interface{})}
	return Rooms[name]
}

func DeleteRoom(name string) {
	delete(Rooms, name)
}
