package model

import (
	"github.com/princeofthesky/example_chat/trace_log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait         = 5 * time.Second
	commandBlockChat  = "71"
	commandEnableChat = "72"
)

type WsConnection struct {
	Id          string
	Name        string
	Ucid        string
	Ws          *websocket.Conn
	mu          sync.Mutex
	StopConn    chan bool
	Stop        bool
	Enable      bool
	Expired     bool
	ExpireTime  int64
	waitingSend [][]byte
	length      int
	rooms       map[string]*Room
	roomMutex   sync.RWMutex
}

func NewConnection(ws *websocket.Conn, username string, userConnectionId string) *WsConnection {
	connect := WsConnection{
		Id:          uuid.NewString(),
		StopConn:    make(chan bool),
		Ws:          ws,
		Ucid:        userConnectionId,
		Name:        username,
		Enable:      true,
		Expired:     false,
		Stop:        false,
		rooms:       make(map[string]*Room),
		waitingSend: [][]byte{},
		length:      0,
	}
	go connect.loopWrite()
	return &connect
}

func (c *WsConnection) loopWrite() {
	defer func() {
		trace_log.Logger.Info().Str("c.Id", c.Id).Msg("Stop loop write ")
		close(c.StopConn)
		c.Ws.Close()
	}()
	for !c.Stop {
		waiting := c.getListWaitingMsg()
		if len(waiting) == 0 {
			time.Sleep(20 * time.Millisecond)
			continue
		}
		c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
		for _, data := range waiting {
			err := c.Ws.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				trace_log.Logger.Error().Err(err).Str("c.Id", c.Id).Msg("err when write message")
				c.Stop = true
				return
			}
		}
	}
}

func (c *WsConnection) getListWaitingMsg() [][]byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	waiting := c.waitingSend
	c.waitingSend = [][]byte{}
	return waiting
}

func (c *WsConnection) AddWaitingMsg(mt int, payload []byte) {
	if payload == nil || c.Stop {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.waitingSend = append(c.waitingSend, payload)
}

func (conn *WsConnection) SubscribeRoom(room *Room) {
	conn.roomMutex.Lock()
	defer conn.roomMutex.Unlock()
	conn.rooms[room.Id] = room
}

func (conn *WsConnection) UnSubscribeRoom(roomId string) {
	conn.roomMutex.Lock()
	defer conn.roomMutex.Unlock()
	delete(conn.rooms, roomId)
}
func (conn *WsConnection) GetListRoom() []*Room {
	conn.roomMutex.RLock()
	defer conn.roomMutex.RUnlock()
	rooms := make([]*Room, len(conn.rooms))
	i := 0
	for _, room := range conn.rooms {
		rooms[i] = room
		i++
	}
	return rooms
}
