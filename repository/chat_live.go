package repository

import (
	"github.com/liyue201/gostl/utils/comparator"
	"github.com/princeofthesky/example_chat/model"
	"github.com/princeofthesky/example_chat/skiplist"
	"github.com/princeofthesky/example_chat/token"
	"github.com/princeofthesky/example_chat/trace_log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	pingPeriod               = 500 * time.Second
	list_connect_last_scribe = "last_subscibe_for_user"
	list_user_active         = "chatroom_online_count_realtime"
	room_base_number         = "chatroom_online_count"
	public_channel           = "public_channel_skyx"
	private_channel          = "private_channel_skyx"
	SystemChannel            = "system_skyx_channel"
	live_channel             = "live_channel_skyx"
)

var chatChannel = []string{"chat_channel_1", "chat_channel_2", "chat_channel_3"}

const (
	commandSubscribe   = "0"
	commandUnsubscribe = "1"
	commandChat        = "2"
	PongWait           = 60 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type SocketLive struct {
	instanceName string
	livePrefix   string
	prefix       string
	//ListConnection[username][conn.id] -> socket : dung de quan ly connection tren he thong
	// dung de gui message den tat ca connection tren he thong
	//
	listUser *skiplist.Skiplist[string, *model.User]
	listRoom *skiplist.Skiplist[string, *model.Room]

	// subcribe /unscribe vao 1 room
	// send message den user trong 1 user room
	// count so luong user online cua tung room theo socket
	// check 1 user online hay khong theo socket
	// cac tinh nang xoa tin nhan cua 1 user trong 1 room
	// xoa tin nhan cua toan bo phong chat
	// xoa 1 tin nhan bat ki
	// chan / mo chat cua 1 user trong tat ca cac room
	// chan / mo chat cua toan bo user
	//
	userMutex sync.RWMutex
	roomMutex sync.RWMutex
}

func NewSocketLive(prefix string, livePrefix string, instanceName string) *SocketLive {

	var socket = SocketLive{
		instanceName: instanceName,
		livePrefix:   livePrefix,
		prefix:       prefix,
		listUser:     skiplist.New[string, *model.User](comparator.StringComparator, skiplist.WithGoroutineSafe()),
		listRoom:     skiplist.New[string, *model.Room](comparator.StringComparator, skiplist.WithGoroutineSafe()),
	}
	return &socket
}

func (sk *SocketLive) CreateConnection(w http.ResponseWriter, r *http.Request, username string, userConnectionId string, tokenInfo *token.TokenJWTInfo) (bool, error) {
	if userConnectionId == "" {
		userConnectionId = uuid.NewString()
	}
	trace_log.Logger.Info().Str("username", username).Str("userConnectionId", userConnectionId).Msg("create new connection")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		trace_log.Logger.Error().Err(err).Msg("err when try run upgrader")
		// fmt.Println(err)
		return false, err
	}
	conn := model.NewConnection(ws, username, userConnectionId)
	user := sk.GetUser(username)
	if user == nil {
		user = sk.GetOrInitUser(username)
	}
	user.AddConnection(conn)
	ws.SetCloseHandler(func(code int, text string) error {
		sk.CloseAConnectionData(user, conn)
		return nil
	})
	go handleReadData(sk, user, conn)
	return true, nil
}

func (sk *SocketLive) CloseAConnectionData(user *model.User, conn *model.WsConnection) {
	user.RemoveConnection(conn.Id)
	rooms := conn.GetListRoom()
	for _, room := range rooms {
		room.RemoveConnection(conn.Id)
		conn.UnSubscribeRoom(room.Id)
	}
	sk.RemoveUser(user)
}
func (sk *SocketLive) GetRoom(channel string) *model.Room {
	user, err := sk.listRoom.Get(channel)
	if err != nil {
		return nil
	}
	return user
}
func (sk *SocketLive) GetUser(userName string) *model.User {
	user, err := sk.listUser.Get(userName)
	if err != nil {
		return nil
	}
	return user
}

func (sk *SocketLive) RemoveUser(user *model.User) {
	sk.userMutex.Lock()
	defer sk.userMutex.Unlock()
	if user.CountConnection() > 0 {
		return
	}
	sk.listUser.Remove(user.Id)
}

func (sk *SocketLive) GetOrInitUser(userId string) *model.User {
	user, err := sk.listUser.Get(userId)
	if err == nil && user != nil {
		return user
	}
	sk.userMutex.Lock()
	defer sk.userMutex.Unlock()
	user, err = sk.listUser.Get(userId)
	if err == nil && user != nil {
		return user
	}
	user = &model.User{
		Id:          userId,
		Connections: skiplist.New[string, *model.WsConnection](comparator.StringComparator, skiplist.WithGoroutineSafe()),
	}
	sk.listUser.Insert(userId, user)
	return user
}
func (sk *SocketLive) GetOrInitRoom(channel string) *model.Room {
	room, err := sk.listRoom.Get(channel)
	if room != nil && err == nil {
		return room
	}
	sk.roomMutex.Lock()
	defer sk.roomMutex.Unlock()
	room, err = sk.listRoom.Get(channel)
	if room != nil && err == nil {
		return room
	}
	room = model.InitANewRoom(channel)
	sk.listRoom.Insert(channel, room)
	return room
}

func (sk *SocketLive) Subscribe(user *model.User, channel string, conn *model.WsConnection) bool {
	room := sk.GetOrInitRoom(channel)
	room.AddConnection(conn)
	conn.SubscribeRoom(room)
	return true
}

func (sk *SocketLive) Unsubscribe(user *model.User, channel string, conn *model.WsConnection) {
	room := sk.GetRoom(channel)
	if room == nil {
		return
	}
	conn.UnSubscribeRoom(channel)
	room.RemoveConnection(conn.Id)
}

func (sub *SocketLive) getPublicChannelList() []string {
	var c []string
	c = append(c, sub.prefix+public_channel)
	return c
}

func (sub *SocketLive) getPrivateChannelList() []string {
	var c []string
	c = append(c, sub.prefix+private_channel)
	return c
}
