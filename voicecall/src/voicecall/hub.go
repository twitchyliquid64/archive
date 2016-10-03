package voicecall


var Hubs map[string]*Hub = make(map[string]*Hub)

type Hub struct {
	conns map[int]*Connection
	Name string
}

func (inst *Hub)attach(conn *Connection){
	inst.conns[conn.getID()] = conn
}

func (inst *Hub)Broadcast(recieverID int, data []byte){
	for key, val := range inst.conns {
		if key != recieverID {
			val.Send(data)
		}
	}
}

func Attach(name string, conn *Connection)*Hub{
	hub := Hubs[name]
	
	if hub == nil{
		temp := make(map[int]*Connection)
		temp[conn.getID()] = conn
		Hubs[name] = &Hub{
			conns: temp,
			Name: name,
		}
		return Hubs[name]
	}else{
		Hubs[name].attach(conn)
		return Hubs[name]
	}
}

