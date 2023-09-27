package model

import (
	"github.com/liyue201/gostl/utils/comparator"
	"github.com/princeofthesky/example_chat/skiplist"
	"sync"
	"time"
)

const maxHash = 100

type Room struct {
	Id          string
	Connections map[int]*skiplist.Skiplist[string, *WsConnection]
	Users       map[string]struct{}
	mu          sync.RWMutex
	Close       bool
}

func InitANewRoom(roomId string) *Room {
	connections := make(map[int]*skiplist.Skiplist[string, *WsConnection])
	for i := 0; i < maxHash; i++ {
		connections[i] = skiplist.New[string, *WsConnection](comparator.StringComparator)
	}
	room := &Room{
		Id:          roomId,
		Connections: connections,
		Users:       make(map[string]struct{}),
	}
	go room.updateListUser()
	return room
}
func (room *Room) updateListUser() {
	for !room.Close {
		time.Sleep(10 * time.Second)
		updated := map[string]struct{}{}
		for _, connections := range room.Connections {
			connections.Traversal(func(connId string, conn *WsConnection) bool {
				updated[conn.Name] = struct{}{}
				return true
			})
		}
		room.Users = updated
	}
}

func (room *Room) AddConnection(conn *WsConnection) {
	if conn == nil {
		return
	}
	hashValue := hashToRange(conn.Id)
	room.Connections[hashValue].Insert(conn.Id, conn)
}

func (room *Room) RemoveConnection(connId string) {
	hashValue := hashToRange(connId)
	room.Connections[hashValue].Remove(connId)
}

func (room *Room) SendMsgToAll(msg *Message) {
	payload := msg.JSONEncode()
	for _, connections := range room.Connections {
		connections.Traversal(func(connId string, conn *WsConnection) bool {
			conn.AddWaitingMsg(1, payload)
			return true
		})
	}
}
