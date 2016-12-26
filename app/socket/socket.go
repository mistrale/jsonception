package socket

var (
	// archive initially, and then new messages as they come in.
	Rooms map[string]*Room = make(map[string]*Room)
)

type Room struct {
	Name string
	Chan chan string
}

func CreateRoom(name string) *Room {
	Rooms[name] = &Room{Name: name, Chan: make(chan string)}
	return Rooms[name]
}

func DeleteRoom(name string) {
	delete(Rooms, name)
}
