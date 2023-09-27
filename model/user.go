package model

import (
	"github.com/princeofthesky/example_chat/skiplist"
)

type User struct {
	Id          string
	Connections *skiplist.Skiplist[string, *WsConnection]
}

func (user *User) AddConnection(conn *WsConnection) {
	user.Connections.Insert(conn.Id, conn)
}

func (user *User) RemoveConnection(connId string) {
	user.Connections.Remove(connId)
}

func (user *User) CountConnection() int {
	return user.Connections.Len()
}
